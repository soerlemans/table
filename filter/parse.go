package filter

import (
	"errors"
	"fmt"
	"strconv"

	a "github.com/soerlemans/table/filter/ast"
	u "github.com/soerlemans/table/util"
)

// Errors:
var (
	ErrNodeEmpty     = errors.New("Expected a node")
	ErrUnexpectedEos = errors.New("Unexpected end of stream")
	ErrExpected      = errors.New("Expected")
	ErrUnsupported   = errors.New("Unsupported")
)

// Import these as we use them frequently.
type NodeList a.NodeList
type Node a.Node

type parseFnNode = func(*TokenStream) (Node, error)
type parseFnNodeList = func(*TokenStream) (NodeList, error)

// TODO: Implement.
func errNodeEmpty(t_location string) error {
	return fmt.Errorf("Expected node returned nil in %s (%w)", t_location, ErrNodeEmpty)
}

func errExpectedString(t_location string, t_string string) error {
	return fmt.Errorf("Expected a \"%s\" in %s (%w)", t_string, t_location, ErrExpected)
}

func errExpectedTokenType(t_location string, t_type TokenType) error {
	return errExpectedString(t_location, t_type.String())
}

func errUnexpectedEos(t_location string) error {
	return fmt.Errorf("Unexpected end of stream in %s (%w)", t_location, ErrUnexpectedEos)
}

// Should be used in conjunction with defer.
func logUnlessNil(t_preabmle string, t_node Node) {
	if t_node != nil {
		u.Logf("Found %s: %T", t_preabmle, t_node)
	}
}

// Check if a node is initialized if the node is not nil.
// And no error is set.
func validNode(t_node Node, t_err error) bool {
	return t_node != nil && t_err == nil
}

func parseList(t_stream *TokenStream, t_fn parseFnNode, t_sep TokenType) (NodeList, error) {
	var list NodeList

	for {
		node, err := t_fn(t_stream)
		if err != nil {
			return list, err
		}

		// If the node is not nil we found an item.
		if node != nil {
			list = append(list, node)
		} else {
			break
		}

		// We must check for any separator symbols.
		if !t_stream.Eos() {
			token := t_stream.Current()

			// If no intermediary pipe symbol was found we should quit.
			if token.Type != t_sep {
				break
			}
		}

	}

	return list, nil
}

// Parsing functions:
func accessorName(t_stream *TokenStream) (Node, error) {
	var node Node
	defer func() { logUnlessNil("accessorName", node) }()

	// Guard clause.
	if t_stream.Eos() {
		return node, nil
	}

	dot := t_stream.Current()

	// TODO: Refactor boilerplate later.
	if dot.Type == ACCESSOR_NAME {
		u.Logf("Found a '%v'", dot.Type)
		t_stream.Next()

		if t_stream.Eos() {
			return node, errUnexpectedEos("accessorName")
		}

		// Identifier and string are both converted to the same value.
		token := t_stream.Current()
		switch token.Type {
		case IDENTIFIER:
			fallthrough
		case STRING:
			node, err := a.InitAccessorName(token.Value)
			if err != nil {
				return node, nil
			}
			break

		default:
			return node, errExpectedString("accessorName", "identifier or string")
		}

	}

	return node, nil
}

func accessorPositional(t_stream *TokenStream) (Node, error) {
	var node Node
	defer func() { logUnlessNil("accessorName", node) }()

	// Guard clause.
	if t_stream.Eos() {
		return node, nil
	}

	dollarSign := t_stream.Current()

	// TODO: Refactor boilerplate later.
	if dollarSign.Type == ACCESSOR_POSITIONAL {
		u.Logf("Found a '%v'", dollarSign.Type)
		t_stream.Next()

		if t_stream.Eos() {
			return node, errUnexpectedEos("accessorPositional")
		}

		// Identifier value must be used to convert to a string later.
		// Possibly add expressions or arithmetic expressions later down the line.
		token := t_stream.Current()
		switch token.Type {
		case IDENTIFIER:
			return node, nil

		case NUMBER:
			integer, err := strconv.ParseInt(token.Value, 10, 64)
			if err != nil {
				return node, nil
			}

			positionalNode, err := a.InitAccessorPositional(integer)
			if err != nil {
				return node, nil
			}

			node = &positionalNode
			break

		default:
			return node, errExpectedString("accessorPositional", "identifier or number")
		}

	}

	return node, nil
}

func column(t_stream *TokenStream) (Node, error) {
	var node Node

	// Check for a name based column accessor.
	if nameNode, err := accessorName(t_stream); validNode(nameNode, err) {
		node = &nameNode
	} else if err != nil {
		return node, err

		// Check for a positional based column accessor.
	} else if posNode, err := accessorPositional(t_stream); validNode(posNode, err) {
		node = &posNode
	} else {
		return node, err
	}

	return node, nil
}

