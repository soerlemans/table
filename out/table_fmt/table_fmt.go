package table_fmt

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

type TableFmtPtr = *TableFmt

const (
	HEAD_UNSET = -1
	TAIL_UNSET = -1
)

// Format output.
type TableFmt interface {
	GetLabel() string

	// Determines which columns should be printed and in what order.
	// No defined order results in printing all columns in their default order.
	SetOrder(t_order []int)
	GetOrder() []int
	ClearOrder()

	SetHead(t_count int)
	GetHead() int
	ClearHead()

	SetTail(t_count int)
	GetTail() int
	ClearTail()

	InBounds(t_index int) bool

	SetHeaders(headers td.TableDataRow)
	GetHeaders() td.TableDataRow
	SetRows(t_rows []td.TableDataRow)
	GetRows() []td.TableDataRow

	AddRow(t_row td.TableDataRow)
	RowLen() int

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

	Head int
	Tail int

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

func (this *BaseTableFmt) SetHead(t_count int) {
	this.Head = t_count
}

func (this *BaseTableFmt) GetHead() int {
	return this.Head
}

func (this *BaseTableFmt) ClearHead() {
	this.Head = HEAD_UNSET
}

func (this *BaseTableFmt) SetTail(t_count int) {
	this.Tail = t_count
}

func (this *BaseTableFmt) GetTail() int {
	return this.Tail
}

func (this *BaseTableFmt) ClearTail() {
	this.Tail = TAIL_UNSET
}

// Check if a certain row index is in bounds of the head and tail.
func (this *BaseTableFmt) InBounds(t_index int) bool {
	// Default we are always inbounds.
	var (
		inBound = false
		headInBound = false
		tailInBound = false
	)

	head := this.GetHead()
	tail := this.GetTail()

	rowCount := this.RowLen()

	if t_index < 0 {
		msg := fmt.Sprintf("Index error (index:%d < 0)!", t_index)
		panic(msg)
	}

	if t_index >= rowCount {
		msg := fmt.Sprintf("Index error (index:%d >= rowCount:%d)!", t_index, rowCount)
		panic(msg)
	}

	// Any negative values are seen as being unset.
	if head > HEAD_UNSET {
		// If the index is below the head count we are in bounds.
		headInBound = (t_index < head)
		u.Logf("InBounds: %v = %d < %d", headInBound, t_index, head)
	} else {
		headInBound = true
	}

	if tail > TAIL_UNSET {
		// We need to subtract one here to account for zero index-ation.
		tailBound := (rowCount - 1) - tail

		// If the index is above the tailBound we are in bounds.
		tailInBound = (t_index > tailBound)
		u.Logf("InBounds: %v = %d > %d", tailInBound, t_index, tailBound)
	} else {
		tailInBound = true
	}

	if headInBound || tailInBound {
		inBound = true
	}

	return inBound
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

func (this *BaseTableFmt) RowLen() int {
	return len(this.Rows)
}

// Generic copying functionality.
func (this *BaseTableFmt) Copy(t_fmt TableFmt) error {
	this.Label = t_fmt.GetLabel()

	headers := t_fmt.GetHeaders()
	this.SetHeaders(headers)

	rows := t_fmt.GetRows()
	this.SetRows(rows)

	order := t_fmt.GetOrder()
	this.SetOrder(order)

	head := t_fmt.GetHead()
	this.SetHead(head)

	tail := t_fmt.GetTail()
	this.SetTail(tail)

	return nil
}
