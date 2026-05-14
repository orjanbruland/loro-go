package loro

import "testing"

func TestTreeParentSingletons(t *testing.T) {
	if _, ok := TreeParentRoot.(TreeParentIdRoot); !ok {
		t.Fatalf("TreeParentRoot = %T, want TreeParentIdRoot", TreeParentRoot)
	}
	if _, ok := TreeParentDeleted.(TreeParentIdDeleted); !ok {
		t.Fatalf("TreeParentDeleted = %T, want TreeParentIdDeleted", TreeParentDeleted)
	}
	if _, ok := TreeParentUnexist.(TreeParentIdUnexist); !ok {
		t.Fatalf("TreeParentUnexist = %T, want TreeParentIdUnexist", TreeParentUnexist)
	}
}

func TestTreeParentNode(t *testing.T) {
	id := TreeId{Peer: 42, Counter: 7}
	got := TreeParentNode(id)
	node, ok := got.(TreeParentIdNode)
	if !ok {
		t.Fatalf("TreeParentNode returned %T, want TreeParentIdNode", got)
	}
	if node.Id != id {
		t.Fatalf("TreeParentNode.Id = %v, want %v", node.Id, id)
	}
}

func TestTreeParentRoot_TreeCreate(t *testing.T) {
	doc := NewLoroDoc()
	tr := doc.GetTree(AsContainerId("tree"))

	root, err := tr.Create(TreeParentRoot)
	if err != nil {
		t.Fatal(err)
	}
	child, err := tr.Create(TreeParentNode(root))
	if err != nil {
		t.Fatal(err)
	}

	parent, err := tr.Parent(child)
	if err != nil {
		t.Fatal(err)
	}
	node, ok := parent.(TreeParentIdNode)
	if !ok {
		t.Fatalf("Parent(child) = %T, want TreeParentIdNode", parent)
	}
	if node.Id != root {
		t.Fatalf("Parent(child).Id = %v, want %v", node.Id, root)
	}
}
