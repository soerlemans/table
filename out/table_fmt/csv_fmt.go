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
	order := this.GetOrder()

	var sep string
	for _, index := range order {
		cell := t_row[index]
		fmt.Printf("%s%s", sep, cell)

		sep = ","
	}
	fmt.Println()

	return nil
}

func (this *CsvFmt) printTableHeader() error {
	return this.printRow(this.Headers)
}

func (this *CsvFmt) printTableRows() error {
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

func (this *CsvFmt) Write() error {
	this.PerformSort()

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
	fmt_.SortCol = SortUnset
	fmt_.NumericSortCol = SortUnset
	fmt_.Head = HeadUnset
	fmt_.Tail = TailUnset

	return fmt_, nil
}
