package loro

import (
	"fmt"
	"math"
)

// Deref unwraps a double pointer returned by some generated bindings methods.
// Returns nil if v is nil or *v is nil; otherwise returns *v.
//
// Useful as an escape hatch for the few `**T` returns not covered by typed
// wrapper methods (e.g. LoroDoc.FrontiersToVv).
func Deref[T any](v **T) *T {
	if v == nil {
		return nil
	}
	return *v
}

type simpleValue struct {
	value LoroValue
}

func (v simpleValue) AsLoroValue() LoroValue {
	return v.value
}

// AsStringValue converts a string to a LoroValueLike which can be used in most
// places that expect a value.
func AsStringValue(v string) LoroValueLike {
	return simpleValue{
		value: LoroValueString{
			Value: v,
		},
	}
}

// AsNilValue converts a null value to a LoroValueLike which can be used in most
// places that expect a value.
func AsNilValue() LoroValueLike {
	return simpleValue{
		value: LoroValueNull{},
	}
}

// AsBoolValue converts a bool to a LoroValueLike which can be used in most
// places that expect a value.
func AsBoolValue(v bool) LoroValueLike {
	return simpleValue{
		value: LoroValueBool{
			Value: v,
		},
	}
}

// AsFloat64Value converts a float64 to a LoroValueLike which can be used in most
// places that expect a value.
func AsFloat64Value(v float64) LoroValueLike {
	return simpleValue{
		value: LoroValueDouble{
			Value: v,
		},
	}
}

// AsInt64Value converts an int64 to a LoroValueLike which can be used in most
// places that expect a value.
func AsInt64Value(v int64) LoroValueLike {
	return simpleValue{
		value: LoroValueI64{
			Value: v,
		},
	}
}

// AsBinaryValue converts a byte slice to a LoroValueLike which can be used in
// most places that expect a value.
func AsBinaryValue(v []byte) LoroValueLike {
	return simpleValue{
		value: LoroValueBinary{
			Value: v,
		},
	}
}

// AsListValue converts a list of LoroValueLike to a LoroValueLike which can be
// used in most places that expect a value.
func AsListValue(v []LoroValueLike) LoroValueLike {
	values := make([]LoroValue, len(v))
	for i, value := range v {
		values[i] = value.AsLoroValue()
	}
	return simpleValue{
		value: LoroValueList{
			Value: values,
		},
	}
}

// AsListValueFromValue converts a vector of LoroValue to a LoroValueLike which
// can be used in most places that expect a value.
func AsListValueFromValue(v []LoroValue) LoroValueLike {
	return simpleValue{
		value: LoroValueList{
			Value: v,
		},
	}
}

// AsListValueFromAny converts a vector of any to a LoroValueLike which
// can be used in most places that expect a value.
func AsListValueFromAny(v []any) (LoroValueLike, error) {
	values := make([]LoroValue, len(v))
	for i, value := range v {
		value, err := AsValue(value)
		if err != nil {
			return nil, err
		}

		values[i] = value.AsLoroValue()
	}

	return AsListValueFromValue(values), nil
}

// AsMapValue converts a map of string to LoroValueLike to a LoroValueLike which
// can be used in most places that expect a value.
func AsMapValue(v map[string]LoroValueLike) LoroValueLike {
	values := make(map[string]LoroValue)
	for k, v := range v {
		values[k] = v.AsLoroValue()
	}
	return simpleValue{
		value: LoroValueMap{
			Value: values,
		},
	}
}

// AsMapValueFromValue converts a map of string to LoroValue to a LoroValueLike
// which can be used in most places that expect a value.
func AsMapValueFromValue(v map[string]LoroValue) LoroValueLike {
	return simpleValue{
		value: LoroValueMap{
			Value: v,
		},
	}
}

// AsMapValueFromAny converts a map of string to any to a LoroValueLike
// which can be used in most places that expect a value.
func AsMapValueFromAny(v map[string]any) (LoroValueLike, error) {
	values := make(map[string]LoroValue)

	for k, v := range v {
		value, err := AsValue(v)
		if err != nil {
			return nil, err
		}
		values[k] = value.AsLoroValue()
	}

	return AsMapValueFromValue(values), nil
}

