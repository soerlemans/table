package ir

import (
	"fmt"
	"strconv"

	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

// TODO: The index and tableData should be wrapped in a struct or something.
// Now we pass them all the time together we should curry this.

// Convert to normal int (watch out with integer slicing).
func toInt(t_str string) (int, error) {
	var result int

	integer, err := strconv.ParseInt(t_str, 10, 64)
	if err != nil {
		return result, err
	}

	result = int(integer)

	return result, nil
}

var variableStore = make(map[string]string)

func resolveValue(t_value Value, t_index int, t_table *td.TableData) (string, error) {
	var value string

	switch t_value.Type {
	case Identifier:
		variable, ok := variableStore[t_value.Value]
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
		cell, err := t_table.CellByColName(t_index, name)
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
			colIndex = t_table.HeaderLength() - colIndex

			u.Logf("resolveValue: Negative indice %d.", colIndex)
		}

		cell, err := t_table.CellByIndices(t_index, colIndex)
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

func execComparison(t_index int, t_table *td.TableData, t_type InstructionType, t_list ValueList) (bool, error) {
	var result bool

	lhs := t_list[0]
	rhs := t_list[1]

	resLhs, err := resolveValue(lhs, t_index, t_table)
	if err != nil {
		return result, err
	}

	resRhs, err := resolveValue(rhs, t_index, t_table)
	if err != nil {
		return result, err
	}

	// If there is a number comparison we must compare differently.
	// For equal and not equal (we assume numbers for all other comp operations).
	isNumberComp := (lhs.Type == Number || rhs.Type == Number)

	// TODO: Refactor boilerplate.
	switch t_type {
	case LessThan:
		intLhs, err := toInt(resLhs)
		if err != nil {
			return result, err
		}

		intRhs, err := toInt(resRhs)
		if err != nil {
			return result, err
		}

		result = (intLhs < intRhs)

		u.Logf("LessThan: %d < %d", intLhs, intRhs)
		break

	case LessThanEqual:
		intLhs, err := toInt(resLhs)
		if err != nil {
			return result, err
		}

		intRhs, err := toInt(resRhs)
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
		intLhs, err := toInt(resLhs)
		if err != nil {
			return result, err
		}

		intRhs, err := toInt(resRhs)
		if err != nil {
			return result, err
		}

		result = (intLhs > intRhs)
		break

	case GreaterThanEqual:
		intLhs, err := toInt(resLhs)
		if err != nil {
			return result, err
		}

		intRhs, err := toInt(resRhs)
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
func ExecIr(instructions InstructionList, t_index int, t_table *td.TableData) error {
	out, err := t_table.RowAsStr(t_index)
	if err != nil {
		return err
	}

	// Extract required global vars.
	rowLength := t_table.RowLength()
	headerLength := t_table.HeaderLength()

	// Set global variables for each execution.
	// These are similar to AWK.
	variableStore["FNR"] = fmt.Sprintf("%d", t_index)
	variableStore["NR"] = fmt.Sprintf("%d", rowLength)

	variableStore["HR"] = fmt.Sprintf("%d", headerLength)

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
			cmp, err := execComparison(t_index, t_table, inst.Type, inst.Operands)
			// Consider if we should return/fail in this case as number errors could just be ignored.
			if err != nil {
				return err
			}

			// If the comparison was false we should skip printing the line.
			if !cmp {
				skip = true
			}
			break

		case MdBlock:
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
