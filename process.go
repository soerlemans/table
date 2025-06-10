package main

import (
	f "github.com/soerlemans/table/filter"
	ir "github.com/soerlemans/table/filter/ir"
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

func Process(t_ctx ProcessContext) error {
	filter := t_ctx.Filter
	vm, err := ir.InitIrVm(&t_ctx.Table)
	if err != nil {
		return err
	}

	// We can modify redundant instructions during execution.
	instructions := filter.Instructions
	err = vm.Exec(*instructions)
	if err != nil {
		return err
	}

	return vm.Write()
}
