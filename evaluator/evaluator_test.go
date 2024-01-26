package evaluator

import (
	"testing"

	"github.com/arthurlee945/monkey.on/lexer"
	"github.com/arthurlee945/monkey.on/object"
	"github.com/arthurlee945/monkey.on/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"84938", 84938},
		{"-25", -25},
		{"-80", -80},
		{"5 + 5 + 5 - 10", 5},
		{"2 * 2 / 2 + 5", 7},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"50 / 2 * 2 + 10", 60},
		{"9 % 2 * 2 + 10", 12},
		{"3 * 3 * 3 + 10", 37},
		{"2 * (5 + 10)", 30},
		{"(15 * 3) % 4 + 10", 11},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.32", 5.32},
		{"8556.23", 8556.23},
		{"-5.32", -5.32},
		{"25 / 2.5 * 3 + 10", 40},
		{"9 % 2 / 2.5 + 10", 10.4},
		{"2.5 + 2.5 + 2.5 + 2.5", 10},
		{"12.24 * (2.6 + 2.4)", 61.2},
		{"(5 * 3) % 3.5 + 10", 11},
		{"(5 + 10.4 * 2 + 15.0 / 4.0) * 2 + -10", 49.1},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		testFloatObject(t, obj, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	input := `"Hello Monkey World"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello Monkey World" {
		t.Fatalf("String has wrong value. got=%q", str.Value)
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	input := `"monkey" + " " + "Says" + " " + "Hi!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "monkey Says Hi!" {
		t.Fatalf("String has wrong value. got=%q", str.Value)
	}
}

func TestEvalArrayLiteral(t *testing.T) {
	input := "[1, 2, 3 * 2, 4 + 5]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 4 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 2)
	testIntegerObject(t, result.Elements[2], 6)
	testIntegerObject(t, result.Elements[3], 9)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 5, 8][0];", 1},
		{"[1, 5, 8][1];", 5},
		{"[1, 5, 8][2];", 8},
		{"let x = 0; [2, 5][x];", 2},
		{"[1, 5, 8][1+1];", 8},
		{"let myArr = [ 6, 7, 8 ]; myArr[0];", 6},
		{"let myArr = [ 6, 7, 8 ]; myArr[0] + myArr[1] + myArr[2];", 21},
		{"let myArr = [1, 3, 9]; let x = myArr[0]; myArr[x];", 3},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"8 < 8", false},
		{"8 > 8", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{`"monkey" == "ape"`, false},
		{`"monkey" == "monkey"`, true},
		{`"monkey" != "monkey"`, false},
		{`"monkey" != "ape"`, true},
		{"1 != 2", true},
		{"12.3 < 2.55", false},
		{"12.3 > 2.55", true},
		{"1.345 == 1.345", true},
		{"1.345 != 1.345", false},
		{"true == true", true},
		{"true != true", false},
		{"false == false", true},
		{"false != false", false},
		{"(1 < 2) != false", true},
		{"(1 < 2) == false", false},
		{"false == (1 > 2)", true},
		{"false != (1 > 2)", false},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestEvalBangExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!8.5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!8", true},
		{"!!2.3", true},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if(true) { 10 }", 10},
		{`if("onepiece") { 10 }`, 10},
		{"if(false){ 25 }", nil},
		{"if(2) { 12 }", 12},
		{"if(1 < 2) { 90 }", 90},
		{"if(1 > 2) { 234 }", nil},
		{"if(false){ 34 }else{88}", 88},
		{"if(8==8){ 20 }else{88}", 20},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 25; 124;", 25},
		{"return 2 * 5 / 2; 9;", 5},
		{"9 + 9; return 3 * 5; -19;", 15},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}	
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 5; x;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 8; let b = a; b;", 8},
		{"let a = 8; let b = a; a%b;", 0},
		{"let a = 5; let b = a; let c = a + b + 25; c;", 35},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x){x+2;};"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T {%+v}", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let ident = fn(x){x;}; ident(5);", 5},
		{"let ident = fn(y){return y;}; ident(8);", 8},
		{"let double = fn(x){ x * 2; }; double(8);", 16},
		{"let add = fn(x, y){x + y; }; add(1, 7)", 8},
		{"let add = fn(x, y){ return x + y;}; add(1, add(5, 10))", 16},
		{"let w = 5; let add = fn(x, y){ return x + y + w;}; add(1, add(5, 10))", 26},
		{"fn(x){x+2;}(5)", 7},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newExtender = fn(x){
		fn(y){x + y}
	}

	let addEight = newExtender(8)
	addEight(10)
	`

	testIntegerObject(t, testEval(input), 18)
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("hello world MONKEY")`, 18},
		{`len("eight")`, 5},
		{`len(1)`, "argument to 'len' not supported, got INTEGER"},
		{`len(1.235)`, "argument to 'len' not supported, got FLOAT"},
		{`len("one", "two")`, "wrong number of arguments. got=2, expected=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"true + 5; 5;", "type mismatch: BOOLEAN + INTEGER"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if( 10 > 1) {true + false}", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
		if (10 > 1){
			if(10 > 1){
				return true + false
			}
			return 1;
		}
		`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"monkey", "identifier not found: monkey"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%d, got=%d", expected, result.Value)
		return false
	}
	return true
}
func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%f, got=%f", expected, result.Value)
		return false
	}
	return true
}
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected=%t, got=%t", expected, result.Value)
		return false
	}
	return true
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not Null. got=%T {%+v}", obj, obj)
		return false
	}
	return true
}
