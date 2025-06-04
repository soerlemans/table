package table_fmt

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
)

type CsvFmt struct {
	// Includes base data and methods.
	BaseTableFmt
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
	fmt_ := CsvFmt{}

	fmt_.Label = t_label
	fmt_.ColMask = make(map[int]bool)

	return fmt_, nil
}
