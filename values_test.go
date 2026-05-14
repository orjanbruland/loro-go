package loro

import (
	"math"
	"reflect"
	"testing"
)

func TestAsValue_Conversions(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want LoroValue
	}{
		{"nil", nil, LoroValueNull{}},
		{"string", "hi", LoroValueString{Value: "hi"}},
		{"bool", true, LoroValueBool{Value: true}},
		{"float64", 1.5, LoroValueDouble{Value: 1.5}},
		{"float32", float32(2.5), LoroValueDouble{Value: 2.5}},
		{"int", 42, LoroValueI64{Value: 42}},
		{"int64", int64(-7), LoroValueI64{Value: -7}},
		{"int32", int32(3), LoroValueI64{Value: 3}},
		{"int16", int16(4), LoroValueI64{Value: 4}},
		{"int8", int8(5), LoroValueI64{Value: 5}},
		{"uint64-ok", uint64(123), LoroValueI64{Value: 123}},
		{"uint", uint(9), LoroValueI64{Value: 9}},
		{"uint32", uint32(10), LoroValueI64{Value: 10}},
		{"uint16", uint16(11), LoroValueI64{Value: 11}},
		{"uint8", uint8(12), LoroValueI64{Value: 12}},
		{"bytes", []byte{1, 2, 3}, LoroValueBinary{Value: []byte{1, 2, 3}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := AsValue(tc.in)
			if err != nil {
				t.Fatalf("AsValue(%v) error: %v", tc.in, err)
			}
			if !reflect.DeepEqual(got.AsLoroValue(), tc.want) {
				t.Fatalf("AsValue(%v) = %#v, want %#v", tc.in, got.AsLoroValue(), tc.want)
			}
		})
	}
}

func TestAsValue_UintOverflow(t *testing.T) {
	if _, err := AsValue(uint64(math.MaxUint64)); err == nil {
		t.Fatalf("expected overflow error for uint64 MaxUint64")
	}
}

func TestAsValue_PassthroughLoroValue(t *testing.T) {
	v := LoroValueString{Value: "x"}
	got, err := AsValue(v)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.AsLoroValue(), v) {
		t.Fatalf("LoroValue passthrough: got %#v, want %#v", got.AsLoroValue(), v)
	}

	like := AsStringValue("y")
	got, err = AsValue(like)
	if err != nil {
		t.Fatal(err)
	}
	if got != like {
		t.Fatal("LoroValueLike passthrough should return same value")
	}
}

func TestAsValue_ListAndMap(t *testing.T) {
	got, err := AsValue([]any{"a", int64(1), true})
	if err != nil {
		t.Fatal(err)
	}
	list, ok := got.AsLoroValue().(LoroValueList)
	if !ok {
		t.Fatalf("expected LoroValueList, got %T", got.AsLoroValue())
	}
	if len(list.Value) != 3 {
		t.Fatalf("len = %d, want 3", len(list.Value))
	}

	got, err = AsValue(map[string]any{"k": "v"})
	if err != nil {
		t.Fatal(err)
	}
	m, ok := got.AsLoroValue().(LoroValueMap)
	if !ok {
		t.Fatalf("expected LoroValueMap, got %T", got.AsLoroValue())
	}
	if !reflect.DeepEqual(m.Value["k"], LoroValueString{Value: "v"}) {
		t.Fatalf("map[k] = %#v", m.Value["k"])
	}
}

