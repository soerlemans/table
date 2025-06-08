package filter

import (
	"errors"
	"fmt"
	// "strconv"

	"github.com/soerlemans/table/filter/ir"
	u "github.com/soerlemans/table/util"
)

// Errors:
var (
	ErrInstructionEmpty = errors.New("Expected an instruction")
	ErrUnexpectedEos    = errors.New("Unexpected end of stream")
	ErrExpected         = errors.New("Expected")
	ErrUnsupported      = errors.New("Unsupported")
	ErrWriteDirective   = errors.New("No further input allowed after a write directive")
)

// Import these as we use them frequently.
type InstList = ir.InstructionList

type InstPtr = *ir.Instruction
type InstListPtr = *InstList

type parseFnInst = func(*TokenStream) (InstPtr, error)

type ValuePtr = *ir.Value
type ValueListPtr = *ir.ValueList

type parseFnValue = func(*TokenStream) (ValuePtr, error)

// TODO: Implement.
func errInstructionEmpty(t_location string) error {
	return fmt.Errorf("Expected Instruction returned nil in %s (%w)", t_location, ErrInstructionEmpty)
}

func errExpectedString(t_location string, t_string string) error {
	return fmt.Errorf("Expected a \"%s\" in %s (%w)", t_string, t_location, ErrExpected)
}

func errExpectedTokenType(t_location string, t_type TokenType) error {
	return errExpectedString(t_location, t_type.String())
}

func errUnexpectedEos(t_location string) error {
	return fmt.Errorf("Unexpected end of stream in %s (%w)", t_location, ErrUnexpectedEos)
}

func errWriteDirective(t_location string) error {
	return fmt.Errorf("After a write directive end of stream is expected in %s (%w)", t_location, ErrWriteDirective)
}

// We need a pointer to a pointer to modify it.
func addInstruction(t_insts *InstListPtr, t_inst InstPtr) error {
	if t_inst == nil {
		panic("InstructionList pointer is nil!")
	}

	if t_inst == nil {
		panic("Instruction is nil!")
	}

	// Init the instruction list.
	if *t_insts == nil {
		list := ir.InitInstructionList()

		*t_insts = &list
	}

	// Push back the value.
	(*t_insts).PushBack(t_inst)

	return nil
}

// Should be used in conjunction with defer.
func logUnlessNil[T any](t_preabmle string, t_ptr *T) {
	if t_ptr != nil {
		u.Logf("Found %s: %v", t_preabmle, *t_ptr)
	}
}

// Check if a Inst is initialized if the Inst is not nil.
// And no error is set.
func validPtr[T any](t_ptr *T, t_err error) bool {
	return t_ptr != nil && t_err == nil
}

// TODO: Refactor two instances of parseLists.
func parseList(t_stream *TokenStream, t_fn parseFnInst, t_sep TokenType) (InstListPtr, error) {
	var list InstListPtr
	defer func() { logUnlessNil("parseList", list) }()

	for {
		inst, err := t_fn(t_stream)
		if err != nil {
			return list, err
		}

		// If the inst is not nil we found an item.
		if inst != nil {
			// Add instruction to the list.
			addInstruction(&list, inst)
		} else {
			break
		}

		// We must check for any separator symbols.
		if !t_stream.Eos() {
			token := t_stream.Current()

			// If no intermediary pipe symbol was found we should quit.
			if token.Type == t_sep {
				u.Logf("Found sep: %s", token.Type)

				// Skip past the separator.
				t_stream.Next()

				// If after the separator we get an eos this unexpected.
				if t_stream.Eos() {
					return list, errUnexpectedEos("parseList")
				}
			} else {
				break
			}
		}
	}

	return list, nil
}

// TODO: Create generic function which unifies ParseList and ParseValueList.
func parseValueList(t_stream *TokenStream, t_fn parseFnValue, t_sep TokenType) (ValueListPtr, error) {
	var list ValueListPtr = new(ir.ValueList)
	defer func() { logUnlessNil("parseValueList", list) }()

	for {
		inst, err := t_fn(t_stream)
		if err != nil {
			return list, err
		}

		// If the inst is not nil we found an item.
		if inst != nil {
			*list = append(*list, *inst)
		} else {
			break
		}

		// We must check for any separator symbols.
		if !t_stream.Eos() {
			token := t_stream.Current()

			// If no intermediary pipe symbol was found we should quit.
			if token.Type == t_sep {
				u.Logf("Found sep: %s", token.Type)

				// Skip past the separator.
				t_stream.Next()

				// If after the separator we get an eos this unexpected.
				if t_stream.Eos() {
					return list, errUnexpectedEos("parseList")
					break
				}
			} else {
				break
			}
		}
	}

	return list, nil
}

