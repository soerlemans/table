package filter

import (
	s "github.com/soerlemans/table/stream"
	u "github.com/soerlemans/table/util"
)

// TokenStream used later in parsing.
type TokenStream = s.SliceStream[Token]

// TokenType specification.
type TokenType uint64

const (
	Identifier TokenType = iota
	Number
	String

	// Function calls are identified by an identifier followed directly by a `(`.
	FunctionCall

	// This is an accessor using `.`.
	AccessorName

	// This is an accessor using `$`.
	AccessorPositional

	// Commas are required to separate expressions.
	Comma

	// This is the pipe that separates multiple expressions.
	Pipe

	// Little token for separating input on output specifiers.
	Colon

	// Comparison Operators:
	LessThan
	LessThanEqual

	Equal
	NotEqual

	GreaterThan
	GreaterThanEqual

	// Logical:
	Not
	And
	Or

	// Keywords:
	// Denotes a conditional expression that will determine if a row should be filtered out.
	// This is the default for every filter statement.
	When

	// Denotes an operation which mutates something.
	// These expression do not filter out any rows.
	Mut

	// Sort on a specific column.
	Sort
	NumericSort

	// Keyword for how many results you want maximum.
	Head
	Tail

	// Output specifiers:
	// Selects which columns to output (optionally specify order)).
	Csv

	// This operation will export the surviving rows to a markdown table.
	Md

	// This operation will export the surviving rows as a pretty table.
	Pretty

	// This operation will export the surviving rows to JSON.
	Json

	// This operation will export the surviving rows to HTML.
	Html
)

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case Number:
		return "Number"
	case String:
		return "String"
	case FunctionCall:
		return "FunctionCall"
	case AccessorName:
		return "."
	case AccessorPositional:
		return "$"
	case Comma:
		return ","
	case Pipe:
		return "|"
	case Colon:
		return ":"
	case LessThan:
		return "<"
	case LessThanEqual:
		return "<="
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	case GreaterThan:
		return ">"
	case GreaterThanEqual:
		return ">="
	case Not:
		return "!"
	case And:
		return "&&"
	case Or:
		return "||"
	case When:
		return "when"
	case Mut:
		return "mut"

	case Sort:
		return "sort"
	case NumericSort:
		return "numeric_sort"

	case Head:
		return "head"
	case Tail:
		return "tail"
	case Csv:
		return "csv"
	case Md:
		return "md"
	case Pretty:
		return "pretty"
	case Json:
		return "json"
	case Html:
		return "html"
	}

	// Optionally return an error?
	return "<Unknown TokenType>"
}

type Token struct {
	Type  TokenType
	Value string
}

func InitToken(t_type TokenType, t_value string) Token {
	token := Token{t_type, t_value}
	defer func() { u.LogStructName("InitToken", token, u.ETC80) }()

	return token
}
