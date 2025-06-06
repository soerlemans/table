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
	ColumnOrder() []int

	// SetHead(t_count int)
	// SetTail(t_count int)

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

	// Order and the column mask are tied closely together.
	Order   []int
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

	// The mask also determined the order.
	this.Order = t_mask

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
	// Clear the order and column mask.
	this.Order = nil
	this.ColMask = make(map[int]bool)
}

func (this *BaseTableFmt) ColumnOrder() []int {
	order := this.Order

	// If no order is set, just use the normal ordering.
	if order == nil {
		order = make([]int, len(this.Headers))
		for i := range this.Headers {
			order[i] = i
		}
	}

	return order
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
