package loro

// Typed accessor methods on LoroMap / LoroList / LoroMovableList that read a
// value at a key/index and assert its type. Each returns the zero value and
// false when the entry is absent or holds a different type.
//
// These wrap the free functions in values.go so callsites don't have to pass
// the double-pointer `**ValueOrContainer` around.

// GetString returns the string value at key, or ("", false) if the key is
// absent or the value is not a string.
func (m *LoroMap) GetString(key string) (string, bool) {
	return GetStringValue(m.Get(key))
}

// GetBool returns the bool value at key, or (false, false) if the key is
// absent or the value is not a bool.
func (m *LoroMap) GetBool(key string) (bool, bool) {
	return GetBoolValue(m.Get(key))
}

// GetFloat64 returns the float64 value at key, or (0, false) if the key is
// absent or the value is not a float64.
func (m *LoroMap) GetFloat64(key string) (float64, bool) {
	return GetFloat64Value(m.Get(key))
}

// GetInt64 returns the int64 value at key, or (0, false) if the key is absent
// or the value is not an int64.
func (m *LoroMap) GetInt64(key string) (int64, bool) {
	return GetInt64Value(m.Get(key))
}

// GetList returns the list value at key as []LoroValue, or (nil, false) if
// the key is absent or the value is not a list.
func (m *LoroMap) GetList(key string) ([]LoroValue, bool) {
	return GetListValue(m.Get(key))
}

// GetListOfAny returns the list value at key as []any, or (nil, false) if
// the key is absent or the value is not a list.
func (m *LoroMap) GetListOfAny(key string) ([]any, bool) {
	return GetListValueOfAny(m.Get(key))
}

// GetMap returns the map value at key as map[string]LoroValue, or (nil, false)
// if the key is absent or the value is not a map.
func (m *LoroMap) GetMap(key string) (map[string]LoroValue, bool) {
	return GetMapValue(m.Get(key))
}

// GetMapOfAny returns the map value at key as map[string]any, or (nil, false)
// if the key is absent or the value is not a map.
func (m *LoroMap) GetMapOfAny(key string) (map[string]any, bool) {
	return GetMapValueOfAny(m.Get(key))
}

// GetAny returns the value at key as an any, or (nil, false) if the key is
// absent. Nested containers are returned as their *LoroX wrapper.
func (m *LoroMap) GetAny(key string) (any, bool) {
	return GetAnyValue(m.Get(key))
}

// IsExplicitlyNil reports whether the value at key is set to an explicit nil
// (LoroValueNull), as opposed to the key being absent.
func (m *LoroMap) IsExplicitlyNil(key string) bool {
	return IsValueExplicitlyNil(m.Get(key))
}

// GetLoroCounter returns the LoroCounter at key, or (nil, false) if the key
// is absent or the value is not a counter container.
func (m *LoroMap) GetLoroCounter(key string) (*LoroCounter, bool) {
	return GetLoroCounterContainer(m.Get(key))
}

// GetLoroMap returns the LoroMap at key, or (nil, false) if the key is absent
// or the value is not a map container.
func (m *LoroMap) GetLoroMap(key string) (*LoroMap, bool) {
	return GetLoroMapContainer(m.Get(key))
}

// GetLoroList returns the LoroList at key, or (nil, false) if the key is
// absent or the value is not a list container.
func (m *LoroMap) GetLoroList(key string) (*LoroList, bool) {
	return GetLoroListContainer(m.Get(key))
}

// GetLoroMovableList returns the LoroMovableList at key, or (nil, false) if
// the key is absent or the value is not a movable list container.
func (m *LoroMap) GetLoroMovableList(key string) (*LoroMovableList, bool) {
	return GetLoroMovableListContainer(m.Get(key))
}

// GetLoroText returns the LoroText at key, or (nil, false) if the key is
// absent or the value is not a text container.
func (m *LoroMap) GetLoroText(key string) (*LoroText, bool) {
	return GetLoroTextContainer(m.Get(key))
}

// GetLoroTree returns the LoroTree at key, or (nil, false) if the key is
// absent or the value is not a tree container.
func (m *LoroMap) GetLoroTree(key string) (*LoroTree, bool) {
	return GetLoroTreeContainer(m.Get(key))
}