// AsValue converts an any to a LoroValueLike which can be used in most
// places that expect a value.
func AsValue(v any) (LoroValueLike, error) {
	if v == nil {
		return AsNilValue(), nil
	}

	switch v := v.(type) {
	case LoroValueLike:
		return v, nil
	case LoroValue:
		return simpleValue{value: v}, nil
	case string:
		return AsStringValue(v), nil
	case bool:
		return AsBoolValue(v), nil
	case float64:
		return AsFloat64Value(v), nil
	case float32:
		return AsFloat64Value(float64(v)), nil
	case int64:
		return AsInt64Value(v), nil
	case int:
		return AsInt64Value(int64(v)), nil
	case int32:
		return AsInt64Value(int64(v)), nil
	case int16:
		return AsInt64Value(int64(v)), nil
	case int8:
		return AsInt64Value(int64(v)), nil
	case uint64:
		if v > math.MaxInt64 {
			return nil, fmt.Errorf("uint64 value %d exceeds math.MaxInt64", v)
		}
		return AsInt64Value(int64(v)), nil
	case uint:
		if uint64(v) > math.MaxInt64 {
			return nil, fmt.Errorf("uint value %d exceeds math.MaxInt64", v)
		}
		return AsInt64Value(int64(v)), nil
	case uint32:
		return AsInt64Value(int64(v)), nil
	case uint16:
		return AsInt64Value(int64(v)), nil
	case uint8:
		return AsInt64Value(int64(v)), nil
	case []byte:
		return AsBinaryValue(v), nil
	case []any:
		return AsListValueFromAny(v)
	case map[string]any:
		return AsMapValueFromAny(v)
	}

	return nil, fmt.Errorf("unsupported value type: %T", v)
}

// GetValue takes ValueOrContainer and returns the LoroValue if it is a value.
func GetValue(v **ValueOrContainer) (LoroValue, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroValue := v0.AsValue()
	if loroValue == nil {
		return nil, false
	}

	return *loroValue, true
}

// IsValueExplicitlyNil checks if a pointer to a ValueOrContainer is explicitly
// nil, meaning it has a LoroValueNull value.
func IsValueExplicitlyNil(v **ValueOrContainer) bool {
	value, ok := GetValue(v)
	if !ok {
		return false
	}

	if _, ok := value.(LoroValueNull); ok {
		return true
	}

	return false
}

// GetStringValue takes a pointer to a ValueOrContainer and returns the string
// value if it is a string.
func GetStringValue(v **ValueOrContainer) (string, bool) {
	value, ok := GetValue(v)
	if !ok {
		return "", false
	}
	return ValueAsString(value)
}

// GetBoolValue takes a pointer to a ValueOrContainer and returns the bool
// value if it is a bool.
func GetBoolValue(v **ValueOrContainer) (bool, bool) {
	value, ok := GetValue(v)
	if !ok {
		return false, false
	}
	return ValueAsBool(value)
}

// GetFloat64Value takes a pointer to a ValueOrContainer and returns the float64
// value if it is a float64.
func GetFloat64Value(v **ValueOrContainer) (float64, bool) {
	value, ok := GetValue(v)
	if !ok {
		return 0, false
	}
	return ValueAsFloat64(value)
}

// GetInt64Value takes a pointer to a ValueOrContainer and returns the int64
// value if it is an int64.
func GetInt64Value(v **ValueOrContainer) (int64, bool) {
	value, ok := GetValue(v)
	if !ok {
		return 0, false
	}
	return ValueAsInt64(value)
}

// GetListValue takes a pointer to a ValueOrContainer and returns the list
// value if it is a list.
func GetListValue(v **ValueOrContainer) ([]LoroValue, bool) {
	value, ok := GetValue(v)
	if !ok {
		return nil, false
	}
	return ValueAsList(value)
}

