package filter_node

// These are pointers, garbage collector keeps them in scope.
type NodeListPtr []*FilterNode
type NodePtr *FilterNode

// Interface for interacting with filter nodes.
type FilterNode interface {
	// TODO: Implement.

	//	eval()
	//	print()
}