// GetString returns the string value at index, or ("", false) if the index is
// out of bounds or the value is not a string.
func (l *LoroList) GetString(index uint32) (string, bool) {
	return GetStringValue(l.Get(index))
}

// GetBool returns the bool value at index, or (false, false) if the index is
// out of bounds or the value is not a bool.
func (l *LoroList) GetBool(index uint32) (bool, bool) {
	return GetBoolValue(l.Get(index))
}

// GetFloat64 returns the float64 value at index, or (0, false) if the index
// is out of bounds or the value is not a float64.
func (l *LoroList) GetFloat64(index uint32) (float64, bool) {
	return GetFloat64Value(l.Get(index))
}

// GetInt64 returns the int64 value at index, or (0, false) if the index is
// out of bounds or the value is not an int64.
func (l *LoroList) GetInt64(index uint32) (int64, bool) {
	return GetInt64Value(l.Get(index))
}

// GetList returns the list value at index as []LoroValue, or (nil, false) if
// the index is out of bounds or the value is not a list.
func (l *LoroList) GetList(index uint32) ([]LoroValue, bool) {
	return GetListValue(l.Get(index))
}

// GetListOfAny returns the list value at index as []any, or (nil, false) if
// the index is out of bounds or the value is not a list.
func (l *LoroList) GetListOfAny(index uint32) ([]any, bool) {
	return GetListValueOfAny(l.Get(index))
}

// GetMap returns the map value at index as map[string]LoroValue, or
// (nil, false) if the index is out of bounds or the value is not a map.
func (l *LoroList) GetMap(index uint32) (map[string]LoroValue, bool) {
	return GetMapValue(l.Get(index))
}

// GetMapOfAny returns the map value at index as map[string]any, or
// (nil, false) if the index is out of bounds or the value is not a map.
func (l *LoroList) GetMapOfAny(index uint32) (map[string]any, bool) {
	return GetMapValueOfAny(l.Get(index))
}

// GetAny returns the value at index as an any, or (nil, false) if the index
// is out of bounds. Nested containers are returned as their *LoroX wrapper.
func (l *LoroList) GetAny(index uint32) (any, bool) {
	return GetAnyValue(l.Get(index))
}

// IsExplicitlyNil reports whether the value at index is set to an explicit
// nil (LoroValueNull), as opposed to the index being out of bounds.
func (l *LoroList) IsExplicitlyNil(index uint32) bool {
	return IsValueExplicitlyNil(l.Get(index))
}

// GetLoroCounter returns the LoroCounter at index, or (nil, false) if the
// index is out of bounds or the value is not a counter container.
func (l *LoroList) GetLoroCounter(index uint32) (*LoroCounter, bool) {
	return GetLoroCounterContainer(l.Get(index))
}

// GetLoroMap returns the LoroMap at index, or (nil, false) if the index is
// out of bounds or the value is not a map container.
func (l *LoroList) GetLoroMap(index uint32) (*LoroMap, bool) {
	return GetLoroMapContainer(l.Get(index))
}

// GetLoroList returns the LoroList at index, or (nil, false) if the index is
// out of bounds or the value is not a list container.
func (l *LoroList) GetLoroList(index uint32) (*LoroList, bool) {
	return GetLoroListContainer(l.Get(index))
}

// GetLoroMovableList returns the LoroMovableList at index, or (nil, false) if
// the index is out of bounds or the value is not a movable list container.
func (l *LoroList) GetLoroMovableList(index uint32) (*LoroMovableList, bool) {
	return GetLoroMovableListContainer(l.Get(index))
}

// GetLoroText returns the LoroText at index, or (nil, false) if the index is
// out of bounds or the value is not a text container.
func (l *LoroList) GetLoroText(index uint32) (*LoroText, bool) {
	return GetLoroTextContainer(l.Get(index))
}

// GetLoroTree returns the LoroTree at index, or (nil, false) if the index is
// out of bounds or the value is not a tree container.
func (l *LoroList) GetLoroTree(index uint32) (*LoroTree, bool) {
	return GetLoroTreeContainer(l.Get(index))
}

