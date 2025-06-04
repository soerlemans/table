package ir

import (
	"fmt"

	tf "github.com/soerlemans/table/out/table_fmt"
	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

type VmPtr = *IrVm

// Short of intermediary representation virtual machine.
type IrVm struct {
	VariableStore map[string]string

	Index int
	Table *td.TableData

	// The formatter is in control of formatting the table in different formats.
	Fmt tf.TableFmt
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

func (this *IrVm) resolveValues(t_values ValueList) ([]string, error) {
	var slice []string

	for _, value := range t_values {
		str, err := this.resolveValue(value)
		if err != nil {
			return slice, err
		}

		slice = append(slice, str)
	}

	return slice, nil
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
		u.Logf("execComparison: %s == %s", resLhs, resRhs)
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
			u.Logln("ExecIr: Switching to md fmt.")
			md, err := tf.InitMdFmt(inst.Label)
			if err != nil {
				return err
			}

			// TODO: Write a copy/switch function.
			// Or maybe use a pointer for the internal data.
			// Update the headers.
			md.SetHeaders(this.Table.Headers)

			// Copy over the old registered rows.
			rows := this.Fmt.GetRows()
			md.SetRows(rows)

			colNames, err := this.resolveValues(inst.Operands)
			if err != nil {
				return err
			}

			indices, err := this.Table.ColNamesToIndices(colNames)
			if err != nil {
				return err
			}

			md.SetMask(indices)

			this.Fmt = &md
			break

		default:
			u.Logln("ExecIr: Error unhandeld InstructionType.")
			// TODO: Error out.
			break
		}

		// If we dont check all instructions we skip the setting of the writer.
		// Then a edge case appears where if non of the lines are matched.
		// We stick with the default CSV writer.
		// if skip {
		// 	break
		// }
	}

	// Only add the row to output if the current line isnt skipped.
	if !skip {
		row, err := table.GetRow(index)
		if err != nil {
			return err
		}

		this.Fmt.AddRow(row)
	}

	return nil
}

func (this *IrVm) Write() error {
	return this.Fmt.Write()
}

func InitIrVm(t_table *td.TableData) (IrVm, error) {
	var vm IrVm
	defer func() { u.LogStructName("InitIrVm", vm, u.ETC80) }()

	// Do init.
	vm.VariableStore = make(map[string]string)
	vm.Table = t_table

	// Set formatter field.
	fmt_, err := tf.InitCsvFmt("ldefault")
	if err != nil {
		return vm, err
	}

	vm.Fmt = &fmt_
	vm.Fmt.SetHeaders(vm.Table.Headers)

	return vm, nil
}
