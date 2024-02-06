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

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	return evaluator.Eval(program)
}

func testIntegerObject(s *Suite, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	s.Require().True(ok)

	s.Require().Equal(expected, result.Value)
}

func testBooleanObject(s *Suite, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	s.Require().True(ok)

	s.Require().Equal(expected, result.Value)
}