// Parsing functions:
func accessorName(t_stream *TokenStream) (ValuePtr, error) {
	var value ValuePtr
	defer func() { logUnlessNil("accessorName", value) }()

	// Guard clause.
	if t_stream.Eos() {
		return value, nil
	}

	dot := t_stream.Current()

	// TODO: Refactor boilerplate later.
	if dot.Type == ACCESSOR_NAME {
		u.Logf("Found a '%v'", dot.Type)
		t_stream.Next()

		if t_stream.Eos() {
			return value, errUnexpectedEos("accessorName")
		}

		// Identifier and string are both converted to the same value.
		token := t_stream.Current()
		switch token.Type {
		case IDENTIFIER:
			fallthrough
		case STRING:
			name := ir.InitValue(ir.FieldByName, token.Value)

			value = &name
			break

		default:
			return value, errExpectedString("accessorName", "identifier or string")
		}

	}

	return value, nil
}

func accessorPositional(t_stream *TokenStream) (ValuePtr, error) {
	var value ValuePtr
	defer func() { logUnlessNil("accessorName", value) }()

	// Guard clause.
	if t_stream.Eos() {
		return value, nil
	}

	dollarSign := t_stream.Current()

	// TODO: Refactor boilerplate later.
	if dollarSign.Type == ACCESSOR_POSITIONAL {
		u.Logf("Found a '%v'", dollarSign.Type)
		t_stream.Next()

		if t_stream.Eos() {
			return value, errUnexpectedEos("accessorPositional")
		}

		// Identifier value must be used to convert to a string later.
		// Possibly add expressions or arithmetic expressions later down the line.
		token := t_stream.Current()
		switch token.Type {
		// FIXME: For now not allowed.
		// case IDENTIFIER:
		// 	return value, nil

		case NUMBER:
			// integer, err := strconv.ParseInt(token.Value, 10, 64)
			pos := ir.InitValue(ir.FieldByPosition, token.Value)

			value = &pos
			break

		default:
			return value, errExpectedString("accessorPositional", "identifier or number")
		}

	}

	return value, nil
}

func column(t_stream *TokenStream) (ValuePtr, error) {
	var value ValuePtr

	// Check for a name based column accessor.
	if name, err := accessorName(t_stream); validPtr(name, err) {
		value = name
	} else if err != nil {
		return value, err

		// Check for a positional based column accessor.
	} else if pos, err := accessorPositional(t_stream); validPtr(pos, err) {
		value = pos
	} else {
		return value, err
	}

	return value, nil
}

func parameter(t_stream *TokenStream) (ValuePtr, error) {
	return rvalue(t_stream)
}

func parameterList(t_stream *TokenStream) (ValueListPtr, error) {
	return parseValueList(t_stream, parameter, COMMA)
}

// Ast:
func rvalue(t_stream *TokenStream) (ValuePtr, error) {
	var value ValuePtr
	defer func() { logUnlessNil("rvalue", value) }()

	// Guard clause.
	if t_stream.Eos() {
		return value, nil
	}

	token := t_stream.Current()

	// TODO: Refactor boilerplate later.
	switch token.Type {
	case NUMBER:
		number := ir.InitValue(ir.Number, token.Value)

		value = &number
		t_stream.Next()
		break

	case STRING:
		str := ir.InitValue(ir.String, token.Value)

		value = &str
		t_stream.Next()
		break

	case IDENTIFIER:
		id := ir.InitValue(ir.Identifier, token.Value)

		value = &id
		t_stream.Next()
		break

	default:
		if col, err := column(t_stream); validPtr(col, err) {
			value = col
			t_stream.Next()
		} else if err != nil {
			return value, err
		}
		break
	}

	return value, nil
}

