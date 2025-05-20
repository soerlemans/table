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

	// The query to perform.
	ProgramText string

	// Internal table representation.
	Table td.TableData

	// Raw input in row format.
	// Input []string
}

// Gives a unique id for every Processcontext.
var idCounter uint64

func initProcessContext(t_text string, t_table td.TableData) ProcessContext {
	ctx := ProcessContext{idCounter, t_text, t_table}
	defer u.LogStructName("initContext", ctx, u.ETC80)

	// Increment id counter.
	idCounter++

	return ctx
}

func Process(t_ctx ProcessContext) []string {
	// Parse filtering code to create a Pipe like data sctructure.
	// Some kind of decorator structure which.
	// Can then be executed like an AST.
	// Just create a single data type for processing.
	// something like a Table structure, consisting of columns, rows, etc.
	f.InitFilter(t_ctx.ProgramText)

	rows := t_ctx.Table.RowsData
	for i, line := range rows {
		u.Printf("line(%d:%d): %s", i, t_ctx.Id, line)
	}

	return nil
}
