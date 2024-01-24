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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
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
