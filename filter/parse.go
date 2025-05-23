package filter

import (
	fn "github.com/soerlemans/table/filter/filter_node"
	// u "github.com/soerlemans/table/util"
)

type NodeList []fn.FilterNode
type Node fn.FilterNode

func item() (Node, error) {}

func itemList() (NodeList, error) {
	for {
	}
}

func program(t_stream *TokenStream) (NodeList, error) {
	var ast NodeList

	return ast, nil
}

// Source code to parse.
func Parse(t_stream *TokenStream) (NodeList, error) {

	return nil, nil
}
