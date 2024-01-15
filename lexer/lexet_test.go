package lexer

import (
	"fmt"

	token "github.com/arthurlee945/monkey.on/token"
)

func test() {
	t := token.Token{Type: "something", Literal: "nana"}
	fmt.Println(t)
}
