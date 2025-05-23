package filter

import (
	s "github.com/soerlemans/table/stream"
	u "github.com/soerlemans/table/util"
)

// TokenStream used later in parsing.
type TokenStream s.SliceStream[Token]

// TokenType specification.
type TokenType uint64
const (
	IDENTIFIER TokenType = iota
	NUMBER
	STRING

	// Function calls are identified by an identifier followed directly by a `(`.
	FUNCTION_CALL

	// This is an accessor using `.`.
	ACCESSOR_NAME

	// This is an accessor using `$`.
	ACCESSOR_POSITIONAL

	// Commas are required to separate expressions.
	COMMA

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

	// Output specifiers:
	// Selects which columns to output (optionally specify order)).
	OUT

	// This operation will export the surviving rows to a markdown table.
	MD

	// This operation will export the surviving rows to JSON.
	JSON
)

// TODO: Document.
type Token struct {
	Type  TokenType
	Value string
}

func InitToken(t_type TokenType, t_value string) Token {
	token := Token{t_type, t_value}
	defer func() { u.LogStructName("InitToken", token, u.ETC80) }()

	return token
}