// GetString returns the string value at index, or ("", false) if the index is
// out of bounds or the value is not a string.
func (ml *LoroMovableList) GetString(index uint32) (string, bool) {
	return GetStringValue(ml.Get(index))
}

// GetBool returns the bool value at index, or (false, false) if the index is
// out of bounds or the value is not a bool.
func (ml *LoroMovableList) GetBool(index uint32) (bool, bool) {
	return GetBoolValue(ml.Get(index))
}

// GetFloat64 returns the float64 value at index, or (0, false) if the index
// is out of bounds or the value is not a float64.
func (ml *LoroMovableList) GetFloat64(index uint32) (float64, bool) {
	return GetFloat64Value(ml.Get(index))
}

// GetInt64 returns the int64 value at index, or (0, false) if the index is
// out of bounds or the value is not an int64.
func (ml *LoroMovableList) GetInt64(index uint32) (int64, bool) {
	return GetInt64Value(ml.Get(index))
}

// GetList returns the list value at index as []LoroValue, or (nil, false) if
// the index is out of bounds or the value is not a list.
func (ml *LoroMovableList) GetList(index uint32) ([]LoroValue, bool) {
	return GetListValue(ml.Get(index))
}

// GetListOfAny returns the list value at index as []any, or (nil, false) if
// the index is out of bounds or the value is not a list.
func (ml *LoroMovableList) GetListOfAny(index uint32) ([]any, bool) {
	return GetListValueOfAny(ml.Get(index))
}

// GetMap returns the map value at index as map[string]LoroValue, or
// (nil, false) if the index is out of bounds or the value is not a map.
func (ml *LoroMovableList) GetMap(index uint32) (map[string]LoroValue, bool) {
	return GetMapValue(ml.Get(index))
}

// GetMapOfAny returns the map value at index as map[string]any, or
// (nil, false) if the index is out of bounds or the value is not a map.
func (ml *LoroMovableList) GetMapOfAny(index uint32) (map[string]any, bool) {
	return GetMapValueOfAny(ml.Get(index))
}

// GetAny returns the value at index as an any, or (nil, false) if the index
// is out of bounds. Nested containers are returned as their *LoroX wrapper.
func (ml *LoroMovableList) GetAny(index uint32) (any, bool) {
	return GetAnyValue(ml.Get(index))
}

// IsExplicitlyNil reports whether the value at index is set to an explicit
// nil (LoroValueNull), as opposed to the index being out of bounds.
func (ml *LoroMovableList) IsExplicitlyNil(index uint32) bool {
	return IsValueExplicitlyNil(ml.Get(index))
}

// GetLoroCounter returns the LoroCounter at index, or (nil, false) if the
// index is out of bounds or the value is not a counter container.
func (ml *LoroMovableList) GetLoroCounter(index uint32) (*LoroCounter, bool) {
	return GetLoroCounterContainer(ml.Get(index))
}

// GetLoroMap returns the LoroMap at index, or (nil, false) if the index is
// out of bounds or the value is not a map container.
func (ml *LoroMovableList) GetLoroMap(index uint32) (*LoroMap, bool) {
	return GetLoroMapContainer(ml.Get(index))
}

// GetLoroList returns the LoroList at index, or (nil, false) if the index is
// out of bounds or the value is not a list container.
func (ml *LoroMovableList) GetLoroList(index uint32) (*LoroList, bool) {
	return GetLoroListContainer(ml.Get(index))
}

// GetLoroMovableList returns the LoroMovableList at index, or (nil, false) if
// the index is out of bounds or the value is not a movable list container.
func (ml *LoroMovableList) GetLoroMovableList(index uint32) (*LoroMovableList, bool) {
	return GetLoroMovableListContainer(ml.Get(index))
}

// GetLoroText returns the LoroText at index, or (nil, false) if the index is
// out of bounds or the value is not a text container.
func (ml *LoroMovableList) GetLoroText(index uint32) (*LoroText, bool) {
	return GetLoroTextContainer(ml.Get(index))
}

// GetLoroTree returns the LoroTree at index, or (nil, false) if the index is
// out of bounds or the value is not a tree container.
func (ml *LoroMovableList) GetLoroTree(index uint32) (*LoroTree, bool) {
	return GetLoroTreeContainer(ml.Get(index))
}
