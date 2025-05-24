package filter

import (
	"errors"

	a "github.com/soerlemans/table/filter/ast"
	u "github.com/soerlemans/table/util"
)

// Errors:
var (
	ErrNodeEmpty     = errors.New("Expected a node")
	ErrUnexpectedEos = errors.New("Unexpected end of stream")
	ErrExpectedToken = errors.New("Expected token")
)

// Import these as we use them frequently.
type NodeList a.NodeList
type Node a.Node

// TODO: Implement.
func errExpectedToken() {}

// Should be used in conjunction with defer.
func logUnlessNil(t_preabmle string, t_node Node) {
	if t_node != nil {
		u.Logf("Found %s: %T", t_preabmle, t_node)
	}
}

func accessorName(t_stream *TokenStream) (Node, error) {
	var node Node

	return node, nil
}

func accessorPositional(t_stream *TokenStream) (Node, error) {
	var node Node

	return node, nil
}

func column(t_stream *TokenStream) (Node, error) {
	var node Node

	// Check for a name based column accessor.
	if namePtr, err := accessorName(t_stream); node != nil {
		node = namePtr
	} else if err != nil {
		return node, err

	// Check for a positional based column accessor.
	} else if posPtr, err := accessorPositional(t_stream); node != nil {
		node = posPtr
	} else if err != nil {
		return node, err
	}

	return node, nil
}

func columnList(t_stream *TokenStream) (NodeList, error) {
	var list NodeList

	for {
		node, err := column(t_stream)
		if err != nil {
			return list, err
		}

		// If the node is not nil we found an item.
		if node != nil {
			list = append(list, node)
		} else {
			break
		}

		// We must check for the intermediary pipe symbol '|'.
		if !t_stream.Eos() {
			token := t_stream.Current()

			// If no intermediary pipe symbol was found we should quit.
			// TODO: Or maybe error if we are not at EOS and find a Pipe.
			if token.Type != COMMA {
				break
			}
		}

	}

	return list, nil
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
	}

	return node, nil
}

// Initialize the comparison.
func initComparison[T a.ComparisonType](t_stream *TokenStream, t_lhs Node) (Node, error) {
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
			// TODO: Throw error.
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

		// If the node is not nil we found an item.
		if node != nil {
			list = append(list, node)
		} else {
			break
		}

		// We must check for the intermediary pipe symbol '|'.
		if !t_stream.Eos() {
			token := t_stream.Current()

			// If no intermediary pipe symbol was found we should quit.
			// TODO: Or maybe error if we are not at EOS and find a Pipe.
			if token.Type != PIPE {
				break
			}
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
	u.Logf("BEGIN PARSING.")
	defer u.Logf("END PARSING.")

	list, err := program(t_stream)
	if err != nil {
		return list, err
	}

	return list, nil
}
