package filter

// TODO: Consider visitor pattern, maybe later down the line.

type Visitor[V any] interface {
	Visit(visitor V)
}
