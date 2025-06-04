package table_fmt

import (
	"fmt"
	"strings"

	td "github.com/soerlemans/table/table_data"
)

type IdentLevel int

const (
	Table IdentLevel = iota
	Section
	Row
	Cell
)

type HtmlFmt struct {
	// Includes base data and methods.
	BaseTableFmt
}

func indent(t_level IdentLevel, t_fmt string, t_args ...interface{}) {
	level := int(t_level)

	// We use 2 spaces for indentation.
	str := strings.Repeat("  ", level)

	format := fmt.Sprintf("%s%s\n", str, t_fmt)

	fmt.Printf(format, t_args...)
}

func (this *HtmlFmt) printRow(t_tag string, t_row td.TableDataRow) error {
	indent(Row, "<tr>")
	for index, cell := range t_row {
		// Check if the column is selected.
		if this.ColumnMasked(index) {

			indent(Cell, "<%s> %s </%s>", t_tag, cell, t_tag)
		}
	}
	indent(Row, "</tr>")

	return nil
}

func (this *HtmlFmt) printTableHeader() error {
	indent(Section, "<thead>")
	err := this.printRow("td", this.Headers)
	indent(Section, "</thead>")

	return err
}

func (this *HtmlFmt) printTableRows() error {
	indent(Section, "<tbody>")
	// Print per row.
	for _, row := range this.Rows {
		// Print cells of the row.
		err := this.printRow("tr", row)
		if err != nil {
			return err
		}
	}
	indent(Section, "</tbody>")

	return nil
}

func (this *HtmlFmt) Write() error {
	indent(Table, "<table>")
	err := this.printTableHeader()
	if err != nil {
		return err
	}

	err = this.printTableRows()
	if err != nil {
		return err
	}

	indent(Table, "</table>")

	return nil
}

func InitHtmlFmt(t_label string) (HtmlFmt, error) {
	fmt_ := HtmlFmt{}

	fmt_.Label = t_label
	fmt_.ColMask = make(map[int]bool)

	return fmt_, nil
}
