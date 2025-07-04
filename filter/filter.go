package filter

import (
	u "github.com/soerlemans/table/util"
)

type Filter struct {
	Instructions InstListPtr
}

// This constructs the AST.
func InitFilter(t_text string) (Filter, error) {
	var filter_ Filter
	defer func() { u.LogStructName("InitFilter", filter_, u.ETC80) }()

	tokenStream, err := Lex(t_text)
	if err != nil {
		return filter_, err
	}

	instructions, err := Parse(&tokenStream)
	if err != nil {
		return filter_, err
	}


	filter_ = Filter{instructions}

	return filter_, nil
}
