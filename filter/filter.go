package filter

import (
	u "github.com/soerlemans/table/util"
)

type Filter struct {
	Program NodeList
}

// This constructs the AST.
func InitFilter(t_text string) (Filter, error) {
	var filter_ Filter
	defer func() { u.LogStructName("initFilter", filter_, u.ETC80) }()

	tokenStream, err := Lex(t_text)
	if err != nil {
		return filter_, err
	}

	nodes, err := Parse(&tokenStream)
	if err != nil {
		return filter_, err
	}

	filter_ = Filter{nodes}

	return filter_, nil
}
