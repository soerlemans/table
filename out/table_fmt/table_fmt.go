package table_fmt

import (
	td "github.com/soerlemans/table/table_data"
)

type TableFmtPtr = *TableFmt

// Format output.
type TableFmt interface {
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

// Embeddable struct for reducing boilerplate.
type StandardTableFmt struct {
	Label string

	// We need to calculate the max column width for every entry.
	ColWidth map[int]int
	Headers  td.TableDataRow

	ColMask map[int]bool

	Rows []td.TableDataRow
}

func (this *StandardTableFmt) GetLabel() string {
	return this.Label
}

func (this *StandardTableFmt) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers
}

// Mark columns to print during write.
func (this *StandardTableFmt) SetMask(t_mask []int) {
	this.ClearMask()

	// TODO: Error handle non existent column indices (for now ignore).
	for _, value := range t_mask {
		this.ColMask[value] = true
	}
}

func (this *StandardTableFmt) ClearMask() {
	// Clear by assigning a new one.
	this.ColMask = make(map[int]bool)
}

func (this *StandardTableFmt) ColumnMasked(t_colIndex int) bool {
	// Guard clause (the mask has no elements then print everything).
	// As we should always print atleast one column.
	if len(this.ColMask) == 0 {
		return true
	}

	// We use the map as a set.
	_, ok := this.ColMask[t_colIndex]

	return ok
}

func (this *StandardTableFmt) SetRows(t_rows []td.TableDataRow) {
	this.Rows = t_rows
}

func (this *StandardTableFmt) GetRows() []td.TableDataRow {
	return this.Rows
}

func (this *StandardTableFmt) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)
}
