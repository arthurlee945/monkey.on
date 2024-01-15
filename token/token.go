package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
type testing struct {
	name string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifier + literals(Func, Var, etc...)
	IDENT = "IDENT" //add, foo, bar...
	INT   = "INT"   //1,2,3 ...

	//OPERATOR
	ASSIGN = "="
	PLUS   = "+"

	//Delimiter
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//Keywords
	FUNCTION = "fUNCTION"
	LET      = "LET"
)
