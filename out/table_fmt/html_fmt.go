package table_fmt

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
)

type HtmlFmt struct {
	// Includes base data and methods.
	BaseTableFmt
}

func (this *HtmlFmt) printRow(t_tag string, t_row td.TableDataRow) error {
	fmt.Println("<tr>")
	for index, cell := range t_row {
		// Check if the column is selected.
		if this.ColumnMasked(index) {
			fmt.Printf("<%s> %s </%s>\n", t_tag, cell, t_tag)
		}
	}
	fmt.Println("</tr>")

	return nil
}

func (this *HtmlFmt) printTableHeader() error {
	fmt.Println("<thead>")
	err := this.printRow("td", this.Headers)
	fmt.Println("</thead>")

	return err
}

func (this *HtmlFmt) printTableRows() error {
	// Print per row.
	for _, row := range this.Rows {
		// Print cells of the row.
		err := this.printRow("tr", row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *HtmlFmt) Write() error {
	fmt.Println("<table>")
	err := this.printTableHeader()
	if err != nil {
		return err
	}

	fmt.Println("</thead>")
	err = this.printTableRows()
	if err != nil {
		return err
	}
	fmt.Println("</thead>")

	fmt.Println("</table>")

	return nil
}

func InitHtmlFmt(t_label string) (HtmlFmt, error) {
	fmt_ := HtmlFmt{}

	fmt_.Label = t_label
	fmt_.ColMask = make(map[int]bool)

	return fmt_, nil
}
