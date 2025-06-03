package writer

import (
	td "github.com/soerlemans/table/table_data"
)

type MdWriter struct {
	// We need to calculate the max column width for every entry..
	ColWidth map[int]int
	Headers  td.TableDataRow

	Rows []td.TableDataRow

	Label string
}

func (this *MdWriter) updateColWidth(t_row td.TableDataRow) {
	for i, cell := range t_row {
		cellWidth := len(cell)

		// Update the max column size.
		if cellWidth > this.ColWidth[i] {
			this.ColWidth[i] = cellWidth
		}
	}
}

func (this *MdWriter) SetHeaders(t_headers td.TableDataRow) {
	this.Headers = t_headers

	this.updateColWidth(t_headers)
}

func (this *MdWriter) AddRow(t_row td.TableDataRow) {
	this.Rows = append(this.Rows, t_row)

	this.updateColWidth(t_row)
}

func (this *MdWriter) GetLabel() string {
	return this.Label
}

func (this *MdWriter) Write() error {
	// TODO: Implement markdown writing logic.

	return nil
}

func InitMdWriter(t_label string) (MdWriter, error) {
	return MdWriter{}, nil
}
