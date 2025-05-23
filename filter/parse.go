package filter

import (
	"errors"

	fn "github.com/soerlemans/table/filter/filter_node"
	// u "github.com/soerlemans/table/util"
)

// Errors:
var (
	ErrNodeEmpty = errors.New("")
)

// Import these as we use them frequently.
type NodeListPtr fn.NodeListPtr
type NodePtr fn.NodePtr

func item() (NodePtr, error) {
	var node NodePtr

	return node, nil
}

func itemList(t_stream *TokenStream) (NodeListPtr, error) {
	var list NodeListPtr

	for {
		nodePtr, err := item()
		if err != nil {
			return list, err
		}

		if nodePtr != nil {
			list = append(list, nodePtr)
		} else {
			break
		}
	}

	return list, nil
}

// This function is here purely just to match the grammary.yy.
func program(t_stream *TokenStream) (NodeListPtr, error) {
	return itemList(t_stream)
}

// Source code to parse.
func Parse(t_stream *TokenStream) (NodeListPtr, error) {
	list, err := program(t_stream)
	if err != nil {
		return list, err
	}

	return list, nil
}
