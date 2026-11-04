//! Rebuilds a document with the content of every op in a deleted container
//! replaced by placeholders, keeping op IDs, lamports and deps. Works on the
//! owned `JsonSchema` structs, so no JSON text is produced and floats stay
//! exact.

use std::collections::{HashMap, HashSet};
use std::panic::{catch_unwind, AssertUnwindSafe};
use std::sync::Arc;

use loro::{
    ContainerID, ContainerTrait, ContainerType, JsonFutureOp, JsonListOp, JsonMapOp,
    JsonMovableListOp, JsonOpContent, JsonSchema, JsonTextOp, JsonTreeOp, LoroDoc, LoroValue,
    VersionVector,
};
use uniffi::Handle;

const REPLACEMENT_CHAR: char = '\u{FFFD}';

const OK: i32 = 0;
const ERR_SHALLOW: i32 = 1;
const ERR_DETACHED: i32 = 2;
const ERR_UNKNOWN_OP: i32 = 3;
const ERR_REBUILD: i32 = 4;
const ERR_VERIFY: i32 = 5;
const ERR_PANIC: i32 = 6;

#[repr(C)]
pub struct LoroGoRedactResult {
    /// uniffi handle of the rebuilt `LoroDoc`, owned by the caller; 0 when
    /// nothing needed redacting or on error.
    pub doc: u64,
    pub redacted_ops: u64,
    pub redacted_containers: u64,
    pub skipped_revivable_containers: u64,
    pub error_code: i32,
    pub error_msg: *mut u8,
    pub error_len: usize,
}

#[derive(Default)]
struct Stats {
    redacted_ops: u64,
    redacted_containers: u64,
    skipped_revivable_containers: u64,
}

struct Failure {
    code: i32,
    msg: String,
}

impl Failure {
    fn new(code: i32, msg: impl Into<String>) -> Self {
        Self {
            code,
            msg: msg.into(),
        }
    }
}

/// # Safety
///
/// `doc` must be a uniffi `LoroDoc` handle that this call may consume: the
/// Go side passes the clone `incrementPointer` returns, as every generated
/// method does.
#[no_mangle]
pub unsafe extern "C" fn loro_go_redact_deleted_containers(
    doc: u64,
    include_revivable: u8,
) -> LoroGoRedactResult {
    let src: Arc<loro_ffi::LoroDoc> = Handle::from_raw_unchecked(doc).into_arc();
    let outcome = catch_unwind(AssertUnwindSafe(|| {
        redact_deleted_containers(&src, include_revivable != 0)
    }));
    drop(src);

    let mut result = LoroGoRedactResult {
        doc: 0,
        redacted_ops: 0,
        redacted_containers: 0,
        skipped_revivable_containers: 0,
        error_code: OK,
        error_msg: std::ptr::null_mut(),
        error_len: 0,
    };
    let failure = match outcome {
        Ok(Ok((rebuilt, stats))) => {
            result.redacted_ops = stats.redacted_ops;
            result.redacted_containers = stats.redacted_containers;
            result.skipped_revivable_containers = stats.skipped_revivable_containers;
            if let Some(rebuilt) = rebuilt {
                result.doc = Handle::from_arc(Arc::new(rebuilt)).as_raw();
            }
            return result;
        }
        Ok(Err(f)) => f,
        Err(panic) => {
            let msg = panic
                .downcast_ref::<&str>()
                .map(|s| s.to_string())
                .or_else(|| panic.downcast_ref::<String>().cloned())
                .unwrap_or_else(|| "panic".to_string());
            Failure::new(ERR_PANIC, msg)
        }
    };
    let mut bytes = failure.msg.into_bytes().into_boxed_slice();
    result.error_code = failure.code;
    result.error_len = bytes.len();
    result.error_msg = bytes.as_mut_ptr();
    std::mem::forget(bytes);
    result
}

