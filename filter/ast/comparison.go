package ast

type ComparisonType interface {
	LessThan | LessThanEqual | Equal | NotEqual | GreaterThan | GreaterThanEqual
}

type LessThan struct {
	BinaryExpr
}

type LessThanEqual struct {
	BinaryExpr
}

type Equal struct {
	BinaryExpr
}

type NotEqual struct {
	BinaryExpr
}

type GreaterThan struct {
	BinaryExpr
}

type GreaterThanEqual struct {
	BinaryExpr
}

func InitComparison[T ComparisonType](t_lhs Node, t_rhs Node) T {
	var comp T

	comp = T{
		BinaryExpr: BinaryExpr{t_lhs, t_rhs},
	}

	return comp
}
