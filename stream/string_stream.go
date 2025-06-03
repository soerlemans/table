package stream

type StringStream struct {
	View  *string
	Index int
}

func (this *StringStream) Current() rune {
	view := *(this.View)

	return rune(view[this.Index])
}

func (this *StringStream) Prev() {
	this.Index--
}

// Increment the index.
func (this *StringStream) Next() {
	this.Index++
}

// Peek the next rune.
func (this *StringStream) Peek() (rune, bool) {
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

// Append to the stream.
func (this *StringStream) Append(t_value string) {
	*this.View = *this.View + t_value
}

// Short for End Of StringStream.
func (this *StringStream) Eos() bool {
	return this.Index >= len(*this.View)
}

func (this *StringStream) Len() int {
	return len(*this.View)
}

func InitStringStream(t_view *string) StringStream {
	return StringStream{t_view, 0}
}
