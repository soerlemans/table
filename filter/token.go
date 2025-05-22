package filter

import (
	u "github.com/soerlemans/table/util"
)

// TODO: Document.
type Token struct {
	Type  TokenType
	Value string
}

func InitToken(t_type TokenType, t_value string) Token {
	token := Token{t_type, t_value}
	defer u.LogStructName("InitToken", token, u.ETC80)

	return token
}
