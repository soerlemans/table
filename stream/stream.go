package stream

// TODO: Make use of generics instead.

// type Indexable[T any] interface {
// 	~string | ~[]T
// }

type Stream[V any] interface {
	// TODO: Consider:
	// Current() (V, error)

	Current() V
	Prev()
	Next()

	Peek() (V, bool)
	Append(V)

	Eos() bool
	Len() int
}
