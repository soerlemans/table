package table_fmt

import (
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
		this.writef("%s%s", sep, cell)

		sep = ","
	}
	this.writeln()

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
	base, err := InitBaseTableFmt(t_label)
	format := CsvFmt{BaseTableFmt: base}
	if err != nil {
		return format, err
	}

	return format, nil
}
