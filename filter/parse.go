package filter

import (
	"errors"

	fn "github.com/soerlemans/table/filter/filter_node"
	// u "github.com/soerlemans/table/util"
)

// Errors:
var (
	ErrNodeEmpty     = errors.New("Expected a node")
	ErrUnexpectedEos = errors.New("Unexpected end of stream")
	ErrExpectedToken = errors.New("Expected token")
)

// Import these as we use them frequently.
type NodeList fn.NodeList
type Node fn.Node

// TODO: Implement.
func errExpectedToken() {}

func keyword(t_stream *TokenStream) (Node, error) {
	var node Node

	return node, nil
}

func rvalue(t_stream *TokenStream) (Node, error) {
	var node Node

	return node, nil
}

// Initialize the comparison.
func initComparison[T fn.ComparisonType](t_stream *TokenStream, t_lhs Node) (Node, error) {
	var node Node

	t_stream.Next()
	if t_stream.Eos() {
		// TODO: Throw error.
	}

	rhs, err := rvalue(t_stream)
	if err != nil {
		return node, err
	}

	if rhs == nil {
		// TODO: Error expected a node.
		return node, nil
	}

	comp := fn.InitComparison[T](t_lhs, rhs)
	node = &comp

	return node, nil
}

func expr(t_stream *TokenStream) (Node, error) {
	var node Node

	lhs, err := rvalue(t_stream)
	if err != nil {
		return node, err
	}

	if lhs != nil {
		if t_stream.Eos() {
			// TODO: Throw error.
		}

		token := t_stream.Current()

		switch token.Type {
		case LESS_THAN:

			break

		case LESS_THAN_EQUAL:
			break

		case EQUAL:
			break

		case NOT_EQUAL:
			break

		case GREATER_THAN:
			break

		case GREATER_THAN_EQUAL:
			break
		}
	}

	return node, nil
}

func item(t_stream *TokenStream) (Node, error) {
	var node Node

	// Check for keywords.
	if keywordPtr, err := keyword(t_stream); node != nil {
		node = keywordPtr
	} else if err != nil {
		return node, err

		// Check for an epxression.
	} else if exprPtr, err := expr(t_stream); node != nil {
		node = exprPtr
	} else if err != nil {
		return node, err
	}

	return node, nil
}

func itemList(t_stream *TokenStream) (NodeList, error) {
	var list NodeList

	for {
		node, err := item(t_stream)
		if err != nil {
			return list, err
		}

		if node != nil {
			list = append(list, node)
		} else {
			break
		}
	}

	return list, nil
}

// This function is here purely just to match the grammary.yy.
func program(t_stream *TokenStream) (NodeList, error) {
	return itemList(t_stream)
}

// Source code to parse.
func Parse(t_stream *TokenStream) (NodeList, error) {
	list, err := program(t_stream)
	if err != nil {
		return list, err
	}

	return list, nil
}
