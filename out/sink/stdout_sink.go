package sink

import (
	"fmt"
)

// Very straightforward implementation.
type StdoutSink struct {
}

func (this *StdoutSink) Writef(t_fmt string, t_args ...interface{}) error {
	fmt.Printf(t_fmt, t_args...)

	return nil
}

func (this *StdoutSink) Writeln(t_args ...interface{}) error {
	fmt.Println(t_args...)

	return nil
}

func InitStdoutSink() StdoutSink {
	return StdoutSink{}
}
