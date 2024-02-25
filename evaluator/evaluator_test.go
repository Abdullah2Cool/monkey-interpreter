package evaluator_test

import (
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func (s *Suite) SetupTest() {
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
func (s *Suite) TestEvalIntegerExpression() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(s, evaluated, tt.expected)
	}
}

func (s *Suite) TestEvalBooleanExpression() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(s, evaluated, tt.expected)
	}
}

func (s *Suite) TestBangOperator() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(s, evaluated, tt.expected)
	}
}

func (s *Suite) TestIfElseExpressions() {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(s, evaluated, int64(integer))
		} else {
			testNullObject(s, evaluated)
		}
	}
}

func (s *Suite) TestReturnStatements() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},

		{
			`
	if (10 > 1) {
		if (10 > 1) {
			return 10;
		}
		return 1;
	}
	`,
			10,
		},
		{
			`
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
let f = fn(x) {
   let result = x + 10;
   return result;
   return 10;
};
f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(s, evaluated, tt.expected)
	}
}

func (s *Suite) TestErrorHandling() {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
	if (10 > 1) {
		if (10 > 1) {
			return true + false;
		}
		return 1;
	}
	`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)

		s.Require().Truef(ok, "no error object returned. got=%T(%+v)", evaluated, evaluated)

		s.Require().Equal(tt.expectedMessage, errObj.Message)
	}
}

func (s *Suite) TestLetStatements() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(s, testEval(tt.input), tt.expected)
	}
}

func (s *Suite) TestFunctionObject() {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	s.Require().Truef(ok, "object is not Function. got=%T (%+v)", evaluated, evaluated)

	s.Require().Len(fn.Parameters, 1)

	s.Require().Equal("x", fn.Parameters[0].String())

	s.Require().Equal("(x + 2)", fn.Body.String())
}

func (s *Suite) TestFunctionApplication() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(s, testEval(tt.input), tt.expected)
	}
}

func (s *Suite) TestClosures() {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};
	let addTwo = newAdder(2);
	addTwo(2);`

	testIntegerObject(s, testEval(input), 4)
}

func (s *Suite) TestStringLiteral() {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	s.Require().Truef(ok, "expected *object.String but got %T", evaluated)

	s.Require().Equal("Hello World!", str.Value)
}

func (s *Suite) TestStringConcatenation() {
	input := `"Hello" + " " +  "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	s.Require().Truef(ok, "expected *object.String but got %T (%+v)", evaluated, evaluated)

	s.Require().Equal("Hello World!", str.Value)
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testIntegerObject(s *Suite, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	s.Require().True(ok, "expected *object.Integer but got %T", obj)

	s.Require().Equal(expected, result.Value)
}

func testBooleanObject(s *Suite, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	s.Require().Truef(ok, "expected *object.Boolean but got %T", obj)

	s.Require().Equal(expected, result.Value)
}

func testNullObject(s *Suite, obj object.Object) {
	s.Require().Equal(evaluator.NULL, obj)
}
