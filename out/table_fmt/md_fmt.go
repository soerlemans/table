package table_fmt

import (
	"errors"
	"fmt"
	"strings"

	td "github.com/soerlemans/table/table_data"
)

const (
	MdColSep = '|'
	MdRowSep = '-'
)

var (
	ErrColWidthIndex = errors.New("Column width index does not exist.")
)

type MdFmt struct {
	// Includes base data and methods.
	BaseTableFmt

	ColWidth map[int]int
}

func (this *MdFmt) errColWidthIndex(t_index int) error {
	return fmt.Errorf("No %d index in ColWidth %+v (%w)", t_index, this.ColWidth, ErrColWidthIndex)
}

func (this *MdFmt) updateColWidth(t_row td.TableDataRow) {
	for i, cell := range t_row {
		cellWidth := len(cell)

		// Update the max column size.
		if cellWidth > this.ColWidth[i] {
			this.ColWidth[i] = cellWidth
		}
	}
}

func (this *MdFmt) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers

	this.updateColWidth(t_headers)
}

func (this *MdFmt) SetRows(t_rows []td.TableDataRow) {
	this.Rows = t_rows

	// We need to update the mex column width for every line now.
	for _, row := range this.Rows {
		this.updateColWidth(row)
	}
}

func (this *MdFmt) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)

	this.updateColWidth(t_row)
}

func (this *MdFmt) printRow(t_row td.TableDataRow) error {
	order := this.GetOrder()

	for _, index := range order {
		cell := t_row[index]

		colWidth, ok := this.ColWidth[index]
		if !ok {
			return this.errColWidthIndex(index)
		}

		// Check if the column is selected.
		fmt.Printf("| %-*s ", colWidth, cell)
	}
	fmt.Println("|")

	return nil
}

func (this *MdFmt) printTableHeader() error {
	return this.printRow(this.Headers)
}

func (this *MdFmt) printTableHeaderSep() error {
	order := this.GetOrder()

	for _, index := range order {
		colWidth, ok := this.ColWidth[index]
		if !ok {
			return this.errColWidthIndex(index)
		}

		// Check if the column is selected.
		colSep := strings.Repeat("-", colWidth)
		fmt.Printf("| %s ", colSep)
	}
	fmt.Println("|")

	return nil
}

func (this *MdFmt) printTableRows() error {

	// Print per row.
	for index, row := range this.Rows {
		// Skip if we are not in bounds.
		if !this.InBounds(index) {
			continue
		}

		// Print cells of the row.
		err := this.printRow(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *MdFmt) Write() error {
	this.PerformSort()

	err := this.printTableHeader()
	if err != nil {
		return err
	}

	err = this.printTableHeaderSep()
	if err != nil {
		return err
	}

	err = this.printTableRows()
	if err != nil {
		return err
	}

	return nil
}

// Generic copying functionality.
func (this *MdFmt) Copy(t_fmt TableFmt) error {
	// We need to enforce the shadowed functions of the MdFmt struct.
	// Not the BaseTableFmt Copy().
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

	return nil
}

func InitMdFmt(t_label string) (MdFmt, error) {
	fmt_ := MdFmt{}

	fmt_.Label = t_label
	fmt_.ColWidth = make(map[int]int)
	fmt_.SortCol = SortUnset
	fmt_.NumericSortCol = SortUnset
	fmt_.Head = HeadUnset
	fmt_.Tail = TailUnset

	return fmt_, nil
}
