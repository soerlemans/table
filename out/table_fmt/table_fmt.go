package table_fmt

import (
	"fmt"
	"sort"
	"strconv"

	s "github.com/soerlemans/table/out/sink"
	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

type TableFmtPtr = *TableFmt

const (
	HeadUnset = -1
	TailUnset = -1
	SortUnset = -1
)

// Format output.
type TableFmt interface {
	GetLabel() string

	// Order related:
	// Determines which columns should be printed and in what order.
	// No defined order results in printing all columns in their default order.
	SetOrder(t_order []int)
	GetOrder() []int
	ClearOrder()

	// Bounds related:
	SetHead(t_count int)
	GetHead() int
	ClearHead()

	SetTail(t_count int)
	GetTail() int
	ClearTail()

	InBounds(t_index int) bool

	// Sorting related:
	SetSort(t_col int)
	GetSort() int
	ClearSort()
	Sort()

	SetNumericSort(t_col int)
	GetNumericSort() int
	ClearNumericSort()
	NumericSort()

	PerformSort()

	// Primary data related:
	SetHeaders(headers td.TableDataRow)
	GetHeaders() td.TableDataRow

	SetRows(t_rows []td.TableDataRow)
	GetRows() []td.TableDataRow

	AddRow(t_row td.TableDataRow)
	RowLen() int

	// Set the sink.
	SetSink(t_sink s.Sink)
	GetSink() s.Sink

	// Copying over:
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

	SortCol        int
	NumericSortCol int

	Head int
	Tail int

	// We need to calculate the max column width for every entry.
	Headers td.TableDataRow
	Rows    []td.TableDataRow

	Sink s.Sink
}

// Private:
func (this *BaseTableFmt) writef(t_fmt string, t_args ...interface{}) {
	this.Sink.Writef(t_fmt, t_args...)
}

func (this *BaseTableFmt) writeln(t_args ...interface{}) {
	this.Sink.Writeln(t_args...)
}

// Public:
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
	this.Head = HeadUnset
}

func (this *BaseTableFmt) SetTail(t_count int) {
	this.Tail = t_count
}

func (this *BaseTableFmt) GetTail() int {
	return this.Tail
}

func (this *BaseTableFmt) ClearTail() {
	this.Tail = TailUnset
}

// Check if a certain row index is in bounds of the head and tail.
func (this *BaseTableFmt) InBounds(t_index int) bool {
	// Default we are always inbounds.
	var (
		inBound     = false
		headIsUnset = false
		tailIsUnset = false
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
	if head > HeadUnset {
		// If the index is below the head count we are in bounds.
		inBound = (t_index < head)
		u.Logf("InBounds: %v = %d < %d", inBound, t_index, head)
		if inBound {
			return inBound
		}
	} else {
		headIsUnset = true
	}

	if tail > TailUnset {
		// We need to subtract one here to account for zero index-ation.
		tailBound := (rowCount - 1) - tail

		// If the index is above the tailBound we are in bounds.
		inBound = (t_index > tailBound)
		u.Logf("InBounds: %v = %d > %d", inBound, t_index, tailBound)
		if inBound {
			return inBound
		}
	} else {
		tailIsUnset = true
	}

	if headIsUnset && tailIsUnset {
		inBound = true
	}

	return inBound
}

// Sorting related:
func (this *BaseTableFmt) SetSort(t_col int) {
	this.SortCol = t_col
}

func (this *BaseTableFmt) GetSort() int {
	return this.SortCol
}

func (this *BaseTableFmt) ClearSort() {
	this.SortCol = SortUnset
}

func (this *BaseTableFmt) Sort() {
	col := this.SortCol

	less := func(index1, index2 int) bool {
		return this.Rows[index1][col] < this.Rows[index2][col]
	}

	// Use sort.Slice with a custom less function
	sort.Slice(this.Rows, less)
}

func (this *BaseTableFmt) SetNumericSort(t_col int) {
	this.NumericSortCol = t_col
}

func (this *BaseTableFmt) GetNumericSort() int {
	return this.NumericSortCol
}

func (this *BaseTableFmt) ClearNumericSort() {
	this.NumericSortCol = SortUnset
}

func (this *BaseTableFmt) NumericSort() {
	col := this.NumericSortCol

	less := func(index1, index2 int) bool {
		// Convert string to int for the column to sort.
		num1, err := strconv.Atoi(this.Rows[index1][col])
		if err != nil {
			// If not a num then we want it at the end.
			// Then we want the not num at the end.
			return false
		}

		num2, err := strconv.Atoi(this.Rows[index2][col])
		if err != nil {
			// If only the second arg is not a num.
			// Then we want the not num at the end.
			return true
		}

		return num1 < num2
	}

	// Use sort.Slice with a custom less function
	sort.Slice(this.Rows, less)
}

func (this *BaseTableFmt) PerformSort() {
	// Sort and numerical sort are mutually exclusive operations.
	if this.NumericSortCol != SortUnset {
		this.NumericSort()
	} else if this.SortCol != SortUnset {
		this.Sort()
	}
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

func (this *BaseTableFmt) SetSink(t_sink s.Sink) {
	this.Sink = t_sink
}

func (this *BaseTableFmt) GetSink() s.Sink {
	return this.Sink
}

// Generic copying functionality.
func (this *BaseTableFmt) Copy(t_fmt TableFmt) error {
	// This is a generic implementation for copying over all data.
	this.Label = t_fmt.GetLabel()

	headers := t_fmt.GetHeaders()
	this.SetHeaders(headers)

	rows := t_fmt.GetRows()
	this.SetRows(rows)

	order := t_fmt.GetOrder()
	this.SetOrder(order)

	sort_ := t_fmt.GetSort()
	this.SetSort(sort_)

	numSort := t_fmt.GetNumericSort()
	this.SetNumericSort(numSort)

	head := t_fmt.GetHead()
	this.SetHead(head)

	tail := t_fmt.GetTail()
	this.SetTail(tail)

	sink := t_fmt.GetSink()
	this.SetSink(sink)

	return nil
}

func InitBaseTableFmt(t_label string) (BaseTableFmt, error) {
	format := BaseTableFmt{}

	format.Label = t_label
	format.SortCol = SortUnset
	format.NumericSortCol = SortUnset
	format.Head = HeadUnset
	format.Tail = TailUnset

	sink := s.InitStdoutSink()
	format.Sink = &sink

	return format, nil
}
