package loro

// This file adds single-pointer wrappers around generated methods that return
// `**T` (UniFFI's representation of `Option<Arc<T>>`). Each wrapper returns
// `*T`, which is nil when the value is absent — the idiomatic Go shape.
//
// The generated `**T` methods (Get, Doc, GetAttached, GetCursor, TryGet*, etc.)
// remain available; these helpers are additive. For uncovered cases use
// loro.Deref.

// Lookup returns the value at key, or nil if the key is absent.
//
//	if v := m.Lookup("name"); v != nil {
//		// use v
//	}
func (m *LoroMap) Lookup(key string) *ValueOrContainer {
	return Deref(m.Get(key))
}

// At returns the value at index, or nil if the index is out of bounds.
func (l *LoroList) At(index uint32) *ValueOrContainer {
	return Deref(l.Get(index))
}

// At returns the value at index, or nil if the index is out of bounds.
func (ml *LoroMovableList) At(index uint32) *ValueOrContainer {
	return Deref(ml.Get(index))
}

// OwnerDoc returns the LoroDoc this container is attached to, or nil if it is
// detached.
func (c *LoroCounter) OwnerDoc() *LoroDoc { return Deref(c.Doc()) }

// OwnerDoc returns the LoroDoc this container is attached to, or nil if it is
// detached.
func (l *LoroList) OwnerDoc() *LoroDoc { return Deref(l.Doc()) }

// OwnerDoc returns the LoroDoc this container is attached to, or nil if it is
// detached.
func (m *LoroMap) OwnerDoc() *LoroDoc { return Deref(m.Doc()) }

// OwnerDoc returns the LoroDoc this container is attached to, or nil if it is
// detached.
func (ml *LoroMovableList) OwnerDoc() *LoroDoc { return Deref(ml.Doc()) }

// OwnerDoc returns the LoroDoc this container is attached to, or nil if it is
// detached.
func (t *LoroText) OwnerDoc() *LoroDoc { return Deref(t.Doc()) }

// OwnerDoc returns the LoroDoc this container is attached to, or nil if it is
// detached.
func (tr *LoroTree) OwnerDoc() *LoroDoc { return Deref(tr.Doc()) }

// Attached returns the live attached version of this container, or nil if the
// container has been detached and has no attached counterpart.
func (c *LoroCounter) Attached() *LoroCounter { return Deref(c.GetAttached()) }

// Attached returns the live attached version of this container, or nil if the
// container has been detached and has no attached counterpart.
func (l *LoroList) Attached() *LoroList { return Deref(l.GetAttached()) }

// Attached returns the live attached version of this container, or nil if the
// container has been detached and has no attached counterpart.
func (m *LoroMap) Attached() *LoroMap { return Deref(m.GetAttached()) }

// Attached returns the live attached version of this container, or nil if the
// container has been detached and has no attached counterpart.
func (ml *LoroMovableList) Attached() *LoroMovableList { return Deref(ml.GetAttached()) }

// Attached returns the live attached version of this container, or nil if the
// container has been detached and has no attached counterpart.
func (t *LoroText) Attached() *LoroText { return Deref(t.GetAttached()) }

// Attached returns the live attached version of this container, or nil if the
// container has been detached and has no attached counterpart.
func (tr *LoroTree) Attached() *LoroTree { return Deref(tr.GetAttached()) }

// Cursor creates a stable cursor at pos with the given side, or nil if a
// cursor cannot be created at that position.
func (l *LoroList) Cursor(pos uint32, side Side) *Cursor {
	return Deref(l.GetCursor(pos, side))
}

// Cursor creates a stable cursor at pos with the given side, or nil if a
// cursor cannot be created at that position.
func (ml *LoroMovableList) Cursor(pos uint32, side Side) *Cursor {
	return Deref(ml.GetCursor(pos, side))
}

// Cursor creates a stable cursor at pos with the given side, or nil if a
// cursor cannot be created at that position.
func (t *LoroText) Cursor(pos uint32, side Side) *Cursor {
	return Deref(t.GetCursor(pos, side))
}

// FindCounter returns the counter at id, or nil if no such container exists.
// Unlike GetCounter, it does not create the container if absent.
func (d *LoroDoc) FindCounter(id ContainerIdLike) *LoroCounter {
	return Deref(d.TryGetCounter(id))
}

// FindList returns the list at id, or nil if no such container exists.
// Unlike GetList, it does not create the container if absent.
func (d *LoroDoc) FindList(id ContainerIdLike) *LoroList {
	return Deref(d.TryGetList(id))
}

// FindMap returns the map at id, or nil if no such container exists.
// Unlike GetMap, it does not create the container if absent.
func (d *LoroDoc) FindMap(id ContainerIdLike) *LoroMap {
	return Deref(d.TryGetMap(id))
}

// FindMovableList returns the movable list at id, or nil if no such container
// exists. Unlike GetMovableList, it does not create the container if absent.
func (d *LoroDoc) FindMovableList(id ContainerIdLike) *LoroMovableList {
	return Deref(d.TryGetMovableList(id))
}

// FindText returns the text container at id, or nil if no such container
// exists. Unlike GetText, it does not create the container if absent.
func (d *LoroDoc) FindText(id ContainerIdLike) *LoroText {
	return Deref(d.TryGetText(id))
}

// FindTree returns the tree at id, or nil if no such container exists.
// Unlike GetTree, it does not create the container if absent.
func (d *LoroDoc) FindTree(id ContainerIdLike) *LoroTree {
	return Deref(d.TryGetTree(id))
}

// FindByPath returns the value or container at the given index path, or nil
// if the path does not resolve.
func (d *LoroDoc) FindByPath(path []Index) *ValueOrContainer {
	return Deref(d.GetByPath(path))
}

// FindByStrPath returns the value or container at the given string path, or
// nil if the path does not resolve.
func (d *LoroDoc) FindByStrPath(path string) *ValueOrContainer {
	return Deref(d.GetByStrPath(path))
}

// FindContainer returns the value or container with the given id, or nil if
// no such container exists.
func (d *LoroDoc) FindContainer(id ContainerId) *ValueOrContainer {
	return Deref(d.GetContainer(id))
}
