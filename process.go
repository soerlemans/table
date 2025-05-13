package main

import (
	"github.com/soerlemans/table/util"
)

// Execution context of a procss required for running the program.
type Context struct {
	// Identifier.
	Id uint64

	Table TableData

	// Raw input in row format.
	// Input []string
}

// Gives a unique id for every context.
var idCounter uint64

func initContext(t_table TableData) Context {
	ctx := Context{idCounter, t_table}

	// inputLength := len(t_input)
	util.Logf("initContext: %+v", ctx)

	// Increment id counter.
	idCounter++

	return ctx
}

func Process(t_ctx Context) []string {
	// Parse filtering code to create a Pipe like data sctructure.
	// Some kind of decorator structure which.
	// Can then be executed like an AST.
	// Just create a single data type for processing.
	// something like a Table structure, consisting of columns, rows, etc.

	// for i, line := range t_ctx.Input {
	// 	util.Printf("line(%d:%d): %s", i, t_ctx.Id, line)
	// }

	return nil
}
