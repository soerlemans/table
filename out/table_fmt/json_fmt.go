package table_fmt

import (
	td "github.com/soerlemans/table/table_data"
)

type JsonFmt struct {
	// Includes base data and methods.
	BaseTableFmt
}

func (this *JsonFmt) printRow(t_row td.TableDataRow) error {
	var sep string
	order := this.GetOrder()

	this.writef("{ ")
	for _, index := range order {
		colName := this.Headers[index]
		value := t_row[index]

		this.writef("%s\"%s\": \"%s\"", sep, colName, value)
		sep = ", "
	}
	this.writeln(" }")

	return nil
}

func (this *JsonFmt) printTableRows() error {
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

func (this *JsonFmt) Write() error {
	this.PerformSort()

	err := this.printTableRows()
	if err != nil {
		return err
	}

	return nil
}

func InitJsonFmt(t_label string) (JsonFmt, error) {
	base, err := InitBaseTableFmt(t_label)
	format := JsonFmt{BaseTableFmt: base}
	if err != nil {
		return format, err
	}

	return format, nil
}
