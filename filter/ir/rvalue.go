package ir

import (
	"strconv"
)

type String struct {
	Value string
}

func InitString(t_value string) (String, error) {
	return String{t_value}, nil
}

type Number struct {
	Value int64
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

// TODO: Probably move this to lvalue some day.
type Identifier struct {
	Name string
}

func InitIdentifier(t_value string) (Identifier, error) {
	return Identifier{t_value}, nil
}
