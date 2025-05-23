package filter_node

type ComparisonType interface {
	LessThan | LessThanEqual | Equal | NotEqual | GreaterThan | GreaterThanEqual
}

type LessThan struct{}
type LessThanEqual struct{}

type Equal struct{}
type NotEqual struct{}

type GreaterThan struct{}
// func (this *GreaterThan) eval() {}

type GreaterThanEqual struct{}

// func (this *GreaterThanEqual) eval() {}

func InitComparison[T ComparisonType](t_lhs NodePtr, t_rhs NodePtr) T {
	return T{}
}
