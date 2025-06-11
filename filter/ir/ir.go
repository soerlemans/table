package ir

import (
	"fmt"
)

// Convenience alias.
type ValueList = []Value
type InstListPtr = *InstructionList

type InstructionType int

const (
	// TODO: figure out if needed.
	// LoadVariable
	// StoreVariable
	// LoadField

	LessThan InstructionType = iota
	LessThanEqual
	Equal
	NotEqual
	GreaterThan
	GreaterThanEqual

	When
	Mut

	Sort
	NumericSort

	Head
	Tail

	Csv
	Md
	Pretty
	Json
	Html

	WriteDirective
)

func (t InstructionType) String() string {
	switch t {
	case LessThan:
		return "LessThan"
	case LessThanEqual:
		return "LessThanEqual"
	case Equal:
		return "Equal"
	case NotEqual:
		return "NotEqual"
	case GreaterThan:
		return "GreaterThan"
	case GreaterThanEqual:
		return "GreaterThanEqual"

	case When:
		return "When"
	case Mut:
		return "Mut"

	case Sort:
		return "Sort"
	case NumericSort:
		return "NumericSort"

	case Head:
		return "Head"
	case Tail:
		return "Tail"

	case Csv:
		return "Csv"
	case Md:
		return "Md"
	case Pretty:
		return "Pretty"
	case Json:
		return "Json"
	case Html:
		return "Html"

	case WriteDirective:
		return "WriteDirective"
	}

	// Return unhandled case.
	msg := fmt.Sprintf("<Unknown InstructionType (%v)>", t)
	return msg
}

type ValueType int

const (
	Identifier ValueType = iota

	String
	Number

	FieldByName
	FieldByPosition
)

func (t ValueType) String() string {
	switch t {
	case Identifier:
		return "Identifier"

	case String:
		return "String"
	case Number:
		return "Number"

	case FieldByName:
		return "FieldByName"
	case FieldByPosition:
		return "FieldByPosition"
	}

	// Optionally return an error?
	msg := fmt.Sprintf("<Unknown ValueType (%v)>", t)
	return msg
}

type Value struct {
	Type  ValueType
	Value string
}

func (this *Value) String() string {
	return fmt.Sprintf("%s", this.Value)
}

// TODO:
type Instruction struct {
	Label    string
	Type     InstructionType
	Operands ValueList
}

func (this *Instruction) String() string {
	return fmt.Sprintf("%s : %s <= %v", this.Label, this.Type, this.Operands)
}

// Initialization:
func InitValue(t_type ValueType, t_value string) Value {
	return Value{t_type, t_value}
}

func InitInstruction(t_type InstructionType, t_operands ...Value) Instruction {
	return InitInstructionByList(t_type, t_operands)
}

var labelCount int

func InitInstructionByList(t_type InstructionType, t_operands ValueList) Instruction {
	label := fmt.Sprintf("l%d", labelCount)
	labelCount++

	return Instruction{label, t_type, t_operands}
}
