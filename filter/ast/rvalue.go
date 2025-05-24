package ast

import (
	"strconv"
)

type String struct {
	Value string
}

type Number struct {
	Value int64
}

// TODO: Probably move this to lvalue some day.
type Identifier struct {
	Name string
}

func InitString(t_value string) (String, error) {
	return String{t_value}, nil
}

func InitNumber(t_value string) (Number, error) {
	var number Number
	// TODO: Log.
	// defer func(){}()

	integer, err := strconv.ParseInt(t_value, 10, 64)
	if err != nil {
		return number, err
	}

	number = Number{integer}

	return number, nil
}

func InitIdentifier(t_value string) (Identifier, error) {
	return Identifier{t_value}, nil
}
