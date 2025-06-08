package sink

import (
	"fmt"
)

// Very straightforward implementation.
type StdoutSink struct {
}

func (this *StdoutSink) Writef(t_fmt string, t_args ...interface{}) {
	fmt.Printf(t_fmt, t_args...)
}

func (this *StdoutSink) Writeln(t_args ...interface{}) {
	fmt.Println(t_args...)
}

func InitStdoutSink() StdoutSink {
	return StdoutSink{}
}
