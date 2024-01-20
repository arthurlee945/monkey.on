package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arthurlee945/monkey.on/ast"
	"github.com/arthurlee945/monkey.on/lexer"
)

// LET TEST
func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 8348583;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statments does not contain 3 statements. got=%d", len(program.Statements))
	}

	test := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range test {
		stmt := program.Statements[i]
		if !testLetStatements(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s can not be casted to type *ast.LetStatement. got=%q", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

// RETURN TEST
func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 9302;
	return 08;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement, got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got=%q", returnStmt.TokenLiteral())
		}
	}
}

// IDENTIFIER TEST
func TestIdentifierExpression(t *testing.T) {
	input := "momono;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program expected 1 statement but got - %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.TokenLiteral() != "momono" {
		t.Fatalf("ident.TokenLiteral is not %s. got=%s", "momono", ident.TokenLiteral())
	}
}

// INTEGER LITERAL TEST
func TestIntegerLiteralExpression(t *testing.T) {
	input := `8;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. Expected 1 but got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expression is not ast.IntegerLiteral, got=%T", stmt.Expression)
	}
	if literal.Value != 8 {
		t.Errorf("literal.Value not %d, got=%d", 8, literal.Value)
	}
	if literal.TokenLiteral() != "8" {
		t.Errorf("literal.TokenLiteral not %s, got=%s", "8", literal.TokenLiteral())
	}
}

// Float LITERAL TEST
func TestFloatLiteralExpression(t *testing.T) {
	input := `8.43;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. Expected 1 but got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("Expression is not ast.FloatLiteral, got=%T", stmt.Expression)
	}
	if literal.Value != 8.43 {
		t.Errorf("literal.Value not %f, got =%f", 8.43, literal.Value)
	}
	if literal.TokenLiteral() != "8.43" {
		t.Errorf("literal.TokenLiteral not %s, got=%s", "8.43", literal.TokenLiteral())
	}
}

// Parsing Prefix Expression
type PrefixTest[T int64 | float64] struct {
	input       string
	operator    string
	numberValue T
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixIntTests, prefixFlaotTests := []PrefixTest[int64]{
		{"!8", "!", 8},
		{"-23", "-", 23},
	}, []PrefixTest[float64]{
		{"!7.23", "!", 7.23},
		{"-84.3", "-", 84.3},
	}

	testParsingPrefixExpressions[int64](t, prefixIntTests)
	testParsingPrefixExpressions[float64](t, prefixFlaotTests)
}
func testParsingPrefixExpressions[T int64 | float64](t *testing.T, pt []PrefixTest[T]) {
	for _, tt := range pt {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d Statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got=%s", tt.operator, exp.Operator)
		}
		if !testNumberLiteral[T](t, exp.Right, tt.numberValue) {
			return
		}
	}
}
func testNumberLiteral[T int64 | float64](t *testing.T, nl ast.Expression, value T) bool {
	var num ast.Expression
	var okNum bool
	var expression string

	if _, ok := any(value).(int64); ok {
		num, okNum = nl.(*ast.IntegerLiteral)
		expression = "IntegerLiteral"
	} else {
		num, okNum = nl.(*ast.FloatLiteral)
		expression = "FloatLiteral"
	}

	if !okNum {
		t.Errorf("nl not *ast.%s, got=%T", expression, nl)
		return false
	}

	if expression == "IntegerLiteral" && (num.(*ast.IntegerLiteral)).Value != int64(value) {
		t.Errorf("num.Value not %d, got=%d", int64(value), (num.(*ast.IntegerLiteral)).Value)
		return false
	} else if expression == "FloatLiteral" && (num.(*ast.FloatLiteral)).Value != float64(value) {
		t.Errorf("num.Value not %f, got=%f", float64(value), (num.(*ast.FloatLiteral)).Value)
		return false
	}

	if expression == "IntegerLiteral" && num.TokenLiteral() != fmt.Sprintf("%d", int(value)) {
		t.Errorf("num.TokenLiteral not %d. got=%s", int(value), num.TokenLiteral())
		return false
	} else if expression == "FloatLiteral" && !strings.HasPrefix(fmt.Sprintf("%f", float64(value)), num.TokenLiteral()) {
		t.Errorf("num.TokenLiteral not %f. got=%s", float64(value), num.TokenLiteral())
		return false
	}

	return true

}

// UTILITY
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser encountered %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err)
	}

	t.FailNow()
}
