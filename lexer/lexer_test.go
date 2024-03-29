package lexer

import (
	"testing"

	"github.com/arthurlee945/monkey.on/token"
)

type TestType struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let eight = 8;
	
	let add = fn(x, y){
		x + y;
	};
	
	let result = add(five, eight);
	!-/*8;
	8 < 25 > 8;

	if (5 < 10){
		return true;
	}else{
		return false;
	}
	
	10 == 10;
	8 != 10;

	let ff = 9.32;
	5%5;
	"monkey"
	"monkey paw"
	[1, 23];
	{"monkey" : "paw"}
	`

	//tests := []struct{expectedType token.TokenType expectedLiteral string}
	tests := []TestType{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "eight"},
		{token.ASSIGN, "="},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "eight"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},

		{token.INT, "8"},
		{token.LT, "<"},
		{token.INT, "25"},
		{token.GT, ">"},
		{token.INT, "8"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "8"},
		{token.NOT_EQ, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ff"},
		{token.ASSIGN, "="},
		{token.FLOAT, "9.32"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.MODULO, "%"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.STRING, "monkey"},
		{token.STRING, "monkey paw"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "23"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "monkey"},
		{token.COLON, ":"},
		{token.STRING, "paw"},
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tType := range tests {
		tok := l.NextToken()

		if tok.Type != tType.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q || got=%q", i, tType.expectedType, tok.Type)
		}

		if tok.Literal != tType.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q || got=%q", i, tType.expectedLiteral, tok.Literal)
		}
	}
}
