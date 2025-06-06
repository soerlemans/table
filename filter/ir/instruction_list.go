package ir

import (
	l "container/list"
	"fmt"
)

// For the love of god Golang doesnt have a list which uses generics.
// So we need to enforce our own type safety.
// In a terrible way.
type InstructionList struct {
	list *l.List
}

// PushBack adds an *Instruction to the back of the list.
func (this *InstructionList) PushBack(instr *Instruction) *l.Element {
	if instr == nil {
		panic("cannot add nil Instruction")
	}

	return this.list.PushBack(instr)
}

// PushFront adds an *Instruction to the front of the list.
func (this *InstructionList) PushFront(t_inst *Instruction) *l.Element {
	if t_inst == nil {
		panic("Cannot add nil Instruction!")
	}

	return this.list.PushFront(t_inst)
}

// Remove removes the given element from the list.
func (this *InstructionList) Remove(t_elem *l.Element) {
	this.list.Remove(t_elem)
}

// Front returns the first element of the list.
func (this *InstructionList) Front() *l.Element {
	return this.list.Front()
}

// Back returns the last element of the list.
func (this *InstructionList) Back() *l.Element {
	return this.list.Back()
}

func (this *InstructionList) Len() int {
	return this.list.Len()
}

func (this *InstructionList) String() string {
	text := "["

	var sep string
	for elem := this.list.Front(); elem != nil; elem = elem.Next() {
		value := InstructionListValue(elem)

		if value != nil {
			text += fmt.Sprintf("%s%v", sep, *value)
		} else {
			text += fmt.Sprintf("%s<nil>", sep)
		}

		sep = ", "
	}

	text += "]"

	return text
}

// Value returns the *Instruction stored in element t_elem.
// Panics if the element is nil or does not hold *Instruction.
func InstructionListValue(t_elem *l.Element) *Instruction {
	if t_elem == nil {
		panic("Element is nil!")
	}

	inst, ok := t_elem.Value.(*Instruction)
	if !ok {
		msg := fmt.Sprintf("Element value is not *Instruction but %T!", t_elem.Value)
		panic(msg)
	}

	return inst
}

func InitInstructionList() InstructionList {
	instList := InstructionList{list: l.New()}

	return instList
}
