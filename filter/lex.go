package filter

import (
	"errors"
	"fmt"
	"unicode"

	s "github.com/soerlemans/table/stream"
	u "github.com/soerlemans/table/util"
)

// Character Terminals:
const (
	DoubleQuoteRn = '"'

	DotRn        = '.'
	DollarSignRn = '$'

	CommaRn = ','
	PipeRn  = '|'
	ColonRn = ':'

	LessThanStr      = "<"
	LessThanEqualStr = "<="

	EqualStr    = "=="
	NotEqualStr = "!="

	GreaterThanStr      = ">"
	GreaterThanEqualStr = ">="
)

// Single rune symbols mapped to TokenType.
var singleRuneSymbols = map[string]TokenType{
	string(DotRn):        AccessorName,
	string(DollarSignRn): AccessorPositional,
	string(CommaRn):      Comma,
	string(PipeRn):       Pipe,
	string(ColonRn):      Colon,

	LessThanStr:    LessThan,
	GreaterThanStr: GreaterThan,
}

// Multi rune symbols mapped to TokenType.
var multiRuneSymbols = map[string]TokenType{
	LessThanEqualStr: LessThanEqual,

	EqualStr:    Equal,
	NotEqualStr: NotEqual,

	GreaterThanEqualStr: GreaterThanEqual,
}

// Keywords mapped to TokenType.
var Keywords = map[string]TokenType{
	// Keywords:
	"when": When,
	"mut":  Mut,

	"sort":         Sort,
	"numeric_sort": NumericSort,

	"head": Head,
	"tail": Tail,

	// Writer specifications:
	"csv":    Csv,
	"md":     Md,
	"pretty": Pretty,
	"json":   Json,
	"html":   Html,
}

// Errors:
var (
	ErrInvalidRune           = errors.New("Invalid rune")
	ErrIncorrectStartingRune = errors.New("Incorrect starting rune")
	ErrUnterminated          = errors.New("Unterminated")
)

func incorrectStartingRune(t_preamble string, t_rn rune) error {
	return fmt.Errorf("%s do not start with '%c' (%w).", t_preamble, t_rn, ErrIncorrectStartingRune)
}

func unterminated(t_preamble string, t_rn rune) error {
	return fmt.Errorf("unterminated %s expected '%c' (%w).", t_preamble, t_rn, ErrUnterminated)
}

// func TokenType2Str(t_type TokenType) String {
// }

func lexNumbers(t_stream *s.StringStream) (Token, error) {
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

		token = InitToken(Number, value)
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

func lexIdentifier(t_stream *s.StringStream) (Token, error) {
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

		token = InitToken(Identifier, value)
	} else {
		err := incorrectStartingRune("identifiers", initialRn)

		return token, err
	}

	// Check if we have an identifier that is a keyword.
	checkKeyword(&token)

	return token, nil

}

func lexString(t_stream *s.StringStream) (Token, error) {
	var token Token
	initialRn := t_stream.Current()

	if initialRn == DoubleQuoteRn {
		var value string
		for !t_stream.Eos() {
			// Inch stream forward.
			t_stream.Next()

			// Error handle.
			if t_stream.Eos() {
				err := unterminated("string", DoubleQuoteRn)

				return token, err
			}

			// Get current rune.
			rn := t_stream.Current()
			if rn != DoubleQuoteRn {
				value += string(rn)
			} else {
				// Pass by the double quote rune.
				t_stream.Next()

				// Double quote so end of string.
				break
			}
		}

		token = InitToken(String, value)
	} else {
		return token, incorrectStartingRune("strings", initialRn)
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

func lexSymbol(t_stream *s.StringStream) (Token, bool) {
	var (
		token       Token
		foundSymbol bool // Defaults to false.
	)

	// Handle potential single rune symbols.
	rn := t_stream.Current()
	buf := string(rn)
	token, foundSymbol = lexSingleSymbol(&buf)

	// Handle potential multi rune symbols.
	if !t_stream.Eos() {
		rn, ok := t_stream.Peek()
		if ok {
			buf += string(rn)
			tokenMulti, foundMultiSymbol := lexMultiSymbol(&buf)

			if foundMultiSymbol {
				token = tokenMulti
				foundSymbol = foundMultiSymbol

				// Skip to the next char.
				t_stream.Next()
			}
		}
	}

	if foundSymbol {
		// Skip another character either way.
		t_stream.Next()
	}

	return token, foundSymbol
}

// Lex the program text and return a TokenVec.
func Lex(t_text string) (TokenStream, error) {
	tokenStream := s.InitSliceStreamEmpty[Token]()
	defer func() { u.Logf("tokenStream: %v", tokenStream.View) }()

	u.Logf("ProgramText: \"%s\"", t_text)

	runeStream := s.InitStringStream(&t_text)

	for !runeStream.Eos() {
		rn := rune(runeStream.Current())

		if unicode.IsSpace(rn) {
			u.Logln("Skipping whitespace.")
			runeStream.Next()
			continue
		} else if unicode.IsNumber(rn) {
			token, err := lexNumbers(&runeStream)
			if err != nil {
				return tokenStream, err
			}

			tokenStream.Append(token)
		} else if rn == DoubleQuoteRn {
			// TODO: Lex a string.
			token, err := lexString(&runeStream)
			if err != nil {
				return tokenStream, err
			}

			tokenStream.Append(token)
		} else if unicode.IsLetter(rn) {
			// Deal with possible function call.
			token, err := lexIdentifier(&runeStream)
			if err != nil {
				return tokenStream, err
			}

			tokenStream.Append(token)
		} else {
			token, found := lexSymbol(&runeStream)

			// Error handle, unhandled tokens:
			if !found {
				u.Logf("Unhandled rune: %c", rn)

				err := fmt.Errorf("Invalid rune for lexing '%c' (%w).", rn, ErrInvalidRune)
				return tokenStream, err
			}

			tokenStream.Append(token)
		}
	}

	return tokenStream, nil
}
