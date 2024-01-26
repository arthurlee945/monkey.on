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
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let monkeyWheel = round", "monkeyWheel", "round"},
	}

	for _, tt := range tests {
		program := prepTest(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statments does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]

		if !testLetStatements(t, stmt, tt.expectedIdent) {
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

// RETURN TEST
func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return monkeyPaw;", "monkeyPaw"},
	}

	for _, tt := range tests {
		program := prepTest(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statments does not contain %d statements. got=%d", 1, len(program.Statements))
		}
		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

// IDENTIFIER TEST
func TestIdentifierExpression(t *testing.T) {
	input := "momono;"

	stmt := prepExpressionTest(t, input)

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.TokenLiteral() != "momono" {
		t.Fatalf("ident.TokenLiteral is not %s. got=%s", "momono", ident.TokenLiteral())
	}
}

// BOOLEAN TEST
func TestBooleanExpression(t *testing.T) {
	input := "true;"

	stmt := prepExpressionTest(t, input)

	boolExp, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Boolean. got=%T", stmt.Expression)
	}

	if boolExp.TokenLiteral() != "true" {
		t.Fatalf("boolExp.TokenLiteral is not %s. got=%s", "true", boolExp.TokenLiteral())
	}
}

// INTEGER LITERAL TEST
func TestIntegerLiteralExpression(t *testing.T) {
	input := `8;`

	stmt := prepExpressionTest(t, input)

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

	stmt := prepExpressionTest(t, input)

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

// STRING LITERAL TEST
func TestStringLiteralExpression(t *testing.T) {
	input := `"monkey world"`

	stmt := prepExpressionTest(t, input)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "monkey world" {
		t.Fatalf("literal.Value is not %q. got=%q", "monkey world", literal.Value)
	}

}

// ARRAY PARSE TEST
func TestParsingArrayLiteral(t *testing.T) {
	input := `[1,  6 * 2,  6 - 2, "monkey"]`

	stmt := prepExpressionTest(t, input)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 4 {
		t.Fatalf("len(array.Elements) not 4. got=%d", len(array.Elements))
	}
	testNumberLiteral[int64](t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 6, "*", 2)
	testInfixExpression(t, array.Elements[2], 6, "-", 2)

	str, ok := array.Elements[3].(*ast.StringLiteral)
	if !ok {
		t.Errorf("array.Elements[3] not ast.StringLiteral. got=%T", array.Elements[3])
	}

	if str.Value != "monkey" {
		t.Errorf("str.Value does not match. got=%q", str.Value)
	}
}

// PREFIX TEST
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTest := []struct {
		input       string
		operator    string
		numberValue interface{}
	}{
		{"!8", "!", 8},
		{"-23", "-", 23},
		{"!7.23", "!", 7.23},
		{"-84.3", "-", 84.3},
		{"!momnke", "!", "momnke"},
		{"-paw", "-", "paw"},
		{"!true", "!", true},
		{"!false", "!", false},
	}
	for _, tt := range prefixTest {
		stmt := prepExpressionTest(t, tt.input)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.numberValue) {
			return
		}
	}
}

