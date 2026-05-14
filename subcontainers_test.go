package loro

import "testing"

func TestLoroMap_GetOrCreateLoroChildren(t *testing.T) {
	doc := NewLoroDoc()
	root := doc.GetMap(AsContainerId("root"))

	cm, err := root.GetOrCreateLoroMap("m")
	if err != nil {
		t.Fatalf("GetOrCreateLoroMap: %v", err)
	}
	if !cm.IsAttached() {
		t.Fatal("child map should be attached")
	}

	cl, err := root.GetOrCreateLoroList("l")
	if err != nil {
		t.Fatalf("GetOrCreateLoroList: %v", err)
	}
	if !cl.IsAttached() {
		t.Fatal("child list should be attached")
	}

	cml, err := root.GetOrCreateLoroMovableList("ml")
	if err != nil {
		t.Fatalf("GetOrCreateLoroMovableList: %v", err)
	}
	if !cml.IsAttached() {
		t.Fatal("child movable list should be attached")
	}

	ct, err := root.GetOrCreateLoroText("t")
	if err != nil {
		t.Fatalf("GetOrCreateLoroText: %v", err)
	}
	if !ct.IsAttached() {
		t.Fatal("child text should be attached")
	}

	ctr, err := root.GetOrCreateLoroTree("tr")
	if err != nil {
		t.Fatalf("GetOrCreateLoroTree: %v", err)
	}
	if !ctr.IsAttached() {
		t.Fatal("child tree should be attached")
	}

	cc, err := root.GetOrCreateLoroCounter("c")
	if err != nil {
		t.Fatalf("GetOrCreateLoroCounter: %v", err)
	}
	if !cc.IsAttached() {
		t.Fatal("child counter should be attached")
	}
}

func TestLoroMap_GetOrCreateLoroMap_idempotent(t *testing.T) {
	doc := NewLoroDoc()
	root := doc.GetMap(AsContainerId("root"))

	first, err := root.GetOrCreateLoroMap("m")
	if err != nil {
		t.Fatalf("first GetOrCreateLoroMap: %v", err)
	}
	must(t, first.InsertAny("k", "v"))

	second, err := root.GetOrCreateLoroMap("m")
	if err != nil {
		t.Fatalf("second GetOrCreateLoroMap: %v", err)
	}
	if got, ok := second.GetString("k"); !ok || got != "v" {
		t.Fatalf("second call should yield the existing child, got (%q, %v)", got, ok)
	}
}

func TestLoroList_InsertLoroChildren(t *testing.T) {
	doc := NewLoroDoc()
	l := doc.GetList(AsContainerId("l"))

	cm, err := l.InsertLoroMap(0)
	if err != nil {
		t.Fatalf("InsertLoroMap: %v", err)
	}
	if !cm.IsAttached() {
		t.Fatal("inserted map should be attached")
	}

	cl, err := l.InsertLoroList(1)
	if err != nil {
		t.Fatalf("InsertLoroList: %v", err)
	}
	if !cl.IsAttached() {
		t.Fatal("inserted list should be attached")
	}

	cml, err := l.InsertLoroMovableList(2)
	if err != nil {
		t.Fatalf("InsertLoroMovableList: %v", err)
	}
	if !cml.IsAttached() {
		t.Fatal("inserted movable list should be attached")
	}

	ct, err := l.InsertLoroText(3)
	if err != nil {
		t.Fatalf("InsertLoroText: %v", err)
	}
	if !ct.IsAttached() {
		t.Fatal("inserted text should be attached")
	}

	ctr, err := l.InsertLoroTree(4)
	if err != nil {
		t.Fatalf("InsertLoroTree: %v", err)
	}
	if !ctr.IsAttached() {
		t.Fatal("inserted tree should be attached")
	}

	cc, err := l.InsertLoroCounter(5)
	if err != nil {
		t.Fatalf("InsertLoroCounter: %v", err)
	}
	if !cc.IsAttached() {
		t.Fatal("inserted counter should be attached")
	}
}

func TestLoroMovableList_InsertLoroChildren(t *testing.T) {
	doc := NewLoroDoc()
	ml := doc.GetMovableList(AsContainerId("ml"))

	cm, err := ml.InsertLoroMap(0)
	if err != nil {
		t.Fatalf("InsertLoroMap: %v", err)
	}
	if !cm.IsAttached() {
		t.Fatal("inserted map should be attached")
	}

	cl, err := ml.InsertLoroList(1)
	if err != nil {
		t.Fatalf("InsertLoroList: %v", err)
	}
	if !cl.IsAttached() {
		t.Fatal("inserted list should be attached")
	}

	cml, err := ml.InsertLoroMovableList(2)
	if err != nil {
		t.Fatalf("InsertLoroMovableList: %v", err)
	}
	if !cml.IsAttached() {
		t.Fatal("inserted movable list should be attached")
	}

	ct, err := ml.InsertLoroText(3)
	if err != nil {
		t.Fatalf("InsertLoroText: %v", err)
	}
	if !ct.IsAttached() {
		t.Fatal("inserted text should be attached")
	}

	ctr, err := ml.InsertLoroTree(4)
	if err != nil {
		t.Fatalf("InsertLoroTree: %v", err)
	}
	if !ctr.IsAttached() {
		t.Fatal("inserted tree should be attached")
	}

	cc, err := ml.InsertLoroCounter(5)
	if err != nil {
		t.Fatalf("InsertLoroCounter: %v", err)
	}
	if !cc.IsAttached() {
		t.Fatal("inserted counter should be attached")
	}
}
