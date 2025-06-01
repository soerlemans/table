package ir

import (
	"fmt"

	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
)

func resolveValue(t_value Value) string {
	return ""
}

// TODO: receive an output buffer to write to or something.
// Or something else.
func ExecIr(instructions InstructionList, t_index int, t_table *td.TableData) error {
	out, err := t_table.RowAsStr(t_index)
	if err != nil {
		return err
	}

	for _, inst := range instructions {
		switch inst.Type {
		case Equal:
			break

		default:
			u.Logln("ExecIr: Error unhandeld Instruction.")
			break
		}

	}

	fmt.Println(out)

	return nil
}
