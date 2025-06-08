package table_fmt

import (
	"fmt"
	"strings"

	td "github.com/soerlemans/table/table_data"
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
			// TODO: Return err.
		}

		// Check if the column is selected.
		colSep := strings.Repeat("-", colWidth)
		fmt.Printf("+ %s ", colSep)
	}
	fmt.Println("+")

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
	fmt_ := PrettyFmt{}

	fmt_.Label = t_label
	fmt_.ColWidth = make(map[int]int)
	fmt_.SortCol = SortUnset
	fmt_.NumericSortCol = SortUnset
	fmt_.Head = HeadUnset
	fmt_.Tail = TailUnset

	return fmt_, nil
}
