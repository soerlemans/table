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
	Label string

	// We need to calculate the max column width for every entry..
	ColWidth map[int]int
	Headers  td.TableDataRow

	ColMask map[int]bool

	Rows []td.TableDataRow
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

func (this *MdFmt) GetLabel() string {
	return this.Label
}

func (this *MdFmt) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers

	this.updateColWidth(t_headers)
}

// Mark columns to print during write.
func (this *MdFmt) SetMask(t_mask []int) {
	this.ClearMask()

	// TODO: Error handle non existent column indexes.
	for _, value := range t_mask {
		this.ColMask[value] = true
	}
}

func (this *MdFmt) ClearMask() {
	// Clear by assigning a new one.
	this.ColMask = make(map[int]bool)
}

func (this *MdFmt) ColumnMasked(t_colIndex int) bool {
	// Guard clause (the mask has no elements then print everything).
	// As we should always print atleast one column.
	if len(this.ColMask) == 0 {
		return true
	}

	// We use the map as a set.
	_, ok := this.ColMask[t_colIndex]

	return ok
}

func (this *MdFmt) SetRows(t_rows []td.TableDataRow) {
	this.Rows = t_rows

	// We need to update the mex column width for every line now.
	for _, row := range this.Rows {
		this.updateColWidth(row)
	}
}

func (this *MdFmt) GetRows() []td.TableDataRow {
	return this.Rows
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

func InitMdFmt(t_label string) (MdFmt, error) {
	writer := MdFmt{}

	writer.Label = t_label
	writer.ColWidth = make(map[int]int)
	writer.ColMask = make(map[int]bool)

	return writer, nil
}
