package ast

import (
	"testing"

	"github.com/arthurlee945/monkey.on/token"
)

func TestString(t *testing.T) {
	program := Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "momonkey"},
					Value: "momonkey",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "apenation"},
					Value: "apenation",
				},
			},
		},
	}
	if program.String() != "let momonkey = apenation" {
		t.Errorf("program.String() is wrong. got=%q", program.String())
	}
}
