package table_fmt

import (
	td "github.com/soerlemans/table/table_data"
)

type TableFmtPtr = *TableFmt

// Format output.
type TableFmt interface {
	GetLabel() string

	// Determine which columns to mask, this determines if they will be printed.
	SetMask(t_filter []int)
	GetMask() []int

	ClearMask()
	ColumnMasked(t_colIndex int) bool

	SetHeaders(headers td.TableDataRow)
	GetHeaders() td.TableDataRow
	SetRows(t_rows []td.TableDataRow)
	GetRows() []td.TableDataRow

	AddRow(t_row td.TableDataRow)
	Copy(t_fmt TableFmt) error

	Write() error
}

// Embeddable struct for reducing boilerplate.
// Does not implement Write().
type BaseTableFmt struct {
	Label string

	ColMask map[int]bool

	// We need to calculate the max column width for every entry.
	Headers td.TableDataRow
	Rows    []td.TableDataRow
}

func (this *BaseTableFmt) GetLabel() string {
	return this.Label
}

// Mark columns to print during write.
func (this *BaseTableFmt) SetMask(t_mask []int) {
	this.ClearMask()

	// TODO: Error handle non existent column indices (for now ignore).
	for _, value := range t_mask {
		this.ColMask[value] = true
	}
}

func (this *BaseTableFmt) GetMask() []int {
	var mask []int

	for key, value := range this.ColMask {
		if value {
			mask = append(mask, key)
		}
	}

	return mask
}

func (this *BaseTableFmt) ClearMask() {
	// Clear by assigning a new one.
	this.ColMask = make(map[int]bool)
}

func (this *BaseTableFmt) ColumnMasked(t_colIndex int) bool {
	// Guard clause (the mask has no elements then print everything).
	// As we should always print atleast one column.
	if len(this.ColMask) == 0 {
		return true
	}

	// We use the map as a set.
	_, ok := this.ColMask[t_colIndex]

	return ok
}

func (this *BaseTableFmt) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers
}

func (this *BaseTableFmt) GetHeaders() td.TableDataRow {
	return this.Headers
}

func (this *BaseTableFmt) SetRows(t_rows []td.TableDataRow) {
	this.Rows = t_rows
}

func (this *BaseTableFmt) GetRows() []td.TableDataRow {
	return this.Rows
}

func (this *BaseTableFmt) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)
}

// Generic copying functionality.
func (this *BaseTableFmt) Copy(t_fmt TableFmt) error {
	this.Label = t_fmt.GetLabel()

	headers := t_fmt.GetHeaders()
	this.SetHeaders(headers)

	rows := t_fmt.GetRows()
	this.SetRows(rows)

	mask := t_fmt.GetMask()
	this.SetMask(mask)

	return nil
}
