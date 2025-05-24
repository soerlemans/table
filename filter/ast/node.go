package ast

// Convenience alias.
type NodeList = []Node

// Interface for interacting with filter nodes.
type Node interface {
	//	eval()
	//	print()
}

type BinaryExpr struct {
	lhs Node
	rhs Node
}