func TestAsValue_Unsupported(t *testing.T) {
	if _, err := AsValue(struct{}{}); err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestDeref(t *testing.T) {
	if Deref[int](nil) != nil {
		t.Fatal("Deref(nil) should be nil")
	}
	var inner *int
	if Deref(&inner) != nil {
		t.Fatal("Deref(&nil) should be nil")
	}
	n := 7
	innerPtr := &n
	if got := Deref(&innerPtr); got == nil || *got != 7 {
		t.Fatalf("Deref(&&7) = %v", got)
	}
}

// newDocWithMap returns a doc and a root map populated with a small set of
// values for accessor/getter tests.
func newDocWithMap(t *testing.T) (*LoroDoc, *LoroMap) {
	t.Helper()
	doc := NewLoroDoc()
	m := doc.GetMap(AsContainerId("m"))
	must(t, m.Insert("s", AsStringValue("hello")))
	must(t, m.Insert("b", AsBoolValue(true)))
	must(t, m.Insert("f", AsFloat64Value(1.5)))
	must(t, m.Insert("i", AsInt64Value(42)))
	must(t, m.Insert("n", AsNilValue()))
	must(t, m.Insert("list", AsListValue([]LoroValueLike{AsStringValue("a"), AsInt64Value(1)})))
	must(t, m.Insert("map", AsMapValue(map[string]LoroValueLike{"k": AsStringValue("v")})))
	return doc, m
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetValue_Typed(t *testing.T) {
	_, m := newDocWithMap(t)

	if v, ok := GetStringValue(m.Get("s")); !ok || v != "hello" {
		t.Fatalf("GetStringValue = (%q, %v)", v, ok)
	}
	if _, ok := GetStringValue(m.Get("i")); ok {
		t.Fatal("GetStringValue on int should be false")
	}
	if _, ok := GetStringValue(m.Get("missing")); ok {
		t.Fatal("GetStringValue on missing should be false")
	}

	if v, ok := GetBoolValue(m.Get("b")); !ok || v != true {
		t.Fatalf("GetBoolValue = (%v, %v)", v, ok)
	}
	if _, ok := GetBoolValue(m.Get("s")); ok {
		t.Fatal("GetBoolValue on string should be false")
	}

	if v, ok := GetFloat64Value(m.Get("f")); !ok || v != 1.5 {
		t.Fatalf("GetFloat64Value = (%v, %v)", v, ok)
	}
	if _, ok := GetFloat64Value(m.Get("i")); ok {
		t.Fatal("GetFloat64Value on int should be false (no implicit conversion)")
	}

	if v, ok := GetInt64Value(m.Get("i")); !ok || v != 42 {
		t.Fatalf("GetInt64Value = (%v, %v)", v, ok)
	}

	if list, ok := GetListValue(m.Get("list")); !ok || len(list) != 2 {
		t.Fatalf("GetListValue = (%v, %v)", list, ok)
	}
	if mp, ok := GetMapValue(m.Get("map")); !ok || len(mp) != 1 {
		t.Fatalf("GetMapValue = (%v, %v)", mp, ok)
	}
}

func TestGetListValueOfAny(t *testing.T) {
	_, m := newDocWithMap(t)
	values, ok := GetListValueOfAny(m.Get("list"))
	if !ok {
		t.Fatal("GetListValueOfAny ok = false")
	}
	want := []any{"a", int64(1)}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("got %#v want %#v", values, want)
	}
}

func TestGetMapValueOfAny(t *testing.T) {
	_, m := newDocWithMap(t)
	values, ok := GetMapValueOfAny(m.Get("map"))
	if !ok {
		t.Fatal("GetMapValueOfAny ok = false")
	}
	if !reflect.DeepEqual(values, map[string]any{"k": "v"}) {
		t.Fatalf("got %#v", values)
	}
}

func TestIsValueExplicitlyNil(t *testing.T) {
	_, m := newDocWithMap(t)
	if !IsValueExplicitlyNil(m.Get("n")) {
		t.Fatal("expected explicit nil for key n")
	}
	if IsValueExplicitlyNil(m.Get("missing")) {
		t.Fatal("missing key should not be explicit nil")
	}
	if IsValueExplicitlyNil(m.Get("s")) {
		t.Fatal("string value should not be explicit nil")
	}
}

func TestGetAnyValue_BasicAndContainer(t *testing.T) {
	doc, m := newDocWithMap(t)

	if v, ok := GetAnyValue(m.Get("s")); !ok || v != "hello" {
		t.Fatalf("GetAnyValue(string) = (%v, %v)", v, ok)
	}
	if v, ok := GetAnyValue(m.Get("n")); !ok || v != nil {
		t.Fatalf("GetAnyValue(null) = (%v, %v)", v, ok)
	}

	// Container value: insert a child map and check the wrapper is returned.
	_, err := m.InsertMapContainer("child", NewLoroMap())
	if err != nil {
		t.Fatal(err)
	}
	v, ok := GetAnyValue(m.Get("child"))
	if !ok {
		t.Fatal("GetAnyValue(child) ok = false")
	}
	if _, ok := v.(*LoroMap); !ok {
		t.Fatalf("GetAnyValue(child) returned %T, want *LoroMap", v)
	}

	if _, ok := GetAnyValue(m.Get("missing")); ok {
		t.Fatal("GetAnyValue(missing) should be false")
	}

	_ = doc
}
