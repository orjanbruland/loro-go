package loro

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// String formatting for the generated types. All implementations satisfy
// fmt.Stringer so that `%v` / `%s` on logged values produces readable output
// instead of struct literals or bare integers.
//
// ID-shaped types use Loro's canonical `<counter>@<peer>` ordering so log
// output is greppable across the Rust and Go sides.

func (ContainerTypeText) String() string        { return "Text" }
func (ContainerTypeMap) String() string         { return "Map" }
func (ContainerTypeList) String() string        { return "List" }
func (ContainerTypeMovableList) String() string { return "MovableList" }
func (ContainerTypeTree) String() string        { return "Tree" }
func (ContainerTypeCounter) String() string     { return "Counter" }
func (e ContainerTypeUnknown) String() string   { return fmt.Sprintf("Unknown(%d)", e.Kind) }

func (e ContainerIdRoot) String() string {
	return fmt.Sprintf("cid:root-%s:%s", e.Name, e.ContainerType)
}

func (e ContainerIdNormal) String() string {
	return fmt.Sprintf("cid:%d@%d:%s", e.Counter, e.Peer, e.ContainerType)
}

func (id Id) String() string     { return fmt.Sprintf("%d@%d", id.Counter, id.Peer) }
func (id TreeId) String() string { return fmt.Sprintf("%d@%d", id.Counter, id.Peer) }
func (id IdLp) String() string   { return fmt.Sprintf("%dL@%d", id.Lamport, id.Peer) }

func (s CounterSpan) String() string { return fmt.Sprintf("[%d..%d)", s.Start, s.End) }
func (s IdSpan) String() string      { return fmt.Sprintf("%d:%s", s.Peer, s.Counter) }

func (e TreeParentIdNode) String() string { return fmt.Sprintf("node(%s)", e.Id) }
func (TreeParentIdRoot) String() string   { return "root" }
func (TreeParentIdDeleted) String() string {
	return "deleted"
}
func (TreeParentIdUnexist) String() string { return "unexist" }

func (e IndexKey) String() string  { return fmt.Sprintf("%q", e.Key) }
func (e IndexSeq) String() string  { return fmt.Sprintf("[%d]", e.Index) }
func (e IndexNode) String() string { return fmt.Sprintf("node(%s)", e.Target) }

func (LoroValueNull) String() string        { return "null" }
func (e LoroValueBool) String() string      { return fmt.Sprintf("%t", e.Value) }
func (e LoroValueDouble) String() string    { return fmt.Sprintf("%g", e.Value) }
func (e LoroValueI64) String() string       { return fmt.Sprintf("%d", e.Value) }
func (e LoroValueBinary) String() string    { return fmt.Sprintf("binary(%dB)", len(e.Value)) }
func (e LoroValueString) String() string    { return fmt.Sprintf("%q", e.Value) }
func (e LoroValueContainer) String() string { return fmt.Sprint(e.Value) }

func (e LoroValueList) String() string {
	parts := make([]string, len(e.Value))
	for i, v := range e.Value {
		parts[i] = fmt.Sprint(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func (e LoroValueMap) String() string {
	keys := make([]string, 0, len(e.Value))
	for k := range e.Value {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%q: %s", k, e.Value[k])
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func (e TextDeltaRetain) String() string {
	if e.Attributes != nil {
		return fmt.Sprintf("retain(%d, %s)", e.Retain, formatAttrs(*e.Attributes))
	}
	return fmt.Sprintf("retain(%d)", e.Retain)
}

func (e TextDeltaInsert) String() string {
	if e.Attributes != nil {
		return fmt.Sprintf("insert(%q, %s)", e.Insert, formatAttrs(*e.Attributes))
	}
	return fmt.Sprintf("insert(%q)", e.Insert)
}

func (e TextDeltaDelete) String() string { return fmt.Sprintf("delete(%d)", e.Delete) }

func formatAttrs(m map[string]LoroValue) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%q: %s", k, m[k])
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func (e TreeExternalDiffCreate) String() string {
	return fmt.Sprintf("create(parent=%s, index=%d, fi=%q)", e.Parent, e.Index, e.FractionalIndex)
}

func (e TreeExternalDiffMove) String() string {
	return fmt.Sprintf("move(parent=%s, index=%d, fi=%q, from=%s/%d)",
		e.Parent, e.Index, e.FractionalIndex, e.OldParent, e.OldIndex)
}

func (e TreeExternalDiffDelete) String() string {
	return fmt.Sprintf("delete(from=%s/%d)", e.OldParent, e.OldIndex)
}

func (s Side) String() string {
	switch s {
	case SideLeft:
		return "Left"
	case SideMiddle:
		return "Middle"
	case SideRight:
		return "Right"
	default:
		return fmt.Sprintf("Side(%d)", uint(s))
	}
}

func (u UndoOrRedo) String() string {
	switch u {
	case UndoOrRedoUndo:
		return "Undo"
	case UndoOrRedoRedo:
		return "Redo"
	default:
		return fmt.Sprintf("UndoOrRedo(%d)", uint(u))
	}
}

func (k EventTriggerKind) String() string {
	switch k {
	case EventTriggerKindLocal:
		return "Local"
	case EventTriggerKindImport:
		return "Import"
	case EventTriggerKindCheckout:
		return "Checkout"
	default:
		return fmt.Sprintf("EventTriggerKind(%d)", uint(k))
	}
}

func (e ExpandType) String() string {
	switch e {
	case ExpandTypeBefore:
		return "Before"
	case ExpandTypeAfter:
		return "After"
	case ExpandTypeBoth:
		return "Both"
	case ExpandTypeNone:
		return "None"
	default:
		return fmt.Sprintf("ExpandType(%d)", uint(e))
	}
}

func (p PosType) String() string {
	switch p {
	case PosTypeBytes:
		return "Bytes"
	case PosTypeUnicode:
		return "Unicode"
	case PosTypeUtf16:
		return "Utf16"
	case PosTypeEvent:
		return "Event"
	case PosTypeEntity:
		return "Entity"
	default:
		return fmt.Sprintf("PosType(%d)", uint(p))
	}
}

func (o Ordering) String() string {
	switch o {
	case OrderingLess:
		return "Less"
	case OrderingEqual:
		return "Equal"
	case OrderingGreater:
		return "Greater"
	default:
		return fmt.Sprintf("Ordering(%d)", uint(o))
	}
}

// --- Object types ---

func (f *Frontiers) String() string {
	if f == nil {
		return "Frontiers(<nil>)"
	}
	ids := f.ToVec()
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = id.String()
	}
	return "Frontiers[" + strings.Join(parts, ", ") + "]"
}

func (v *VersionVector) String() string {
	if v == nil {
		return "VersionVector(<nil>)"
	}
	m := v.ToHashmap()
	peers := make([]uint64, 0, len(m))
	for p := range m {
		peers = append(peers, p)
	}
	sort.Slice(peers, func(i, j int) bool { return peers[i] < peers[j] })
	parts := make([]string, len(peers))
	for i, p := range peers {
		parts[i] = fmt.Sprintf("%d:%d", p, m[p])
	}
	return "VersionVector{" + strings.Join(parts, ", ") + "}"
}

func (c *Cursor) String() string {
	if c == nil {
		return "Cursor(<nil>)"
	}
	return "Cursor(" + hex.EncodeToString(c.Encode()) + ")"
}
