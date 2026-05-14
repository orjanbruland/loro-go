package loro

import "testing"

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

	if v := doc.FindByPath([]Index{IndexKey{Key: "m"}, IndexKey{Key: "k"}}); v == nil {
		t.Fatal("FindByPath should resolve m.k")
	}
	if v := doc.FindByPath([]Index{IndexKey{Key: "missing"}}); v != nil {
		t.Fatal("FindByPath on missing should be nil")
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
