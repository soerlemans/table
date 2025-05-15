package filter

import (
	fn "github.com/soerlemans/table/filter/filter_node"
	u "github.com/soerlemans/table/util"
)

type Filter struct {
	Nodes []fn.FilterNode
}

// This creates our kind of AST thingy.
func initFilter(t_text string) (Filter, error)
{
	var filter_ Filter
	defer u.LogStructName("initFilter", filter_, u.ETC80)

	tokenStream, err := lex(t_text)
	if err != nil {
		return err
	}

	nodes, err := parse(tokenStream)
	if err != nil {
		return err
	}

	filter_ = Filter{nodes}

	return filter_
}
