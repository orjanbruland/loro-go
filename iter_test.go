package loro

import (
	"sort"
	"testing"
)

func TestLoroList_All(t *testing.T) {
	doc := NewLoroDoc()
	l := doc.GetList(AsContainerId("l"))
	must(t, l.Push(AsStringValue("a")))
	must(t, l.Push(AsStringValue("b")))
	must(t, l.Push(AsStringValue("c")))

	var got []string
	for i, v := range l.All() {
		s, ok := GetStringValue(&v)
		if !ok {
			t.Fatalf("element %d not a string", i)
		}
		got = append(got, s)
	}
	if want := []string{"a", "b", "c"}; !equalStringSlice(got, want) {
		t.Fatalf("All = %v, want %v", got, want)
	}

	// Early break stops iteration after the first element.
	var seen int
	for range l.All() {
		seen++
		break
	}
	if seen != 1 {
		t.Fatalf("early break seen = %d, want 1", seen)
	}

	// Empty list yields nothing.
	empty := doc.GetList(AsContainerId("empty"))
	for range empty.All() {
		t.Fatal("empty list should not yield")
	}
}

func TestLoroMovableList_All(t *testing.T) {
	doc := NewLoroDoc()
	ml := doc.GetMovableList(AsContainerId("ml"))
	must(t, ml.Push(AsInt64Value(1)))
	must(t, ml.Push(AsInt64Value(2)))

	var got []int64
	for _, v := range ml.All() {
		n, ok := GetInt64Value(&v)
		if !ok {
			t.Fatal("expected int")
		}
		got = append(got, n)
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("All = %v", got)
	}
}

func TestLoroMap_All(t *testing.T) {
	doc := NewLoroDoc()
	m := doc.GetMap(AsContainerId("m"))
	must(t, m.Insert("a", AsInt64Value(1)))
	must(t, m.Insert("b", AsInt64Value(2)))
	must(t, m.Insert("c", AsInt64Value(3)))

	got := map[string]int64{}
	for k, v := range m.All() {
		n, ok := GetInt64Value(&v)
		if !ok {
			t.Fatalf("key %q not int", k)
		}
		got[k] = n
	}
	want := map[string]int64{"a": 1, "b": 2, "c": 3}
	if len(got) != len(want) {
		t.Fatalf("All = %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("got[%q] = %d, want %d", k, got[k], v)
		}
	}

	// Early break.
	seen := 0
	for range m.All() {
		seen++
		break
	}
	if seen != 1 {
		t.Fatalf("early break seen = %d, want 1", seen)
	}
}

func TestLoroTree_Iterators(t *testing.T) {
	doc := NewLoroDoc()
	tr := doc.GetTree(AsContainerId("tree"))

	root, err := tr.Create(TreeParentRoot)
	if err != nil {
		t.Fatal(err)
	}
	child1, err := tr.Create(TreeParentNode(root))
	if err != nil {
		t.Fatal(err)
	}
	child2, err := tr.Create(TreeParentNode(root))
	if err != nil {
		t.Fatal(err)
	}

	var nodes []TreeId
	for id := range tr.AllNodes() {
		nodes = append(nodes, id)
	}
	if len(nodes) != 3 {
		t.Fatalf("AllNodes len = %d, want 3", len(nodes))
	}

	var roots []TreeId
	for id := range tr.AllRoots() {
		roots = append(roots, id)
	}
	if len(roots) != 1 || roots[0] != root {
		t.Fatalf("AllRoots = %v, want [%v]", roots, root)
	}

	var children []TreeId
	for id := range tr.AllChildren(TreeParentNode(root)) {
		children = append(children, id)
	}
	sort.Slice(children, func(i, j int) bool {
		return treeIdLess(children[i], children[j])
	})
	expected := []TreeId{child1, child2}
	sort.Slice(expected, func(i, j int) bool {
		return treeIdLess(expected[i], expected[j])
	})
	if len(children) != 2 || children[0] != expected[0] || children[1] != expected[1] {
		t.Fatalf("AllChildren = %v, want %v", children, expected)
	}

	// AllChildren on a leaf node yields nothing.
	leafCount := 0
	for range tr.AllChildren(TreeParentNode(child1)) {
		leafCount++
	}
	if leafCount != 0 {
		t.Fatalf("leaf children count = %d, want 0", leafCount)
	}
}

func treeIdLess(a, b TreeId) bool {
	if a.Peer != b.Peer {
		return a.Peer < b.Peer
	}
	return a.Counter < b.Counter
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
