package writer

import (
	td "github.com/soerlemans/table/table_data"
)

type WriterPtr = *Writer
type FmtWriterPtr = *FmtWriter
type RawWriterPtr = *RawWriter

// A writer is an unified interface for writing the results.
// Of the table command in some unique way.
type Writer interface {
	GetLabel() string
	Write() error
}

// A format writer outputs structured rows in a formatted way.
type FmtWriter interface {
	SetHeaders(headers td.TableDataRow)
	AddRow(t_row td.TableDataRow)

	GetLabel() string
	Write() error
}

// A raw writer receives raw data.
type RawWriter interface {
	SetHeaders(t_row string) int
	AddRow(t_row string) int

	GetLabel() string
	Write() error
}
