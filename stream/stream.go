package stream

// TODO: Make use of generics instead.

// type Indexable[T any] interface {
// 	~string | ~[]T
// }

type Stream[V any] interface {
	Current() V
	Prev()
	Next()
	Peek() (V, bool)
	Eos() bool
	Append(V)
}