/// # Safety
///
/// `data` and `len` must come from one `LoroGoRedactResult`, freed once.
#[no_mangle]
pub unsafe extern "C" fn loro_go_free_redact_error(data: *mut u8, len: usize) {
    if !data.is_null() {
        drop(Box::from_raw(std::ptr::slice_from_raw_parts_mut(data, len)));
    }
}

fn redact_deleted_containers(
    src: &LoroDoc,
    include_revivable: bool,
) -> Result<(Option<loro_ffi::LoroDoc>, Stats), Failure> {
    if src.is_shallow() {
        return Err(Failure::new(ERR_SHALLOW, "document is shallow"));
    }
    // A detached doc reports deletion at its checkout, not at head.
    if src.is_detached() {
        return Err(Failure::new(ERR_DETACHED, "document is detached"));
    }

    // The export would commit an open transaction past the version vector
    // read here.
    src.commit();
    let oplog_vv = src.oplog_vv();
    let mut schema = src
        .export_json_updates_without_peer_compression(&VersionVector::new(), &oplog_vv)
        .map_err(|e| Failure::new(ERR_REBUILD, e.to_string()))?;
    let parents = container_parents(&schema);

    let mut stats = Stats::default();
    // Whether each container's ops are redacted.
    let mut decisions: HashMap<ContainerID, bool> = HashMap::new();
    let mut redacted: HashSet<ContainerID> = HashSet::new();
    for change in schema.changes.iter_mut() {
        for op in change.ops.iter_mut() {
            let redact = match decisions.get(&op.container) {
                Some(d) => *d,
                None => {
                    let d =
                        should_redact(src, &op.container, &parents, include_revivable, &mut stats);
                    decisions.insert(op.container.clone(), d);
                    d
                }
            };
            if !redact {
                continue;
            }
            if redact_op(&mut op.content)
                .map_err(|msg| Failure::new(ERR_UNKNOWN_OP, format!("{}: {msg}", op.container)))?
            {
                stats.redacted_ops += 1;
                redacted.insert(op.container.clone());
            }
        }
    }
    stats.redacted_containers = redacted.len() as u64;
    if stats.redacted_ops == 0 {
        return Ok((None, stats));
    }

    let rebuilt_ffi = loro_ffi::LoroDoc::new();
    // loro_ffi's inherent methods shadow the ones on the Loro doc it wraps.
    let rebuilt: &LoroDoc = &rebuilt_ffi;
    let status = rebuilt
        .import_json_updates(schema)
        .map_err(|e| Failure::new(ERR_REBUILD, e.to_string()))?;
    if status.pending.is_some() {
        return Err(Failure::new(ERR_REBUILD, "rebuilt history has pending ops"));
    }
    // Ops still waiting for a dependency are not in the exported history.
    let pending = src
        .export_pending_updates()
        .map_err(|e| Failure::new(ERR_REBUILD, e.to_string()))?;
    if let Some(bytes) = pending {
        rebuilt
            .import(&bytes)
            .map_err(|e| Failure::new(ERR_REBUILD, e.to_string()))?;
    }
    rebuilt
        .set_peer_id(src.peer_id())
        .map_err(|e| Failure::new(ERR_REBUILD, e.to_string()))?;

    if rebuilt.oplog_vv() != oplog_vv {
        return Err(Failure::new(ERR_VERIFY, "oplog version vector differs"));
    }
    if rebuilt.oplog_frontiers() != src.oplog_frontiers() {
        return Err(Failure::new(ERR_VERIFY, "oplog frontiers differ"));
    }
    if rebuilt.get_deep_value() != src.get_deep_value() {
        return Err(Failure::new(ERR_VERIFY, "deep value differs"));
    }
    Ok((Some(rebuilt_ffi), stats))
}

fn should_redact(
    doc: &LoroDoc,
    id: &ContainerID,
    parents: &HashMap<ContainerID, (ContainerID, ContainerType)>,
    include_revivable: bool,
    stats: &mut Stats,
) -> bool {
    let deleted = doc
        .get_container(id.clone())
        .is_some_and(|c| c.is_deleted());
    if !deleted {
        return false;
    }
    if include_revivable || !is_revivable(id, parents) {
        return true;
    }
    stats.skipped_revivable_containers += 1;
    false
}