// Initialize the comparison.
func initComparison(t_stream *TokenStream, t_type ir.InstructionType, t_lhs ValuePtr) (InstPtr, error) {
	var inst InstPtr

	t_stream.Next()
	if t_stream.Eos() {
		return inst, errUnexpectedEos("initComparison")
	}

	rhs, err := rvalue(t_stream)
	if err != nil {
		return inst, err
	}

	if rhs == nil {
		return inst, errExpectedString("initComparison", "right hand side of expression")
	}

	comp := ir.InitInstruction(t_type, *t_lhs, *rhs)
	inst = &comp

	return inst, nil
}

func expr(t_stream *TokenStream) (InstPtr, error) {
	var inst InstPtr
	defer func() { logUnlessNil("expr", inst) }()

	lhs, err := rvalue(t_stream)
	if err != nil {
		return inst, err
	}

	if lhs != nil {
		if t_stream.Eos() {
			return inst, errUnexpectedEos("expr")
		}

		token := t_stream.Current()

		// TODO: Cleanup boilerplate.
		switch token.Type {
		case LESS_THAN:
			ltPtr, err := initComparison(t_stream, ir.LessThan, lhs)
			if err != nil {
				return inst, err
			}

			inst = ltPtr
			break

		case LESS_THAN_EQUAL:
			ltePtr, err := initComparison(t_stream, ir.LessThanEqual, lhs)
			if err != nil {
				return inst, err
			}

			inst = ltePtr
			break

		case EQUAL:
			eqPtr, err := initComparison(t_stream, ir.Equal, lhs)
			if err != nil {
				return inst, err
			}

			inst = eqPtr
			break

		case NOT_EQUAL:
			nePtr, err := initComparison(t_stream, ir.NotEqual, lhs)
			if err != nil {
				return inst, err
			}

			inst = nePtr
			break

		case GREATER_THAN:
			gtPtr, err := initComparison(t_stream, ir.GreaterThan, lhs)
			if err != nil {
				return inst, err
			}

			inst = gtPtr
			break

		case GREATER_THAN_EQUAL:
			gtePtr, err := initComparison(t_stream, ir.GreaterThanEqual, lhs)
			if err != nil {
				return inst, err
			}

			inst = gtePtr
			break
		}
	}

	return inst, nil
}

func keyword(t_stream *TokenStream) (InstPtr, error) {
	var (
		inst InstPtr
		list ValueListPtr = new(ir.ValueList)
		err  error
	)
	defer func() { logUnlessNil("keyword", inst) }()

	// Guard clause.
	if t_stream.Eos() {
		return inst, nil
	}

	token := t_stream.Current()
	switch token.Type {
	case WHEN:
		t_stream.Next()
		if !t_stream.Eos() {
			list, err = parameterList(t_stream)
			if err != nil {
				return inst, err
			}
		}

		when := ir.InitInstructionByList(ir.When, *list)

		inst = &when
		break

	case MUT:
		t_stream.Next()
		if !t_stream.Eos() {
			list, err = parameterList(t_stream)
			if err != nil {
				return inst, err
			}
		}

		mut := ir.InitInstructionByList(ir.Mut, *list)

		inst = &mut
		break

	case Sort:
		t_stream.Next()
		if t_stream.Eos() {
			// TODO: Handle end of stream expected a numeric.
		}

		if argPtr, err := rvalue(t_stream); validPtr(argPtr, err) {
			head := ir.InitInstruction(ir.Head, *argPtr)

			inst = &head
		} else {
			// We must receive an argument for head.
			return inst, errExpectedString("keyword", "Number or identifier")
		}
		break

	case NumericSort:
		t_stream.Next()
		if t_stream.Eos() {
			// TODO: Handle end of stream expected a numeric.
		}

		if argPtr, err := rvalue(t_stream); validPtr(argPtr, err) {
			head := ir.InitInstruction(ir.Head, *argPtr)

			inst = &head
		} else {
			// We must receive an argument for head.
			return inst, errExpectedString("keyword", "Number or identifier")
		}
		break

	case HEAD:
		t_stream.Next()
		if t_stream.Eos() {
			// TODO: Handle end of stream expected a numeric.
		}

		if argPtr, err := rvalue(t_stream); validPtr(argPtr, err) {
			head := ir.InitInstruction(ir.Head, *argPtr)

			inst = &head
		} else {
			// We must receive an argument for head.
			return inst, errExpectedString("keyword", "Number or identifier")
		}
		break

	case TAIL:
		t_stream.Next()
		if t_stream.Eos() {
			// TODO: Handle end of stream expected a numeric.
		}

		if argPtr, err := rvalue(t_stream); validPtr(argPtr, err) {
			tail := ir.InitInstruction(ir.Tail, *argPtr)

			inst = &tail
		} else {
			// We must receive an argument for head.
			return inst, errExpectedString("keyword", "Number or identifier")
		}
		break

	default:
		// No error this is fine, just not a keyword.
		break
	}

	return inst, nil
}

