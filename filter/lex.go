package filter

import (
	"errors"
	"fmt"
	"unicode"

	u "github.com/soerlemans/table/util"
)

type TokenType uint64
type TokenVec []Token

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

// Single rune symbols mapped to TokenType.
var singleRuneSymbols = map[string]TokenType{
	string(DOT_RN):         ACCESSOR_NAME,
	string(DOLLAR_SIGN_RN): ACCESSOR_POSITIONAL,
	string(PIPE_RN):        PIPE,
}

// Multi rune symbols mapped to TokenType.
var multiRuneSymbols = map[string]TokenType{
	LESS_THAN_STR:       LESS_THAN,
	LESS_THAN_EQUAL_STR: LESS_THAN_EQUAL,

	EQUAL_STR:     EQUAL,
	NOT_EQUAL_STR: NOT_EQUAL,

	GREATER_THAN_STR:       GREATER_THAN,
	GREATER_THAN_EQUAL_STR: GREATER_THAN_EQUAL,
}

// Keywords mapped to TokenType.
var Keywords = map[string]TokenType{
	"when": WHEN,
	"mut":  MUT,
	"out":  OUT,
	"md":   MD,
	"json": JSON,
}

// Errors:
var (
	ErrInvalidRune           = errors.New("Invalid rune")
	ErrIncorrectStartingRune = errors.New("Incorrect starting rune")
	ErrUnterminated          = errors.New("Unterminated")
)

func incorrectStartingRune(t_preamble string, t_rn rune) error {
	return fmt.Errorf("%s do not start with '%c'(%w).", t_preamble, t_rn, ErrIncorrectStartingRune)
}

func unterminated(t_preamble string, t_rn rune) error {
	return fmt.Errorf("unterminated %s expected '%c'(%w).", t_rn, ErrUnterminated)

}

// func TokenType2Str(t_type TokenType) String {
// }

func lexNumbers(t_stream *Stream) (Token, error) {
	var token Token
	initialRn := t_stream.Current()

	if unicode.IsNumber(initialRn) {
		var value string
		for !t_stream.Eos() {
			// Get current rune.
			rn := t_stream.Current()
			if unicode.IsNumber(rn) {
				value += string(rn)
			} else {
				// Double quote so end of string.
				break
			}

			// Inch stream forward.
			t_stream.Next()

			// Error handle.
			if t_stream.Eos() {
				// If we run out of numbers its fine.
				break
			}
		}

		token = InitToken(STRING, value)
	} else {
		err := incorrectStartingRune("numbers", initialRn)

		return token, err
	}

	return token, nil

}

func checkKeyword(t_token *Token) {
	tokenType, ok := Keywords[t_token.Value]
	if ok {
		u.Logf("Keyword found: %s", t_token.Value)

		t_token.Type = tokenType
	}
}

func lexIdentifier(t_stream *Stream) (Token, error) {
	var token Token
	initialRn := t_stream.Current()

	if unicode.IsLetter(initialRn) {
		var value string
		for !t_stream.Eos() {
			// Get current rune.
			rn := t_stream.Current()
			validIdentifierRune := unicode.IsLetter(rn) || unicode.IsNumber(rn) || rn == '_'
			if validIdentifierRune {
				value += string(rn)
			} else {
				// No more alphanumerics means end of string.
				break
			}

			// Inch stream forward.
			t_stream.Next()

			// Error handle.
			if t_stream.Eos() {
				// If we run out of alphanumerics its fine.
				break
			}

		}

		token = InitToken(STRING, value)
	} else {
		err := incorrectStartingRune("identifiers", initialRn)

		return token, err
	}

	// Check if we have an identifier that is a keyword.
	checkKeyword(&token)

	return token, nil

}

func lexString(t_stream *Stream) (Token, error) {
	var token Token
	initialRn := t_stream.Current()

	if initialRn == DOUBLE_QUOTE_RN {
		var value string
		for !t_stream.Eos() {
			// Inch stream forward.
			t_stream.Next()

			// Error handle.
			if t_stream.Eos() {
				err := unterminated("string", DOUBLE_QUOTE_RN)

				return token, err
			}

			// Get current rune.
			rn := t_stream.Current()
			if rn != DOUBLE_QUOTE_RN {
				value += string(rn)
			} else {
				// Double quote so end of string.
				break
			}
		}

		token = InitToken(STRING, value)
	} else {
		err := incorrectStartingRune("strings", initialRn)

		return token, err
	}

	return token, nil
}

func lexSingleSymbol(t_text *string) (Token, bool) {
	var token Token

	text := *t_text
	tokenType, ok := singleRuneSymbols[text]

	if ok {
		u.Logf("Found a single rune symbol: %s", text)
		token = InitToken(tokenType, text)

		return token, true
	}

	return token, false
}

func lexMultiSymbol(t_text *string) (Token, bool) {
	var token Token

	text := *t_text
	tokenType, ok := multiRuneSymbols[text]

	if ok {
		u.Logf("Found a multi rune symbol: %s", text)
		token = InitToken(tokenType, text)

		return token, true
	}

	return token, false

}

func lexSymbol(t_stream *Stream) (Token, bool) {
	var (
		token       Token
		foundSymbol bool // Defaults to false.
	)

	rn := t_stream.Current()
	buf := string(rn)
	token, foundSymbol = lexSingleSymbol(&buf)

	if !foundSymbol {
		rn, ok := t_stream.Peek()

		if ok {
			buf += string(rn)

			token, foundSymbol = lexMultiSymbol(&buf)

			// Skip to the next char.
			t_stream.Next()
			t_stream.Next()
		}
	}

	return token, foundSymbol
}

// Lex the program text and return a TokenVec.
func Lex(t_text string) (TokenVec, error) {
	var tokenVec TokenVec

	u.Logf("ProgramText: %s", t_text)

	runeStream := initStream(&t_text)

	for !runeStream.Eos() {
		rn := runeStream.Current()

		if unicode.IsSpace(rn) {
			u.Logln("Skipping whitespace.")
			runeStream.Next()
			continue
		} else if unicode.IsNumber(rn) {
			token, err := lexNumbers(&runeStream)
			if err != nil {
				return tokenVec, err
			}

			tokenVec = append(tokenVec, token)
		} else if rn == DOUBLE_QUOTE_RN {
			// TODO: Lex a string.
			token, err := lexString(&runeStream)
			if err != nil {
				return tokenVec, err
			}

			tokenVec = append(tokenVec, token)
		} else if unicode.IsLetter(rn) {
			// Deal with possible function call.
			token, err := lexIdentifier(&runeStream)
			if err != nil {
				return tokenVec, err
			}

			tokenVec = append(tokenVec, token)
		} else {
			token, found := lexSymbol(&runeStream)
			tokenVec = append(tokenVec, token)

			// Error handle unhandled tokens:
			if !found {
				u.Logf("Unhandled rune: %c", rn)

				err := fmt.Errorf("Invalid rune for lexing '%c' (%w).", rn, ErrInvalidRune)
				return tokenVec, err
			}
		}
	}

	return tokenVec, nil
}
