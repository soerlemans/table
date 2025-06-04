package table_fmt

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
)

type CsvFmt struct {
	Label string

	// We need to calculate the max column width for every entry..
	Headers td.TableDataRow

	ColMask map[int]bool

	Rows []td.TableDataRow
}

func (this *CsvFmt) printRow(t_row td.TableDataRow) error {
	var sep string
	for index, cell := range t_row {
		if this.ColumnMasked(index) {
			fmt.Printf("%s%s", sep, cell)

			sep = ","
		}
	}
	fmt.Println()

	return nil
}

func (this *CsvFmt) printTableHeader() error {
	return this.printRow(this.Headers)
}

func (this *CsvFmt) printTableRows() error {
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

func (this *CsvFmt) GetLabel() string {
	return this.Label
}

func (this *CsvFmt) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers
}

// Mark columns to print during write.
func (this *CsvFmt) SetMask(t_mask []int) {
	this.ClearMask()

	// TODO: Error handle non existent column indexes.
	for _, value := range t_mask {
		this.ColMask[value] = true
	}
}

func (this *CsvFmt) ClearMask() {
	// Clear by assigning a new one.
	this.ColMask = make(map[int]bool)
}

func (this *CsvFmt) ColumnMasked(t_colIndex int) bool {
	// Guard clause (the mask has no elements then print everything).
	// As we should always print atleast one column.
	if len(this.ColMask) == 0 {
		return true
	}

	// We use the map as a set.
	_, ok := this.ColMask[t_colIndex]

	return ok
}

func (this *CsvFmt) SetRows(t_rows []td.TableDataRow) {
	this.Rows = t_rows
}

func (this *CsvFmt) GetRows() []td.TableDataRow {
	return this.Rows
}

func (this *CsvFmt) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)
}

func (this *CsvFmt) Write() error {
	err := this.printTableHeader()
	if err != nil {
		return err
	}

	err = this.printTableRows()
	if err != nil {
		return err
	}

	return nil
}

func InitCsvFmt(t_label string) (CsvFmt, error) {
	writer := CsvFmt{}

	writer.Label = t_label
	writer.ColMask = make(map[int]bool)

	return writer, nil
}