func write(t_stream *TokenStream) (InstPtr, error) {
	var (
		inst InstPtr
		list ValueListPtr = new(ir.ValueList)
		err  error
	)

	defer func() { logUnlessNil("write", inst) }()

	// Guard clause.
	if t_stream.Eos() {
		return inst, nil
	}

	token := t_stream.Current()
	switch token.Type {
	case CSV:
		t_stream.Next()
		if !t_stream.Eos() {
			list, err = parameterList(t_stream)
			if err != nil {
				return inst, err
			}
		}

		csv := ir.InitInstructionByList(ir.Csv, *list)

		inst = &csv
		break

	case MD:
		t_stream.Next()
		if !t_stream.Eos() {
			list, err = parameterList(t_stream)
			if err != nil {
				return inst, err
			}
		}

		md := ir.InitInstructionByList(ir.Md, *list)

		inst = &md
		break

	case JSON:
		t_stream.Next()
		if !t_stream.Eos() {
			list, err = parameterList(t_stream)
			if err != nil {
				return inst, err
			}
		}

		json := ir.InitInstructionByList(ir.Json, *list)

		inst = &json
		break

	case HTML:
		t_stream.Next()
		if !t_stream.Eos() {
			list, err = parameterList(t_stream)
			if err != nil {
				return inst, err
			}
		}

		html := ir.InitInstructionByList(ir.Html, *list)

		inst = &html
		break

	default:
		// No error this is fine, just not a keyword.
		break
	}

	return inst, nil
}

func item(t_stream *TokenStream) (InstPtr, error) {
	var inst InstPtr
	defer func() { logUnlessNil("item", inst) }()

	// Check for keywords.
	if keywordPtr, err := keyword(t_stream); validPtr(keywordPtr, err) {
		inst = keywordPtr
	} else if err != nil {
		return inst, err

		// Check for an expression.
	} else if exprPtr, err := expr(t_stream); validPtr(exprPtr, err) {
		inst = exprPtr
	} else if err != nil {
		return inst, err

		// Check for a writer specification.
	} else if writePtr, err := write(t_stream); validPtr(writePtr, err) {
		inst = writePtr

		// If a writer is successfully parsed we should be at the end of stream.
		if !t_stream.Eos() {
			return inst, errWriteDirective("item")
		}
	} else if err != nil {
		return inst, err
	}

	return inst, nil
}

func itemList(t_stream *TokenStream) (InstListPtr, error) {
	return parseList(t_stream, item, PIPE)
}

// This function is here purely just to match the grammary.yy.
func program(t_stream *TokenStream) (InstListPtr, error) {
	return itemList(t_stream)
}

// Source code to parse.
func Parse(t_stream *TokenStream) (InstListPtr, error) {
	var (
		list InstListPtr
		err  error
	)

	u.Logf("BEGIN PARSING.")
	defer u.Logf("END PARSING.")

	if t_stream.Len() > 0 {
		// If we have received tokens start parsing.
		list, err = program(t_stream)
		if err != nil {
			return list, err
		}
	} else {
		// If no token were found, assume default query (csv output).
		vlist := new(ir.ValueList)
		csv := ir.InitInstructionByList(ir.Csv, *vlist)

		// Add instruction to the list.
		addInstruction(&list, &csv)

		u.Logln("No tokens found, using default csv output writer.")
	}

	u.Logf("Instructions: %+v", list.String())

	return list, nil
}