func parameter(t_stream *TokenStream) (Node, error) {
	return rvalue(t_stream)
}

func parameterList(t_stream *TokenStream) (NodeList, error) {
	return parseList(t_stream, parameter, COMMA)
}

// Ast:
func keyword(t_stream *TokenStream) (Node, error) {
	var node Node

	return node, nil
}

func rvalue(t_stream *TokenStream) (Node, error) {
	var node Node
	defer func() { logUnlessNil("rvalue", node) }()

	// Guard clause.
	if t_stream.Eos() {
		return node, nil
	}

	token := t_stream.Current()

	// TODO: Refactor boilerplate later.
	switch token.Type {
	case NUMBER:
		number, err := a.InitNumber(token.Value)
		if err != nil {
			return node, err
		}

		node = &number
		t_stream.Next()
		break

	case STRING:
		str, err := a.InitString(token.Value)
		if err != nil {
			return node, err
		}

		node = &str
		t_stream.Next()
		break

	case IDENTIFIER:
		id, err := a.InitIdentifier(token.Value)
		if err != nil {
			return node, err
		}

		node = &id
		t_stream.Next()
		break

	default:
		columnNode, err := column(t_stream)
		if err != nil {
			return node, err
		}

		node = &columnNode
		t_stream.Next()
		break
	}

	return node, nil
}

// Initialize the comparison.
func initComparison[T a.ComparisonType](t_stream *TokenStream, t_lhs Node) (Node, error) {
	var node Node

	t_stream.Next()
	if t_stream.Eos() {
		return node, errUnexpectedEos("initComparison")
	}

	rhs, err := rvalue(t_stream)
	if err != nil {
		return node, err
	}

	if rhs == nil {
		return node, errExpectedString("initComparison", "right hand side of expression")
	}

	comp := a.InitComparison[T](t_lhs, rhs)
	node = &comp

	return node, nil
}

func expr(t_stream *TokenStream) (Node, error) {
	var node Node
	defer func() { logUnlessNil("expr", node) }()

	lhs, err := rvalue(t_stream)
	if err != nil {
		return node, err
	}

	if lhs != nil {
		if t_stream.Eos() {
			return node, errUnexpectedEos("expr")
		}

		token := t_stream.Current()

		// TODO: Cleanup boilerplate.
		switch token.Type {
		case LESS_THAN:
			ltNode, err := initComparison[a.LessThan](t_stream, lhs)
			if err != nil {
				return node, err
			}

			node = &ltNode
			break

		case LESS_THAN_EQUAL:
			lteNode, err := initComparison[a.LessThanEqual](t_stream, lhs)
			if err != nil {
				return node, err
			}

			node = &lteNode
			break

		case EQUAL:
			eqNode, err := initComparison[a.Equal](t_stream, lhs)
			if err != nil {
				return node, err
			}

			node = &eqNode
			break

		case NOT_EQUAL:
			neNode, err := initComparison[a.NotEqual](t_stream, lhs)
			if err != nil {
				return node, err
			}

			node = &neNode
			break

		case GREATER_THAN:
			gtNode, err := initComparison[a.GreaterThan](t_stream, lhs)
			if err != nil {
				return node, err
			}

			node = &gtNode
			break

		case GREATER_THAN_EQUAL:
			gteNode, err := initComparison[a.GreaterThanEqual](t_stream, lhs)
			if err != nil {
				return node, err
			}

			node = &gteNode
			break
		}
	}

	return node, nil
}

func item(t_stream *TokenStream) (Node, error) {
	var node Node

	// Check for keywords.
	if keywordNode, err := keyword(t_stream); validNode(keywordNode, err) {
		node = &keywordNode
	} else if err != nil {
		return node, err

		// Check for an epxression.
	} else if exprNode, err := expr(t_stream); validNode(exprNode, err) {
		node = &exprNode
	} else if err != nil {
		return node, err
	}

	return node, nil
}

func itemList(t_stream *TokenStream) (NodeList, error) {
	return parseList(t_stream, item, PIPE)
}

// This function is here purely just to match the grammary.yy.
func program(t_stream *TokenStream) (NodeList, error) {
	return itemList(t_stream)
}

// Source code to parse.
func Parse(t_stream *TokenStream) (NodeList, error) {
	u.Logf("BEGIN PARSING.")
	defer u.Logf("END PARSING.")

	list, err := program(t_stream)
	if err != nil {
		return list, err
	}

	return list, nil
}
