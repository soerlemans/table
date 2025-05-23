package stream

type SliceStream[T any] struct {
	View  *[]T
	Index int
}

func (this *SliceStream[T]) Current() T {
	view := *(this.View)

	return view[this.Index]
}

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

	// Only get the rune if its valid.
	if ok {
		val = this.Current()
	}
	this.Prev()

	return val, ok
}

// Short for End Of SliceStream[T].
func (this *SliceStream[T]) Eos() bool {
	return this.Index >= len(*this.View)
}

func InitSliceStream[T any](t_view *[]T) SliceStream[T] {
	return SliceStream[T]{t_view, 0}
}
