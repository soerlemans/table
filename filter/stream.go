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

func (this *Stream) Prev() {
	this.Index--
}

// Increment the index.
func (this *Stream) Next() {
	this.Index++
}

// Peek the next rune.
func (this *Stream) Peek() (rune, bool) {
	var rn rune

	this.Next()
	ok := !this.Eos()

	// Only get the rune if its valid.
	if ok {
		rn = this.Current()
	}
	this.Prev()

	return rn, ok
}

// Short for End Of Stream.
func (this *Stream) Eos() bool {
	return this.Index >= len(*this.View)
}

func initStream(t_view *string) Stream {
	return Stream{t_view, 0}
}
