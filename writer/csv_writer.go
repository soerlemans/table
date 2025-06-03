package writer

import (
	td "github.com/soerlemans/table/table_data"
)

type CsvWriter struct {
	// We need to calculate the max column width for every entry..
	Headers  td.TableDataRow

	Rows []td.TableDataRow

	Label string
}

func (this *CsvWriter) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers
}

func (this *CsvWriter) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)
}

func (this *CsvWriter) GetLabel() string {
	return this.Label
}

func (this *CsvWriter) Write() error {
	// TODO: Implement markdown writing logic.

	return nil
}

func InitCsvWriter(t_label string) (CsvWriter, error) {
	writer := CsvWriter{}

	writer.Label = t_label

	return writer, nil
}
