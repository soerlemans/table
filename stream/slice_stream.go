package stream

type SliceStream[T any] struct {
	View  *[]T
	Index int
}

// Get current item.
func (this *SliceStream[T]) Current() T {
	view := *(this.View)

	return view[this.Index]
}

// Decrement the index.
func (this *SliceStream[T]) Prev() {
	this.Index--
}

// Increment the index.
func (this *SliceStream[T]) Next() {
	this.Index++
}

// Peek the next rune.
func (this *SliceStream[T]) Peek() (T, bool) {
	var val T

	this.Next()
	ok := !this.Eos()

	// Only get the value if its valid.
	if ok {
		val = this.Current()
	}
	this.Prev()

	return val, ok
}

// Append to the stream.
func (this *SliceStream[T]) Append(t_value T) {
	*this.View = append(*this.View, t_value)
}

// Short for End Of SliceStream[T].
func (this *SliceStream[T]) Eos() bool {
	return this.Index >= len(*this.View)
}

func (this *SliceStream[T]) Len() int {
	return len(*this.View)
}

func InitSliceStream[T any](t_view *[]T) SliceStream[T] {
	return SliceStream[T]{t_view, 0}
}

func InitSliceStreamEmpty[T any]() SliceStream[T] {
	// Account for if the stream no elems yet.
	placeholder := []T{}

	return SliceStream[T]{&placeholder, 0}
}
