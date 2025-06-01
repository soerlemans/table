package ir

import (
	"fmt"
	"strconv"

	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

var variableStore map[string]string

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
		integer, err := strconv.ParseInt(t_value.Value, 10, 64)
		if err != nil {
			return value, err
		}

		// Convert to normal int (watch out with integer slicing).
		colIndex := int(integer)

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

// TODO: receive an output buffer to write to or something.
// Or something else.
func ExecIr(instructions InstructionList, t_index int, t_table *td.TableData) error {
	out, err := t_table.RowAsStr(t_index)
	if err != nil {
		return err
	}

	// Set global variables for each execution.

	skip := false
	for _, inst := range instructions {
		switch inst.Type {
		case Equal:
			lhs := inst.Operands[0]
			rhs := inst.Operands[1]

			valueLhs, err := resolveValue(lhs, t_index, t_table)
			if err != nil {
				return err
			}

			valueRhs, err := resolveValue(rhs, t_index, t_table)
			if err != nil {
				return err
			}

			// If not equal skip the line.
			if valueLhs != valueRhs {
				skip = true
			}
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
