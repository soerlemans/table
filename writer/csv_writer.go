package writer

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
)

type CsvWriter struct {
	Label string

	// We need to calculate the max column width for every entry..
	Headers  td.TableDataRow
	Rows []td.TableDataRow
}

func (this *CsvWriter) printRow(t_row td.TableDataRow) error {
	var sep string
	for _, cell := range t_row {
		fmt.Printf("%s%s", sep, cell)

		sep = ","
	}
	fmt.Println()

	return nil
}

func (this *CsvWriter) printTableHeader() error {
	return this.printRow(this.Headers)
}

func (this *CsvWriter) printTableRows() error {
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

func (this *CsvWriter) GetLabel() string {
	return this.Label
}

func (this *CsvWriter) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers
}

func (this *CsvWriter) SetRows(t_rows []td.TableDataRow) {
	this.Rows = t_rows
}

func (this *CsvWriter) GetRows() []td.TableDataRow {
	return this.Rows
}

func (this *CsvWriter) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)
}


func (this *CsvWriter) Write() error {
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

func InitCsvWriter(t_label string) (CsvWriter, error) {
	writer := CsvWriter{}

	writer.Label = t_label

	return writer, nil
}
