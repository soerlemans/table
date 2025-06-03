package ir

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
	w "github.com/soerlemans/table/writer"
)

type VmPtr = *IrVm

// Short of intermediary representation virtual machine.
type IrVm struct {
	VariableStore map[string]string

	Index int
	Table *td.TableData

	// The writer is in control of the final output result of the table program.
	Writer w.Writer
}

// TODO: The index and tableData should be wrapped in a struct or something.
// Now we pass them all the time together we should curry this.

func (this *IrVm) resolveValue(t_value Value) (string, error) {
	var value string

	index := this.Index
	table := this.Table

	switch t_value.Type {
	case Identifier:
		variable, ok := this.VariableStore[t_value.Value]
		if !ok {
			u.Logln("resolveValue: Non existent variable referenced.")

			// TODO: Error out.
			return "", nil
		}

		value = variable
		break

	case String:
		value = t_value.Value
		break

	case Number:
		// TODO: Figure out useful way to process numbers.
		value = t_value.Value
		break

	case FieldByName:
		name := t_value.Value
		cell, err := table.CellByColName(index, name)
		if err != nil {
			return value, err
		}

		value = cell
		break

	case FieldByPosition:
		colIndex, err := toInt(t_value.Value)
		if err != nil {
			return value, err
		}

		// Handle negative indices, to count from the end.
		if colIndex < 0 {
			colIndex = table.HeaderLength() - colIndex

			u.Logf("resolveValue: Negative indice %d.", colIndex)
		}

		cell, err := table.CellByIndices(index, colIndex)
		if err != nil {
			return value, err
		}

		value = cell
		break

	default:
		u.Logln("resolveValue: Error unhandeld ValueType.")
		// TODO: Error out.
		break

	}

	return value, nil
}

func (this *IrVm) binaryExprResolve(t_lhs Value, t_rhs Value) (string, string, error) {
	var (
		lhs string
		rhs string
	)

	lhs, err := this.resolveValue(t_lhs)
	if err != nil {
		return lhs, rhs, err
	}

	rhs, err = this.resolveValue(t_rhs)
	if err != nil {
		return lhs, rhs, err
	}

	return lhs, rhs, nil
}


func (this *IrVm) execComparison(t_type InstructionType, t_list ValueList) (bool, error) {
	var result bool

	// TODO: Err if size is != 2.
	lhs := t_list[0]
	rhs := t_list[1]

	resLhs, resRhs, err := this.binaryExprResolve(lhs, rhs)
	if err != nil {
		return result, err
	}

	// If there is a number comparison we must compare differently.
	// For equal and not equal (we assume numbers for all other comp operations).
	isNumberComp := (lhs.Type == Number || rhs.Type == Number)

	// TODO: Refactor boilerplate.
	switch t_type {
	case LessThan:
		intLhs, intRhs, err := binaryExprToInt(resLhs, resRhs)
		if err != nil {
			return result, err
		}

		result = (intLhs < intRhs)
		break

	case LessThanEqual:
		intLhs, intRhs, err := binaryExprToInt(resLhs, resRhs)
		if err != nil {
			return result, err
		}

		result = (intLhs <= intRhs)
		break

	// Equal and not Equal are always compared as strings.
	case Equal:
		if isNumberComp {
			// TODO: Implement conversion.
			result = (resLhs == resRhs)
		} else {
			result = (resLhs == resRhs)
		}
		break

	case NotEqual:
		if isNumberComp {
			// TODO: Implement conversion.
			result = (resLhs != resRhs)
		} else {
			result = (resLhs != resRhs)
		}
		break

	case GreaterThan:
		intLhs, intRhs, err := binaryExprToInt(resLhs, resRhs)
		if err != nil {
			return result, err
		}

		result = (intLhs > intRhs)
		break

	case GreaterThanEqual:
		intLhs, intRhs, err := binaryExprToInt(resLhs, resRhs)
		if err != nil {
			return result, err
		}

		result = (intLhs > intRhs)
		break
	}

	return result, nil
}

// TODO: receive an output buffer to write to or something.
// Or something else.
func (this *IrVm) ExecIr(instructions InstructionList) error {
	index := this.Index
	table := this.Table

	out, err := table.RowAsStr(index)
	if err != nil {
		return err
	}

	// Extract required global vars.
	rowLength := table.RowLength()
	headerLength := table.HeaderLength()

	// TODO: Refactor this init.
	// Set global variables for each execution.
	// These are similar to AWK.
	this.VariableStore["NF"] = fmt.Sprintf("%d", headerLength)
	this.VariableStore["FNR"] = fmt.Sprintf("%d", index)
	this.VariableStore["NR"] = fmt.Sprintf("%d", rowLength)

	skip := false
	for _, inst := range instructions {
		switch inst.Type {
		case LessThan:
			fallthrough
		case LessThanEqual:
			fallthrough
		case Equal:
			fallthrough
		case NotEqual:
			fallthrough
		case GreaterThan:
			fallthrough
		case GreaterThanEqual:
			cmp, err := this.execComparison(inst.Type, inst.Operands)
			// Consider if we should return/fail in this case as number errors could just be ignored.
			if err != nil {
				return err
			}

			// If the comparison was false we should skip printing the line.
			if !cmp {
				skip = true
			}
			break

		case Md:
			// Convert to markdown.
			break

		default:
			u.Logln("ExecIr: Error unhandeld InstructionType.")
			// TODO: Error out.
			break
		}

		if skip {
			break
		}
	}

	// Only print if we should not skip.
	if !skip {
		fmt.Println(out)
	}

	return nil
}

func InitIrVm(t_table *td.TableData) (IrVm, error) {
	var vm IrVm
	defer func() { u.LogStructName("InitIrVm", vm, u.ETC80) }()

	// Do init.
	vm.VariableStore = make(map[string]string)

	vm.Table = t_table
	writer, err := w.InitCsvWriter("ldefault")
	if err != nil {
		return vm, err
	}

	vm.Writer = &writer

	return vm, nil
}
