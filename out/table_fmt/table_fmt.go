package table_fmt

import (
	td "github.com/soerlemans/table/table_data"
)

type TableFmtPtr = *TableFmt

// Format output.
type TableFmt interface {
	GetLabel() string

	// Determines which columns should be printed and in what order.
	// No defined order results in printing all columns in their default order.
	SetOrder(t_order []int)
	GetOrder() []int
	ClearOrder()

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

	// Determines order of the columns, as well as which columns to print.
	// If empty will print all columns in their regular format.
	Order []int

	// We need to calculate the max column width for every entry.
	Headers td.TableDataRow
	Rows    []td.TableDataRow
}

func (this *BaseTableFmt) GetLabel() string {
	return this.Label
}

// Mark columns to print during write.
func (this *BaseTableFmt) SetOrder(t_order []int) {
	// The mask also determined the order.
	this.Order = t_order
}

func (this *BaseTableFmt) GetOrder() []int {
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

func (this *BaseTableFmt) ClearOrder() {
	// Clear the order and column mask.
	this.Order = nil
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

	mask := t_fmt.GetOrder()
	this.SetOrder(mask)

	return nil
}
