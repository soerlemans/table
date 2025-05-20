package filter

import (
	"unicode"

	u "github.com/soerlemans/table/util"
)

type TokenType uint64
type TokenStream []TokenType

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

// Character Terminals:
const (
	DOUBLE_QUOTE_RN = '"'

	DOT_RN         = '.'
	DOLLAR_SIGN_RN = '$'

	PIPE_RN = '|'

	LESS_THAN_STR       = "<"
	LESS_THAN_EQUAL_STR = "<="

	EQUAL_STR     = "=="
	NOT_EQUAL_STR = "!="

	GREATER_THAN_STR       = ">"
	GREATER_THAN_EQUAL_STR = ">="
)

// TODO: Document.
type Token struct {
	Type  TokenType
	Value string
}

// func TokenType2Str(t_type TokenType) String {
// }

func lexIdentifier(t_text string) {}

// Lex the program text and return a TokenStream.
func Lex(t_text string) (TokenStream, error) {
	var stream TokenStream

	u.Logf("ProgramText: %s", t_text)

	for index, rn := range t_text {
		textView := t_text[index:]

		if unicode.IsSpace(rn) {
			u.Logln("Skipping whitespace.")
			continue
		} else if unicode.IsNumber(rn) {
			// TODO: Lex numbers.
		} else if rn == DOUBLE_QUOTE_RN {
			// TODO: Lex a string.
		} else if unicode.IsLetter(rn) {
			// Deal with possible function call.
			lexIdentifier(textView)
		} else if rn == DOT_RN {
			// TODO: Lex name accessor.
		} else if rn == DOLLAR_SIGN_RN {
			// TODO: Lex positional accessor.
		}
	}

	return stream, nil
}
