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

	// We use two spaces for indentation.
	Identation string = "  "
)

type HtmlFmt struct {
	// Includes base data and methods.
	BaseTableFmt
}

func (this *HtmlFmt) indent(t_level IdentLevel, t_fmt string, t_args ...interface{}) {
	level := int(t_level)

	// We use 2 spaces for indentation.
	str := strings.Repeat(Identation, level)
	format := fmt.Sprintf("%s%s\n", str, t_fmt)

	this.writef(format, t_args...)
}

func (this *HtmlFmt) printRow(t_tag string, t_row td.TableDataRow) error {
	order := this.GetOrder()

	this.indent(Row, "<tr>")
	for _, index := range order {
		cell := t_row[index]

		this.indent(Cell, "<%s> %s </%s>", t_tag, cell, t_tag)
	}
	this.indent(Row, "</tr>")

	return nil
}

func (this *HtmlFmt) printTableHeader() error {
	this.indent(Section, "<thead>")
	err := this.printRow("td", this.Headers)
	this.indent(Section, "</thead>")

	return err
}

func (this *HtmlFmt) printTableRows() error {
	this.indent(Section, "<tbody>")
	// Print per row.
	for index, row := range this.Rows {
		// Skip if we are not in bounds.
		if !this.InBounds(index) {
			continue
		}

		// Print cells of the row.
		err := this.printRow("tr", row)
		if err != nil {
			return err
		}
	}
	this.indent(Section, "</tbody>")

	return nil
}

func (this *HtmlFmt) Write() error {
	this.indent(Table, "<table>")
	err := this.printTableHeader()
	if err != nil {
		return err
	}

	err = this.printTableRows()
	if err != nil {
		return err
	}

	this.indent(Table, "</table>")

	return nil
}

func InitHtmlFmt(t_label string) (HtmlFmt, error) {
	base, err := InitBaseTableFmt(t_label)
	format := HtmlFmt{BaseTableFmt: base}
	if err != nil {
		return format, err
	}

	return format, nil
}
