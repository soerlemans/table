package ir

import (
	"fmt"
)

// Convenience alias.
type InstructionList = []Instruction
type ValueList = []Value

type InstructionType int

const (
	xxx InstructionType = iota

	LoadVariable
	StoreVariable
	LoadField

	LessThan
	LessThanEqual
	Equal
	NotEqual
	GreaterThan
	GreaterThanEqual

	WhenBlock
	MutBlock
	OutBlock
	MdBlock
	JsonBlock
)

type ValueType int

const (
	Identifier ValueType = iota

	String
	Number

	FieldByName
	FieldByPosition
)

type Value struct {
	Type  ValueType
	Value string
}

// TODO:
type Instruction struct {
	Id       string
	Type     InstructionType
	Operands ValueList
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
	label := fmt.Sprintf("%i%d", labelCount)
	labelCount++

	return Instruction{label, t_type, t_operands}
}
