package ir

// TODO: xxx

// Convenience alias.
type InstructionList = []Instruction

type InstructionType int

const (
	InstructionType = iota

	LoadVariable
	LoadField

	LessThan
	LessThanEqual
	Equal
	NotEqual
	GreaterThan
	GreaterThanEqual

	WhenBlock
	MutBlock
	OutBlock
	MdBlock
	JsonBlock
)

type ValueType int

const (
	String ValueType = iota,
	Number,
)

// Interface for interacting with filter nodes.
type Instruction interface {
	//	eval()
	//	print()
}

type BinaryExpr struct {
	lhs Node
	rhs Node
}
