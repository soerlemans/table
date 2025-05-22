package filter

// TODO: Make use of generics instead.

type Stream struct {
	View  *string
	Index int
}

func (this *Stream) Current() rune {
	view := *(this.View)
	ch := view[this.Index]

	return rune(ch)
}

// Increment the index.
func (this *Stream) Next() {
	this.Index++
}

// Short for End Of Stream.
func (this *Stream) Eos() bool {
	return this.Index >= len(*this.View)
}

func initStream(t_view *string) Stream {
	return Stream{t_view, 0}
}
