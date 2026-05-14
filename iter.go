package loro

import "iter"

// Range-over-func iterators for the container types. These let callsites use
// Go 1.23+ for-range without first materializing a slice via ToVec / Values.
// Each iterator stops early when the loop body breaks.

// All yields each (index, element) pair in order. The element is a
// *ValueOrContainer so nested containers are reachable without a second
// lookup. Len is sampled once at the start.
//
//	for i, v := range list.All() {
//		// use v
//	}
func (l *LoroList) All() iter.Seq2[int, *ValueOrContainer] {
	return func(yield func(int, *ValueOrContainer) bool) {
		n := l.Len()
		for i := uint32(0); i < n; i++ {
			if !yield(int(i), l.At(i)) {
				return
			}
		}
	}
}

// All yields each (index, element) pair in order. The element is a
// *ValueOrContainer so nested containers are reachable without a second
// lookup. Len is sampled once at the start.
func (ml *LoroMovableList) All() iter.Seq2[int, *ValueOrContainer] {
	return func(yield func(int, *ValueOrContainer) bool) {
		n := ml.Len()
		for i := uint32(0); i < n; i++ {
			if !yield(int(i), ml.At(i)) {
				return
			}
		}
	}
}

// All yields each (key, value) pair. Iteration order matches the underlying
// Keys() result; a separate FFI lookup is issued per key, so breaking out
// early avoids the remaining lookups.
//
//	for k, v := range m.All() {
//		// use k, v
//	}
func (m *LoroMap) All() iter.Seq2[string, *ValueOrContainer] {
	return func(yield func(string, *ValueOrContainer) bool) {
		for _, key := range m.Keys() {
			if !yield(key, m.Lookup(key)) {
				return
			}
		}
	}
}

// AllNodes yields every node id in the tree, in the order returned by Nodes.
func (t *LoroTree) AllNodes() iter.Seq[TreeId] {
	return func(yield func(TreeId) bool) {
		for _, id := range t.Nodes() {
			if !yield(id) {
				return
			}
		}
	}
}

// AllRoots yields every root node id, in the order returned by Roots.
func (t *LoroTree) AllRoots() iter.Seq[TreeId] {
	return func(yield func(TreeId) bool) {
		for _, id := range t.Roots() {
			if !yield(id) {
				return
			}
		}
	}
}

// AllChildren yields each direct child of parent. Yields nothing if parent is
// unknown or has no children.
func (t *LoroTree) AllChildren(parent TreeParentId) iter.Seq[TreeId] {
	return func(yield func(TreeId) bool) {
		children := t.Children(parent)
		if children == nil {
			return
		}
		for _, id := range *children {
			if !yield(id) {
				return
			}
		}
	}
}
