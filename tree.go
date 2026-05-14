package loro

// TreeParentId singletons. Use these instead of the zero-value struct literals
// (e.g. TreeParentIdRoot{}) when calling tree APIs that take a TreeParentId.
//
//	root, err := tree.Create(loro.TreeParentRoot)
//	child, err := tree.Create(loro.TreeParentNode(root))
//
// TreeParentRoot is the conceptual root of a tree (a node has no parent).
// It is distinct from TreeRoot(name), which returns a ContainerId for a tree
// root container in the document.
var (
	TreeParentRoot    TreeParentId = TreeParentIdRoot{}
	TreeParentDeleted TreeParentId = TreeParentIdDeleted{}
	TreeParentUnexist TreeParentId = TreeParentIdUnexist{}
)

// TreeParentNode returns a TreeParentId that points at the given tree node.
// It replaces the verbose TreeParentIdNode{Id: id} struct literal.
//
//	child, err := tree.Create(loro.TreeParentNode(parent))
func TreeParentNode(id TreeId) TreeParentId {
	return TreeParentIdNode{Id: id}
}
