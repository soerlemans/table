package filter

import (
	td "github.com/soerlemans/table/table_data"
)

// Executes the AST into filtering
type Interpreter struct {
	// We process row by row.
	RowIndex int64

	// Table of data to spit through.
	Table *td.TableData
}

func (this *Interpreter) Visit() {

}
