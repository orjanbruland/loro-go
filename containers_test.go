package loro

import (
	"strings"
	"testing"
)

func TestLoroMap_Lookup(t *testing.T) {
	_, m := newDocWithMap(t)
	if v := m.Lookup("s"); v == nil {
		t.Fatal("Lookup(s) should not be nil")
	}
	if v := m.Lookup("missing"); v != nil {
		t.Fatal("Lookup(missing) should be nil")
	}
}

func TestLoroList_At(t *testing.T) {
	doc := NewLoroDoc()
	l := doc.GetList(AsContainerId("l"))
	must(t, l.Push(AsStringValue("x")))

	if v := l.At(0); v == nil {
		t.Fatal("At(0) should not be nil")
	}
	if v := l.At(100); v != nil {
		t.Fatal("At(oob) should be nil")
	}
}

func TestLoroMovableList_At(t *testing.T) {
	doc := NewLoroDoc()
	ml := doc.GetMovableList(AsContainerId("ml"))
	must(t, ml.Push(AsStringValue("x")))

	if v := ml.At(0); v == nil {
		t.Fatal("At(0) should not be nil")
	}
	if v := ml.At(100); v != nil {
		t.Fatal("At(oob) should be nil")
	}
}

func TestLoroDoc_FindContainers(t *testing.T) {
	doc := NewLoroDoc()
	_ = doc.GetMap(AsContainerId("m"))
	_ = doc.GetList(AsContainerId("l"))
	_ = doc.GetText(AsContainerId("t"))
	_ = doc.GetMovableList(AsContainerId("ml"))
	_ = doc.GetCounter(AsContainerId("c"))
	_ = doc.GetTree(AsContainerId("tr"))

	if v := doc.FindMap(AsContainerId("m")); v == nil {
		t.Fatal("FindMap should resolve a created map")
	}
	if v := doc.FindList(AsContainerId("l")); v == nil {
		t.Fatal("FindList should resolve a created list")
	}
	if v := doc.FindText(AsContainerId("t")); v == nil {
		t.Fatal("FindText should resolve a created text")
	}
	if v := doc.FindMovableList(AsContainerId("ml")); v == nil {
		t.Fatal("FindMovableList should resolve a created movable list")
	}
	if v := doc.FindCounter(AsContainerId("c")); v == nil {
		t.Fatal("FindCounter should resolve a created counter")
	}
	if v := doc.FindTree(AsContainerId("tr")); v == nil {
		t.Fatal("FindTree should resolve a created tree")
	}
}

func TestLoroDoc_FindByPath(t *testing.T) {
	doc := NewLoroDoc()
	m := doc.GetMap(AsContainerId("m"))
	must(t, m.InsertAny("k", "v"))
	doc.Commit()

	if v := doc.FindByPath("m", "k"); v == nil {
		t.Fatal("FindByPath should resolve m.k")
	}
	if v := doc.FindByPath("missing"); v != nil {
		t.Fatal("FindByPath on missing should be nil")
	}
}

func TestLoroDoc_FindByPath_mixedTypes(t *testing.T) {
	doc := NewLoroDoc()

	root := doc.GetMap(AsContainerId("root"))
	users, err := root.InsertListContainer("users", NewLoroList())
	if err != nil {
		t.Fatalf("insert users list: %v", err)
	}
	alice, err := users.InsertMapContainer(0, NewLoroMap())
	if err != nil {
		t.Fatalf("insert alice map: %v", err)
	}
	must(t, alice.Insert("name", AsStringValue("Alice")))

	v := doc.FindByPath("root", "users", 0, "name")
	if v == nil {
		t.Fatal("FindByPath returned nil for present nested path")
	}
	got, ok := GetStringValue(&v)
	if !ok || got != "Alice" {
		t.Fatalf("got (%q, %v), want (Alice, true)", got, ok)
	}
}

func TestLoroDoc_FindByPath_intWidths(t *testing.T) {
	doc := NewLoroDoc()
	root := doc.GetMap(AsContainerId("root"))
	list, err := root.InsertListContainer("items", NewLoroList())
	if err != nil {
		t.Fatalf("insert list: %v", err)
	}
	must(t, list.Push(AsStringValue("first")))

	cases := []struct {
		name string
		idx  any
	}{
		{"int", int(0)},
		{"int8", int8(0)},
		{"int16", int16(0)},
		{"int32", int32(0)},
		{"int64", int64(0)},
		{"uint", uint(0)},
		{"uint8", uint8(0)},
		{"uint16", uint16(0)},
		{"uint32", uint32(0)},
		{"uint64", uint64(0)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := doc.FindByPath("root", "items", c.idx)
			if v == nil {
				t.Fatalf("nil for %s index", c.name)
			}
			got, ok := GetStringValue(&v)
			if !ok || got != "first" {
				t.Fatalf("got (%q, %v), want (first, true)", got, ok)
			}
		})
	}
}

func TestLoroDoc_FindByPath_passthroughIndex(t *testing.T) {
	doc := NewLoroDoc()
	root := doc.GetMap(AsContainerId("root"))
	must(t, root.Insert("k", AsStringValue("v")))

	v := doc.FindByPath(PathKey("root"), PathKey("k"))
	if v == nil {
		t.Fatal("FindByPath returned nil for passthrough Index args")
	}
	got, ok := GetStringValue(&v)
	if !ok || got != "v" {
		t.Fatalf("got (%q, %v), want (v, true)", got, ok)
	}
}

func TestLoroDoc_FindByPath_unsupportedTypePanics(t *testing.T) {
	doc := NewLoroDoc()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for unsupported path part type")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value is not a string: %#v", r)
		}
		if !strings.Contains(msg, "unsupported path part type") {
			t.Fatalf("panic message %q missing expected text", msg)
		}
	}()
	doc.FindByPath("root", 1.5) // float64 isn't supported
}

func TestPathConstructors(t *testing.T) {
	if got, want := PathKey("k"), (IndexKey{Key: "k"}); got != want {
		t.Fatalf("PathKey = %#v, want %#v", got, want)
	}
	if got, want := PathSeq(3), (IndexSeq{Index: 3}); got != want {
		t.Fatalf("PathSeq = %#v, want %#v", got, want)
	}
	id := TreeId{Peer: 1, Counter: 2}
	if got, want := PathNode(id), (IndexNode{Target: id}); got != want {
		t.Fatalf("PathNode = %#v, want %#v", got, want)
	}
}

func TestContainer_OwnerDoc(t *testing.T) {
	doc := NewLoroDoc()
	m := doc.GetMap(AsContainerId("m"))
	if owner := m.OwnerDoc(); owner == nil {
		t.Fatal("attached map should have an OwnerDoc")
	}

	detached := NewLoroMap()
	if owner := detached.OwnerDoc(); owner != nil {
		t.Fatal("detached map OwnerDoc should be nil")
	}
}
