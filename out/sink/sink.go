package sink

import (
// td "github.com/soerlemans/table/table_data"
)

// A format writer outputs structured rows in a formatted way.
type Sink interface {
	Writef(t_fmt string, t_args ...interface{}) error
	Writeln(t_args ...interface{}) error
}
