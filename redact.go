package loro

/*
#include <stdint.h>
#include <stddef.h>

typedef struct {
	uint64_t doc;
	uint64_t redacted_ops;
	uint64_t redacted_containers;
	uint64_t skipped_revivable_containers;
	int32_t error_code;
	uint8_t *error_msg;
	size_t error_len;
} LoroGoRedactResult;

LoroGoRedactResult loro_go_redact_deleted_containers(uint64_t doc, uint8_t include_revivable);
void loro_go_free_redact_error(uint8_t *data, size_t len);
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	// ErrRedactShallow is returned for a shallow document, whose trimmed
	// history cannot be exported for the rebuild.
	ErrRedactShallow = errors.New("loro: cannot redact a shallow document")
	// ErrRedactDetached is returned for a document checked out at an older
	// version, where deletion is not judged at head.
	ErrRedactDetached = errors.New("loro: cannot redact a detached document")
	// ErrRedactUnknownOp is returned when a deleted container holds an op type
	// this Loro version cannot redact.
	ErrRedactUnknownOp = errors.New("loro: deleted container holds an op type that cannot be redacted")
	// ErrRedactVerification is returned when the rebuilt document's version
	// vector, frontiers or deep value differ from the source's.
	ErrRedactVerification = errors.New("loro: redacted document does not match the source")
	// ErrRedactRebuild is returned when exporting or importing the history
	// fails.
	ErrRedactRebuild = errors.New("loro: redaction rebuild failed")
)

// RedactOptions selects which deleted containers RedactDeletedContainers
// redacts.
type RedactOptions struct {
	// IncludeRevivable also redacts containers that sit below a tree or
	// movable-list edge. A later move of that node or element from a peer
	// that has not seen the delete makes such a container live again, and
	// its redacted content then shows as null on this copy while copies with
	// the original history show the old values.
	IncludeRevivable bool
}

// RedactResult describes a RedactDeletedContainers call.
type RedactResult struct {
	// Doc is the rebuilt document, owned by the caller, or nil when no op
	// needed redacting.
	Doc *LoroDoc
	// RedactedOps counts ops whose content was replaced.
	RedactedOps uint64
	// RedactedContainers counts deleted containers with at least one
	// replaced op.
	RedactedContainers uint64
	// SkippedRevivableContainers counts deleted containers left alone
	// because IncludeRevivable was false.
	SkippedRevivableContainers uint64
}

// RedactDeletedContainers rebuilds the document with the content of every op
// in a container that is deleted at head replaced, the way RedactJsonUpdates
// replaces it: text becomes U+FFFD per character, list, map and mark values
// become null, and child container references, deletes, moves and tree ops
// stay. Op IDs, lamports, deps and timestamps are unchanged, so peers holding
// the original history keep syncing with the rebuilt document both ways, as
// long as they exchange deltas from each other's version vector: Loro checks
// an import that brings new ops against every op ID it shares with the
// document, and a redacted op never matches its original, so such an import
// fails with UsedOpId.
//
// The history never leaves Rust as JSON text, so float values are exact and
// one call redacts any number of op spans per peer.
//
// The rebuilt document carries the source's peer ID and pending ops, and its
// oplog version vector, frontiers and deep value are checked against the
// source's before it is returned; a mismatch returns ErrRedactVerification.
// The source document is left as it was, apart from committing an open
// transaction. Document configuration (text styles,
// timestamp recording, merge interval) is not copied.
func (d *LoroDoc) RedactDeletedContainers(opts RedactOptions) (RedactResult, error) {
	var includeRevivable C.uint8_t
	if opts.IncludeRevivable {
		includeRevivable = 1
	}
	handle := d.ffiObject.incrementPointer("*LoroDoc")
	defer d.ffiObject.decrementPointer()
	res := C.loro_go_redact_deleted_containers(handle, includeRevivable)

	if res.error_code != 0 {
		msg := C.GoStringN((*C.char)(unsafe.Pointer(res.error_msg)), C.int(res.error_len))
		C.loro_go_free_redact_error(res.error_msg, res.error_len)
		var base error
		switch res.error_code {
		case 1:
			base = ErrRedactShallow
		case 2:
			base = ErrRedactDetached
		case 3:
			base = ErrRedactUnknownOp
		case 5:
			base = ErrRedactVerification
		default:
			base = ErrRedactRebuild
		}
		return RedactResult{}, fmt.Errorf("%w: %s", base, msg)
	}

	out := RedactResult{
		RedactedOps:                uint64(res.redacted_ops),
		RedactedContainers:         uint64(res.redacted_containers),
		SkippedRevivableContainers: uint64(res.skipped_revivable_containers),
	}
	if res.doc != 0 {
		out.Doc = FfiConverterLoroDocINSTANCE.Lift(res.doc)
	}
	return out, nil
}
