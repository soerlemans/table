package main

import (
	f "github.com/soerlemans/table/filter"
	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

// Execution context of a procss required for running the program.
type ProcessContext struct {
	// Identifier.
	Id uint64

	// Pointer to filters
	Filter *f.Filter

	// Internal table representation.
	Table td.TableData
}

// Gives a unique id for every Processcontext.
var idCounter uint64

func initProcessContext(t_filter *f.Filter, t_table td.TableData) ProcessContext {
	ctx := ProcessContext{idCounter, t_filter, t_table}
	defer func() { u.LogStructName("initContext", ctx, u.ETC80) }()

	// Increment id counter.
	idCounter++

	return ctx
}

func Process(t_ctx ProcessContext) []string {
	// Parse filtering code to create a Pipe like data structure.
	// Some kind of decorator structure which.
	// Can then be executed like an AST.
	// Just create a single data type for processing.
	// something like a Table structure, consisting of columns, rows, etc.

	rows := t_ctx.Table.RowsData
	for index, _ := range rows {
		filter := t_ctx.Filter
		tablePtr := &t_ctx.Table

		// Execute filters on the current row index.
		filter.Exec(index, tablePtr)

		// u.Printf("line(%d:%d): %s", index, t_ctx.Id, line)
	}

	return nil
}
