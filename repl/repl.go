package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/arthurlee945/monkey.on/lexer"
	"github.com/arthurlee945/monkey.on/parser"
	"github.com/arthurlee945/monkey.on/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

func StartParser(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		for _, stmt := range program.Statements {
			fmt.Printf("%+v\n", stmt.String())
		}

	}
}