// INFIX TEST
func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"8 + 8", 8, "+", 8},
		{"8 - 8", 8, "-", 8},
		{"8 * 8", 8, "*", 8},
		{"8 / 8", 8, "/", 8},
		{"8 % 8", 8, "%", 8},
		{"8 > 8", 8, ">", 8},
		{"8 < 8", 8, "<", 8},
		{"8 != 8", 8, "!=", 8},
		{"8 == 8", 8, "==", 8},
		{"monkey + paw;", "monkey", "+", "paw"},
		{"monkey - paw;", "monkey", "-", "paw"},
		{"monkey * paw;", "monkey", "*", "paw"},
		{"monkey / paw;", "monkey", "/", "paw"},
		{"monkey > paw;", "monkey", ">", "paw"},
		{"monkey < paw;", "monkey", "<", "paw"},
		{"monkey == paw;", "monkey", "==", "paw"},
		{"monkey != paw;", "monkey", "!=", "paw"},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range infixTests {
		stmt := prepExpressionTest(t, tt.input)
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

// PRECEDENCE TEST
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"-a + b % 5 ", "((-a) + (b % 5))"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a * b % c", "(a * (b % c))"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c % 5 / d + e - c", "(((a + ((b * (c % 5)) / d)) + e) - c)"},
		{"3 + 6; -25 * 6 % 3", "(3 + 6)((-25) * (6 % 3))"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 12 != !a > 5", "((5 < 12) != ((!a) > 5))"},
		{"3.63432 + 1.436 * 35", "(3.63432 + (1.436 * 35))"},
		{"3 + 4 * 5 == 3 * 1 + 4 % 2 % 2", "((3 + (4 * 5)) == ((3 * 1) + ((4 % 2) % 2)))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		//Boolean
		{"true", "true"},
		{"false", "false"},
		{"3 > 7 == false", "((3 > 7) == false)"},
		{"6 < 10 == true", "((6 < 10) == true)"},
		//GROUPING
		{"1 + (3 + 6) - 3", "((1 + (3 + 6)) - 3)"},
		{"(5 + 18) * 3", "((5 + 18) * 3)"},
		{"2 / (8 + 1)", "(2 / (8 + 1))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == false)", "(!(true == false))"},
		{"(5 + 5 % 3) * 3 + 12.23 / 3", "(((5 + (5 % 3)) * 3) + (12.23 / 3))"},
		//CALL EXPRESSION
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a , b, 1, 35 * 2, 4 + 5 + 88 + add(3, 5 * 9), 53)", "add(a, b, 1, (35 * 2), (((4 + 5) + 88) + add(3, (5 * 9))), 53)"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		program := prepTest(t, tt.input)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("exptected=%q, got=%q", tt.expected, actual)
		}
	}
}

// CONDITION TEST
func TestIfExpression(t *testing.T) {
	input := `if(x < y) { x }`
	stmt := prepExpressionTest(t, input)
	exp, ok := stmt.Expression.(*ast.IFExpression)

	if !ok {
		t.Fatalf("stmt.Expresion is not ast.ConditionalExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	stmt := prepExpressionTest(t, input)
	exp, ok := stmt.Expression.(*ast.IFExpression)

	if !ok {
		t.Fatalf("stmt.Expresion is not ast.ConditionalExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Statements[0] is not ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

// FUNCTION TEST
func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(a, b){ a + b; }`
	stmt := prepExpressionTest(t, input)
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral, got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters not %d parameters, got=%d", 2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "a")
	testLiteralExpression(t, function.Parameters[1], "b")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements not %d statement, got=%d", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement, got=%T", function.Body.Statements)
	}
	testInfixExpression(t, bodyStmt.Expression, "a", "+", "b")
}
func TestFunctionParameterParsing(t *testing.T) {
	test := []struct {
		input          string
		expectedParams []string
	}{
		{"fn(){};", []string{}},
		{"fn(x){};", []string{"x"}},
		{"fn(x, y, z){};", []string{"x", "y", "z"}},
	}

	for _, tt := range test {
		stmt := prepExpressionTest(t, tt.input)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// CALLEXPRESSION TEST
func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 + 3, 8 + 2)"
	stmt := prepExpressionTest(t, input)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression, got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments, expected=%d, got=%d", 3, len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "+", 3)
	testInfixExpression(t, exp.Arguments[2], 8, "+", 2)
}

func TestCallExpressionArgumentParsing(t *testing.T) {
	testCase := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{"add()", "add", []string{}},
		{"remove(2)", "remove", []string{"2"}},
		{"add(5, 4 - 5, 23, 8 * 12 + 3)", "add", []string{"5", "(4 - 5)", "23", "((8 * 12) + 3)"}},
	}

	for _, tt := range testCase {
		stmt := prepExpressionTest(t, tt.input)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("Wrong number of arguement, expected=%d, got=%d", len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Fatalf("argument %d is wrong. expected=%s. got=%s", i, arg, exp.Arguments[i].String())
			}
		}
	}
}

// ----------------------- UTILITY ----------------------
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

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %q, got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testNumberLiteral[int64](t, exp, int64(v))
	case float64:
		return testNumberLiteral[float64](t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("Undefined type of exp. type of exp not handled, got=%T", exp)
	return false
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

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolExp, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean, got=%T", exp)
		return false
	}

	if boolExp.Value != value {
		t.Errorf("boolExp.value not %t, got=%t", value, boolExp.Value)
		return false
	}

	if boolExp.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolExp.TokenLiteral not %t, got=%s", value, boolExp.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s, got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

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

func prepTest(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return program
}

func prepExpressionTest(t *testing.T, input string) *ast.ExpressionStatement {
	program := prepTest(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statments[0] id not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	return stmt
}