/// Whether a later op could make the container reachable again: Loro lets a
/// tree or movable-list move bring back what sits below it. A chain whose
/// parent is not in the history counts as revivable.
fn is_revivable(
    id: &ContainerID,
    parents: &HashMap<ContainerID, (ContainerID, ContainerType)>,
) -> bool {
    let mut current = id;
    loop {
        if current.is_root() {
            return false;
        }
        let Some((parent, parent_type)) = parents.get(current) else {
            return true;
        };
        if matches!(
            parent_type,
            ContainerType::Tree | ContainerType::MovableList
        ) {
            return true;
        }
        current = parent;
    }
}

/// Parent of every container created in the history: container values
/// inserted into lists and maps, and the meta map of every tree node.
fn container_parents(schema: &JsonSchema) -> HashMap<ContainerID, (ContainerID, ContainerType)> {
    let mut parents = HashMap::new();
    for change in &schema.changes {
        for op in &change.ops {
            let parent = (op.container.clone(), op.container.container_type());
            let mut add = |v: &LoroValue| {
                if let LoroValue::Container(child) = v {
                    parents.insert(child.clone(), parent.clone());
                }
            };
            match &op.content {
                JsonOpContent::List(JsonListOp::Insert { value, .. })
                | JsonOpContent::MovableList(JsonMovableListOp::Insert { value, .. }) => {
                    value.iter().for_each(&mut add)
                }
                JsonOpContent::MovableList(JsonMovableListOp::Set { value, .. })
                | JsonOpContent::Map(JsonMapOp::Insert { value, .. }) => add(value),
                JsonOpContent::Tree(JsonTreeOp::Create { target, .. }) => {
                    parents.insert(target.associated_meta_container(), parent.clone());
                }
                _ => {}
            }
        }
    }
    parents
}

/// Replaces the op's content the way `loro::json::redact` does, over the
/// whole op. Returns whether anything changed, so an already redacted
/// history needs no rebuild.
fn redact_op(content: &mut JsonOpContent) -> Result<bool, &'static str> {
    let changed = match content {
        JsonOpContent::List(JsonListOp::Insert { value, .. })
        | JsonOpContent::MovableList(JsonMovableListOp::Insert { value, .. }) => value
            .iter_mut()
            .fold(false, |acc, v| redact_value(v) || acc),
        JsonOpContent::MovableList(JsonMovableListOp::Set { value, .. })
        | JsonOpContent::Map(JsonMapOp::Insert { value, .. }) => redact_value(value),
        JsonOpContent::Text(JsonTextOp::Insert { text, .. }) => {
            if text.chars().all(|c| c == REPLACEMENT_CHAR) {
                false
            } else {
                *text = text.chars().map(|_| REPLACEMENT_CHAR).collect();
                true
            }
        }
        JsonOpContent::Text(JsonTextOp::Mark { style_value, .. }) => redact_value(style_value),
        JsonOpContent::Future(wrapper) => match wrapper.value {
            // A counter's value is a number, never content worth removing.
            JsonFutureOp::Counter(_) => false,
            JsonFutureOp::Unknown(_) => return Err("op type unknown to this Loro version"),
        },
        JsonOpContent::List(JsonListOp::Delete { .. })
        | JsonOpContent::MovableList(
            JsonMovableListOp::Delete { .. } | JsonMovableListOp::Move { .. },
        )
        | JsonOpContent::Map(JsonMapOp::Delete { .. })
        | JsonOpContent::Text(JsonTextOp::Delete { .. } | JsonTextOp::MarkEnd)
        | JsonOpContent::Tree(_) => false,
    };
    Ok(changed)
}

fn redact_value(v: &mut LoroValue) -> bool {
    match v {
        LoroValue::Container(_) | LoroValue::Null => false,
        _ => {
            *v = LoroValue::Null;
            true
        }
    }
}
