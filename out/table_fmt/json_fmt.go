package table_fmt

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
)

type JsonFmt struct {
	// Includes base data and methods.
	BaseTableFmt
}

func (this *JsonFmt) printRow(t_row td.TableDataRow) error {
	var sep string
	order := this.GetOrder()

	fmt.Printf("{ ")
	for _, index := range order {
		colName := this.Headers[index]
		value := t_row[index]

		fmt.Printf("%s\"%s\": \"%s\"", sep, colName, value)
		sep = ", "
	}
	fmt.Println(" }")

	return nil
}

func (this *JsonFmt) printTableRows() error {
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

func (this *JsonFmt) Write() error {
	err := this.printTableRows()
	if err != nil {
		return err
	}

	return nil
}

func InitJsonFmt(t_label string) (JsonFmt, error) {
	fmt_ := JsonFmt{}

	fmt_.Label = t_label

	return fmt_, nil
}