// GetListValueOfAny takes a pointer to a ValueOrContainer and returns a slice
// of any if it is a list.
func GetListValueOfAny(v **ValueOrContainer) ([]any, bool) {
	list, ok := GetListValue(v)
	if !ok {
		return nil, false
	}

	values := make([]any, len(list))
	for i, v := range list {
		v, ok := ValueAsAny(v)
		if !ok {
			return nil, false
		}
		values[i] = v
	}
	return values, true
}

// GetMapValue takes a pointer to a ValueOrContainer and returns the map
// value if it is a map.
func GetMapValue(v **ValueOrContainer) (map[string]LoroValue, bool) {
	value, ok := GetValue(v)
	if !ok {
		return nil, false
	}
	return ValueAsMap(value)
}

// GetMapValueOfAny takes a pointer to a ValueOrContainer and returns a map
// of string to any if it is a map.
func GetMapValueOfAny(v **ValueOrContainer) (map[string]any, bool) {
	m, ok := GetMapValue(v)
	if !ok {
		return nil, false
	}

	values := make(map[string]any, len(m))
	for k, v := range m {
		v, ok := ValueAsAny(v)
		if !ok {
			return nil, false
		}
		values[k] = v
	}
	return values, true
}

// GetAnyValue takes a pointer to a ValueOrContainer and returns the value
// as an any.
func GetAnyValue(v **ValueOrContainer) (any, bool) {
	if v == nil {
		return nil, false
	}

	value, ok := GetValue(v)
	if ok {
		return ValueAsAny(value)
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	if loroCounter := v0.AsLoroCounter(); loroCounter != nil {
		return *loroCounter, true
	}

	if loroMap := v0.AsLoroMap(); loroMap != nil {
		return *loroMap, true
	}

	if loroList := v0.AsLoroList(); loroList != nil {
		return *loroList, true
	}

	if loroMovableList := v0.AsLoroMovableList(); loroMovableList != nil {
		return *loroMovableList, true
	}

	if loroText := v0.AsLoroText(); loroText != nil {
		return *loroText, true
	}

	if loroTree := v0.AsLoroTree(); loroTree != nil {
		return *loroTree, true
	}

	return nil, false
}

// ValueIsNil reports whether the LoroValue is an explicit null.
func ValueIsNil(value LoroValue) bool {
	_, ok := value.(LoroValueNull)
	return ok
}

// ValueAsString returns the string contents of a LoroValue if it holds a string.
func ValueAsString(value LoroValue) (string, bool) {
	if s, ok := value.(LoroValueString); ok {
		return s.Value, true
	}
	return "", false
}

// ValueAsBool returns the bool contents of a LoroValue if it holds a bool.
func ValueAsBool(value LoroValue) (bool, bool) {
	if b, ok := value.(LoroValueBool); ok {
		return b.Value, true
	}
	return false, false
}

// ValueAsInt64 returns the int64 contents of a LoroValue if it holds an int64.
// No implicit conversion from float is performed.
func ValueAsInt64(value LoroValue) (int64, bool) {
	if i, ok := value.(LoroValueI64); ok {
		return i.Value, true
	}
	return 0, false
}

// ValueAsFloat64 returns the float64 contents of a LoroValue if it holds a
// float. No implicit conversion from int is performed.
func ValueAsFloat64(value LoroValue) (float64, bool) {
	if f, ok := value.(LoroValueDouble); ok {
		return f.Value, true
	}
	return 0, false
}

// ValueAsBinary returns the byte slice contents of a LoroValue if it holds
// binary data.
func ValueAsBinary(value LoroValue) ([]byte, bool) {
	if b, ok := value.(LoroValueBinary); ok {
		return b.Value, true
	}
	return nil, false
}

// ValueAsList returns the list contents of a LoroValue if it holds a list.
func ValueAsList(value LoroValue) ([]LoroValue, bool) {
	if l, ok := value.(LoroValueList); ok {
		return l.Value, true
	}
	return nil, false
}

// ValueAsMap returns the map contents of a LoroValue if it holds a map.
func ValueAsMap(value LoroValue) (map[string]LoroValue, bool) {
	if m, ok := value.(LoroValueMap); ok {
		return m.Value, true
	}
	return nil, false
}

// ValueAsAny converts a LoroValue to its idiomatic Go representation,
// recursing into lists and maps. Returns false if the value is of an
// unsupported variant.
func ValueAsAny(value LoroValue) (any, bool) {
	switch value := value.(type) {
	case LoroValueNull:
		return nil, true
	case LoroValueString:
		return value.Value, true
	case LoroValueBool:
		return value.Value, true
	case LoroValueDouble:
		return value.Value, true
	case LoroValueI64:
		return value.Value, true
	case LoroValueBinary:
		return value.Value, true
	case LoroValueList:
		values := make([]any, len(value.Value))
		for i, v := range value.Value {
			v, ok := ValueAsAny(v)
			if !ok {
				return nil, false
			}
			values[i] = v
		}
		return values, true
	case LoroValueMap:
		values := make(map[string]any, len(value.Value))
		for k, v := range value.Value {
			v, ok := ValueAsAny(v)
			if !ok {
				return nil, false
			}
			values[k] = v
		}
		return values, true
	}

	return nil, false
}

// GetLoroCounterContainer takes a pointer to a ValueOrContainer and returns the
// LoroCounter value if it is a LoroCounter.
func GetLoroCounterContainer(v **ValueOrContainer) (*LoroCounter, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroCounter := v0.AsLoroCounter()
	if loroCounter == nil {
		return nil, false
	}

	loroCounter0 := *loroCounter
	if loroCounter0 == nil {
		return nil, false
	}

	return loroCounter0, true
}

// GetLoroMapContainer takes a pointer to a ValueOrContainer and returns the
// LoroMap value if it is a LoroMap.
func GetLoroMapContainer(v **ValueOrContainer) (*LoroMap, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroMap := v0.AsLoroMap()
	if loroMap == nil {
		return nil, false
	}

	loroMap0 := *loroMap
	if loroMap0 == nil {
		return nil, false
	}

	return loroMap0, true
}

// GetLoroListContainer takes a pointer to a ValueOrContainer and returns the
// LoroList value if it is a LoroList.
func GetLoroListContainer(v **ValueOrContainer) (*LoroList, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroList := v0.AsLoroList()
	if loroList == nil {
		return nil, false
	}

	loroList0 := *loroList
	if loroList0 == nil {
		return nil, false
	}

	return loroList0, true
}

// GetLoroMovableListContainer takes a pointer to a ValueOrContainer and returns the
// LoroMovableList value if it is a LoroMovableList.
func GetLoroMovableListContainer(v **ValueOrContainer) (*LoroMovableList, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroMovableList := v0.AsLoroMovableList()
	if loroMovableList == nil {
		return nil, false
	}

	loroMovableList0 := *loroMovableList
	if loroMovableList0 == nil {
		return nil, false
	}

	return loroMovableList0, true
}

// GetLoroTextContainer takes a pointer to a ValueOrContainer and returns the
// LoroText value if it is a LoroText.
func GetLoroTextContainer(v **ValueOrContainer) (*LoroText, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroText := v0.AsLoroText()
	if loroText == nil {
		return nil, false
	}

	loroText0 := *loroText
	if loroText0 == nil {
		return nil, false
	}

	return loroText0, true
}

// GetLoroTreeContainer takes a pointer to a ValueOrContainer and returns the
// LoroTree value if it is a LoroTree.
func GetLoroTreeContainer(v **ValueOrContainer) (*LoroTree, bool) {
	if v == nil {
		return nil, false
	}

	v0 := *v
	if v0 == nil {
		return nil, false
	}

	loroTree := v0.AsLoroTree()
	if loroTree == nil {
		return nil, false
	}

	loroTree0 := *loroTree
	if loroTree0 == nil {
		return nil, false
	}

	return loroTree0, true
}
