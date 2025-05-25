package writer

import (
	td "github.com/soerlemans/table/table_data"
)

// A writer is an unified interface for writing the output of.
type Writer interface {
	Add(t_row td.TableDataRow)
	Write()
}
