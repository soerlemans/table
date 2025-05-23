package filter

import (
	fn "github.com/soerlemans/table/filter/filter_node"
	u "github.com/soerlemans/table/util"
)

type Filter struct {
	Nodes []fn.FilterNode
}

// This creates our kind of AST thingy.
func InitFilter(t_text string) (Filter, error) {
	var filter_ Filter
	defer func() { u.LogStructName("initFilter", filter_, u.ETC80) }()

	tokenVec, err := Lex(t_text)
	if err != nil {
		return filter_, err
	}

	nodes, err := Parse(tokenVec)
	if err != nil {
		return filter_, err
	}

	filter_ = Filter{nodes}

	return filter_, nil
}
