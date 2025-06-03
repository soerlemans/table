package writer

import (
	td "github.com/soerlemans/table/table_data"
)

type WriterPtr = *Writer

// A format writer outputs structured rows in a formatted way.
type Writer interface {
	SetHeaders(headers td.TableDataRow)
	AddRow(t_row td.TableDataRow)

	GetLabel() string
	Write() error
}
