package loro

// Helpers that create and attach a fresh child container in one call, hiding
// the detached-construct-then-attach dance.
//
//	child, err := m.GetOrCreateLoroMap("settings")
//	// instead of:
//	child, err := m.GetOrCreateMapContainer("settings", loro.NewLoroMap())
//
// For LoroList and LoroMovableList the equivalent helpers are named
// InsertLoro<T> and take the position to insert at.

// GetOrCreateLoroMap returns the LoroMap child at key, creating an attached
// empty one if absent.
func (m *LoroMap) GetOrCreateLoroMap(key string) (*LoroMap, error) {
	return m.GetOrCreateMapContainer(key, NewLoroMap())
}

// GetOrCreateLoroList returns the LoroList child at key, creating an attached
// empty one if absent.
func (m *LoroMap) GetOrCreateLoroList(key string) (*LoroList, error) {
	return m.GetOrCreateListContainer(key, NewLoroList())
}

// GetOrCreateLoroMovableList returns the LoroMovableList child at key,
// creating an attached empty one if absent.
func (m *LoroMap) GetOrCreateLoroMovableList(key string) (*LoroMovableList, error) {
	return m.GetOrCreateMovableListContainer(key, NewLoroMovableList())
}

// GetOrCreateLoroText returns the LoroText child at key, creating an attached
// empty one if absent.
func (m *LoroMap) GetOrCreateLoroText(key string) (*LoroText, error) {
	return m.GetOrCreateTextContainer(key, NewLoroText())
}

// GetOrCreateLoroTree returns the LoroTree child at key, creating an attached
// empty one if absent.
func (m *LoroMap) GetOrCreateLoroTree(key string) (*LoroTree, error) {
	return m.GetOrCreateTreeContainer(key, NewLoroTree())
}

// GetOrCreateLoroCounter returns the LoroCounter child at key, creating an
// attached zero-valued one if absent.
func (m *LoroMap) GetOrCreateLoroCounter(key string) (*LoroCounter, error) {
	return m.GetOrCreateCounterContainer(key, NewLoroCounter())
}

// InsertLoroMap inserts a fresh attached LoroMap at pos and returns it.
func (l *LoroList) InsertLoroMap(pos uint32) (*LoroMap, error) {
	return l.InsertMapContainer(pos, NewLoroMap())
}

// InsertLoroList inserts a fresh attached LoroList at pos and returns it.
func (l *LoroList) InsertLoroList(pos uint32) (*LoroList, error) {
	return l.InsertListContainer(pos, NewLoroList())
}

// InsertLoroMovableList inserts a fresh attached LoroMovableList at pos and
// returns it.
func (l *LoroList) InsertLoroMovableList(pos uint32) (*LoroMovableList, error) {
	return l.InsertMovableListContainer(pos, NewLoroMovableList())
}

// InsertLoroText inserts a fresh attached LoroText at pos and returns it.
func (l *LoroList) InsertLoroText(pos uint32) (*LoroText, error) {
	return l.InsertTextContainer(pos, NewLoroText())
}

// InsertLoroTree inserts a fresh attached LoroTree at pos and returns it.
func (l *LoroList) InsertLoroTree(pos uint32) (*LoroTree, error) {
	return l.InsertTreeContainer(pos, NewLoroTree())
}

// InsertLoroCounter inserts a fresh attached LoroCounter at pos and returns
// it.
func (l *LoroList) InsertLoroCounter(pos uint32) (*LoroCounter, error) {
	return l.InsertCounterContainer(pos, NewLoroCounter())
}

// InsertLoroMap inserts a fresh attached LoroMap at pos and returns it.
func (ml *LoroMovableList) InsertLoroMap(pos uint32) (*LoroMap, error) {
	return ml.InsertMapContainer(pos, NewLoroMap())
}

// InsertLoroList inserts a fresh attached LoroList at pos and returns it.
func (ml *LoroMovableList) InsertLoroList(pos uint32) (*LoroList, error) {
	return ml.InsertListContainer(pos, NewLoroList())
}

// InsertLoroMovableList inserts a fresh attached LoroMovableList at pos and
// returns it.
func (ml *LoroMovableList) InsertLoroMovableList(pos uint32) (*LoroMovableList, error) {
	return ml.InsertMovableListContainer(pos, NewLoroMovableList())
}

// InsertLoroText inserts a fresh attached LoroText at pos and returns it.
func (ml *LoroMovableList) InsertLoroText(pos uint32) (*LoroText, error) {
	return ml.InsertTextContainer(pos, NewLoroText())
}

// InsertLoroTree inserts a fresh attached LoroTree at pos and returns it.
func (ml *LoroMovableList) InsertLoroTree(pos uint32) (*LoroTree, error) {
	return ml.InsertTreeContainer(pos, NewLoroTree())
}

// InsertLoroCounter inserts a fresh attached LoroCounter at pos and returns
// it.
func (ml *LoroMovableList) InsertLoroCounter(pos uint32) (*LoroCounter, error) {
	return ml.InsertCounterContainer(pos, NewLoroCounter())
}
