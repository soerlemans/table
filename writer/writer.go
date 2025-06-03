package writer

import (
	td "github.com/soerlemans/table/table_data"
)

type WriterPtr = *Writer

// A format writer outputs structured rows in a formatted way.
type Writer interface {
	GetLabel() string

	SetHeaders(headers td.TableDataRow)

	// Determine which columns to mask, this determines if they will be printed.
	SetMask(t_filter []int)
	ClearMask()
	ColumnMasked(t_colIndex int) bool

	SetRows(t_rows []td.TableDataRow)
	GetRows() []td.TableDataRow

	AddRow(t_row td.TableDataRow)

	Write() error
}
