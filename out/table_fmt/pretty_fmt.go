package table_fmt

import (
	"strings"
)

const (
	PrettyCorner = '+'
	PrettyColSep = '|'
	PrettyRowSep = '-'
)

type PrettyFmt struct {
	// A pretty table and markdown table are very similar.
	MdFmt

	ColWidth map[int]int
}

func (this *PrettyFmt) printTableHeaderSep() error {
	order := this.GetOrder()

	for _, index := range order {
		colWidth, ok := this.ColWidth[index]
		if !ok {
			return this.errColWidthIndex(index)
		}

		// Check if the column is selected.
		colSep := strings.Repeat("-", colWidth)
		this.writef("+ %s ", colSep)
	}
	this.writeln("+")

	return nil
}

func (this *PrettyFmt) Write() error {
	this.PerformSort()

	err := this.printTableHeaderSep()
	if err != nil {
		return err
	}

	err = this.printTableHeader()
	if err != nil {
		return err
	}

	err = this.printTableHeaderSep()
	if err != nil {
		return err
	}

	err = this.printTableRows()
	if err != nil {
		return err
	}

	err = this.printTableHeaderSep()
	if err != nil {
		return err
	}

	return nil
}

func InitPrettyFmt(t_label string) (PrettyFmt, error) {
	base, err := InitMdFmt(t_label)
	format := PrettyFmt{MdFmt: base}
	if err != nil {
		return format, err
	}

	return format, nil
}
