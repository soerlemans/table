package table_fmt

import (
	"fmt"
	"strings"

	td "github.com/soerlemans/table/table_data"
)

const (
	mdColSep = '|'
)

type MdFmt struct {
	// Includes base data and methods.
	BaseTableFmt

	ColWidth map[int]int
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
	for index, cell := range t_row {
		colWidth, ok := this.ColWidth[index]
		if !ok {
			// TODO: Return err.
		}

		// Check if the column is selected.
		if this.ColumnMasked(index) {
			fmt.Printf("| %-*s ", colWidth, cell)
		}
	}
	fmt.Println("|")

	return nil
}

func (this *MdFmt) printTableHeader() error {
	return this.printRow(this.Headers)
}

func (this *MdFmt) printTableHeaderSep() error {
	for index, _ := range this.Headers {
		colWidth, ok := this.ColWidth[index]
		if !ok {
			// TODO: Return err.
		}

		// Check if the column is selected.
		if this.ColumnMasked(index) {
			colSep := strings.Repeat("-", colWidth)
			fmt.Printf("| %s ", colSep)
		}
	}
	fmt.Println("|")

	return nil
}

func (this *MdFmt) printTableRows() error {
	// Print per row.
	for _, row := range this.Rows {
		// Print cells of the row.
		err := this.printRow(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *MdFmt) Write() error {
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
	this.Label = t_fmt.GetLabel()

	headers := t_fmt.GetHeaders()
	this.SetHeaders(headers)

	rows := t_fmt.GetRows()
	this.SetRows(rows)

	mask := t_fmt.GetMask()
	this.SetMask(mask)

	return nil
}

func InitMdFmt(t_label string) (MdFmt, error) {
	fmt_ := MdFmt{}

	fmt_.Label = t_label
	fmt_.ColWidth = make(map[int]int)
	fmt_.ColMask = make(map[int]bool)

	return fmt_, nil
}
