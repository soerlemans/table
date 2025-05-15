package filter

import (
	u "github.com/soerlemans/table/util"
)

type TokenType uint64
type TokenStream []TokenType

const (
	IDENTIFIER TokenType = iota
	STRING
	NUMBER

	// Function calls are identified by an identifier followed directly by a `(`.
	FUNCTION_CALL

	// This is an accessor using `.`.
	ACCESSOR_NAME

	// This is an accessor using `$`.
	ACCESSOR_POSITIONAL

	// This is the pipe that separates multiple expressions.
	PIPE

	// Comparison Operators:
	LESS_THAN
	LESS_THAN_EQUAL

	EQUAL
	NOT_EQUAL

	GREATER_THAN
	GREATER_THAN_EQUAL

	// Logical:
	NOT
	AND
	OR

	// Keywords:
	// Denotes a conditional expression that will determine if a row should be filtered out.
	// This is the default for every filter statement.
	WHEN

	// Denotes an operation which mutates something.
	// These expression do not filter out any rows.
	MUT

	// This operation will export the surviving rows to a markdown table.
	MD

	// This operation will export the surviving rows to JSON.
	JSON
)

// Lex the program text and return a TokenStream.
func Lex(t_text string) (TokenStream, error) {
	var stream TokenStream

	u.Logf("ProgramText: %s", t_text)

	// for i, rune_ := range t_text {

	// }

	return stream, nil
}
