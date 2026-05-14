package loro

import (
	"reflect"
	"testing"
)

func TestLoroMap_TypedGetters(t *testing.T) {
	_, m := newDocWithMap(t)

	if v, ok := m.GetString("s"); !ok || v != "hello" {
		t.Fatalf("GetString = (%q, %v)", v, ok)
	}
	if _, ok := m.GetString("i"); ok {
		t.Fatal("GetString on int should be false")
	}
	if _, ok := m.GetString("missing"); ok {
		t.Fatal("GetString on missing should be false")
	}

	if v, ok := m.GetBool("b"); !ok || !v {
		t.Fatalf("GetBool = (%v, %v)", v, ok)
	}
	if v, ok := m.GetFloat64("f"); !ok || v != 1.5 {
		t.Fatalf("GetFloat64 = (%v, %v)", v, ok)
	}
	if v, ok := m.GetInt64("i"); !ok || v != 42 {
		t.Fatalf("GetInt64 = (%v, %v)", v, ok)
	}

	if l, ok := m.GetList("list"); !ok || len(l) != 2 {
		t.Fatalf("GetList = (%v, %v)", l, ok)
	}
	if l, ok := m.GetListOfAny("list"); !ok || !reflect.DeepEqual(l, []any{"a", int64(1)}) {
		t.Fatalf("GetListOfAny = (%v, %v)", l, ok)
	}
	if mp, ok := m.GetMap("map"); !ok || len(mp) != 1 {
		t.Fatalf("GetMap = (%v, %v)", mp, ok)
	}
	if mp, ok := m.GetMapOfAny("map"); !ok || !reflect.DeepEqual(mp, map[string]any{"k": "v"}) {
		t.Fatalf("GetMapOfAny = (%v, %v)", mp, ok)
	}

	if v, ok := m.GetAny("s"); !ok || v != "hello" {
		t.Fatalf("GetAny(s) = (%v, %v)", v, ok)
	}
	if !m.IsExplicitlyNil("n") {
		t.Fatal("IsExplicitlyNil(n) should be true")
	}
	if m.IsExplicitlyNil("missing") {
		t.Fatal("IsExplicitlyNil(missing) should be false")
	}
}

func TestLoroMap_ContainerGetters(t *testing.T) {
	_, m := newDocWithMap(t)
	if _, err := m.InsertMapContainer("child_map", NewLoroMap()); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertListContainer("child_list", NewLoroList()); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertTextContainer("child_text", NewLoroText()); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertMovableListContainer("child_movable", NewLoroMovableList()); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertCounterContainer("child_counter", NewLoroCounter()); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertTreeContainer("child_tree", NewLoroTree()); err != nil {
		t.Fatal(err)
	}

	if v, ok := m.GetLoroMap("child_map"); !ok || v == nil {
		t.Fatalf("GetLoroMap = (%v, %v)", v, ok)
	}
	if v, ok := m.GetLoroList("child_list"); !ok || v == nil {
		t.Fatalf("GetLoroList = (%v, %v)", v, ok)
	}
	if v, ok := m.GetLoroText("child_text"); !ok || v == nil {
		t.Fatalf("GetLoroText = (%v, %v)", v, ok)
	}
	if v, ok := m.GetLoroMovableList("child_movable"); !ok || v == nil {
		t.Fatalf("GetLoroMovableList = (%v, %v)", v, ok)
	}
	if v, ok := m.GetLoroCounter("child_counter"); !ok || v == nil {
		t.Fatalf("GetLoroCounter = (%v, %v)", v, ok)
	}
	if v, ok := m.GetLoroTree("child_tree"); !ok || v == nil {
		t.Fatalf("GetLoroTree = (%v, %v)", v, ok)
	}

	// Wrong-type queries return ok=false.
	if _, ok := m.GetLoroMap("child_list"); ok {
		t.Fatal("GetLoroMap on a list should be false")
	}
	if _, ok := m.GetLoroList("child_map"); ok {
		t.Fatal("GetLoroList on a map should be false")
	}
	if _, ok := m.GetLoroMap("missing"); ok {
		t.Fatal("GetLoroMap on missing should be false")
	}
}

func TestLoroList_TypedGetters(t *testing.T) {
	doc := NewLoroDoc()
	l := doc.GetList(AsContainerId("l"))
	must(t, l.Push(AsStringValue("s")))
	must(t, l.Push(AsInt64Value(7)))
	must(t, l.Push(AsBoolValue(true)))
	must(t, l.Push(AsFloat64Value(2.5)))
	must(t, l.Push(AsNilValue()))

	if v, ok := l.GetString(0); !ok || v != "s" {
		t.Fatalf("GetString(0) = (%q, %v)", v, ok)
	}
	if _, ok := l.GetString(1); ok {
		t.Fatal("GetString on int index should be false")
	}
	if _, ok := l.GetString(100); ok {
		t.Fatal("GetString on out-of-bounds should be false")
	}

	if v, ok := l.GetInt64(1); !ok || v != 7 {
		t.Fatalf("GetInt64(1) = (%v, %v)", v, ok)
	}
	if v, ok := l.GetBool(2); !ok || !v {
		t.Fatalf("GetBool(2) = (%v, %v)", v, ok)
	}
	if v, ok := l.GetFloat64(3); !ok || v != 2.5 {
		t.Fatalf("GetFloat64(3) = (%v, %v)", v, ok)
	}
	if !l.IsExplicitlyNil(4) {
		t.Fatal("IsExplicitlyNil(4) should be true")
	}
	if l.IsExplicitlyNil(100) {
		t.Fatal("IsExplicitlyNil(oob) should be false")
	}
}

func TestLoroList_ContainerGetters(t *testing.T) {
	doc := NewLoroDoc()
	l := doc.GetList(AsContainerId("l"))
	if _, err := l.InsertMapContainer(0, NewLoroMap()); err != nil {
		t.Fatal(err)
	}
	if _, err := l.InsertTextContainer(1, NewLoroText()); err != nil {
		t.Fatal(err)
	}
	if v, ok := l.GetLoroMap(0); !ok || v == nil {
		t.Fatalf("GetLoroMap(0) = (%v, %v)", v, ok)
	}
	if v, ok := l.GetLoroText(1); !ok || v == nil {
		t.Fatalf("GetLoroText(1) = (%v, %v)", v, ok)
	}
	if _, ok := l.GetLoroMap(1); ok {
		t.Fatal("GetLoroMap on text should be false")
	}
}

func TestLoroMovableList_TypedGetters(t *testing.T) {
	doc := NewLoroDoc()
	ml := doc.GetMovableList(AsContainerId("ml"))
	must(t, ml.Push(AsStringValue("a")))
	must(t, ml.Push(AsInt64Value(99)))

	if v, ok := ml.GetString(0); !ok || v != "a" {
		t.Fatalf("GetString(0) = (%q, %v)", v, ok)
	}
	if v, ok := ml.GetInt64(1); !ok || v != 99 {
		t.Fatalf("GetInt64(1) = (%v, %v)", v, ok)
	}
	if _, ok := ml.GetString(99); ok {
		t.Fatal("GetString OOB should be false")
	}
}
