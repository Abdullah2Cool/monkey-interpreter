package parser_test

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"

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

func (s *Suite) TestLetStatements() {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := parser.New(lex)

		program := p.ParseProgram()
		s.Require().NotNil(program, "program was nil")
		s.Require().Len(p.Errors(), 0, "parser had errors")

		s.Require().Len(program.Statements, 1)

		stmt := program.Statements[0]
		testLetStatement(s, stmt, tt.expectedIdentifier)

		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(s, val, tt.expectedValue)
	}
}

func testLetStatement(s *Suite, stmt ast.Statement, name string) {
	s.Require().Equal("let", stmt.TokenLiteral())

	letStmt, ok := stmt.(*ast.LetStatement)
	s.Require().Truef(ok, "s not *ast.LetStatement. got=%T", stmt)

	s.Require().Equal(letStmt.Name.Value, name)
	s.Require().Equal(letStmt.Name.TokenLiteral(), name)
}

func (s *Suite) TestReturnStatements() {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := parser.New(lex)

		program := p.ParseProgram()
		s.Require().NotNil(program, "program was nil")
		s.Require().Len(p.Errors(), 0, "parser had errors")

		s.Require().Len(program.Statements, 1)

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		s.Require().Truef(ok, "s not *ast.ReturnStatement. got=%T", stmt)

		s.Require().Equal("return", returnStmt.TokenLiteral())

		testLiteralExpression(s, returnStmt.ReturnValue, tt.expectedValue)
	}
}

func (s *Suite) TestIdentifierExpression() {
	input := "foobar;"

	lex := lexer.New(input)
	p := parser.New(lex)

	program := p.ParseProgram()
	s.Require().NotNil(program, "program was nil")
	s.Require().Len(p.Errors(), 0, "parser had errors")

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", stmt)

	ident, ok := stmt.Expression.(*ast.Identifier)
	s.Require().Truef(ok, "s not *ast.Identifier. got=%T", stmt)

	s.Require().Equal("foobar", ident.TokenLiteral())
}

func (s *Suite) TestIntegerLiteralExpression() {
	input := "5;"

	lex := lexer.New(input)
	p := parser.New(lex)

	program := p.ParseProgram()
	s.Require().NotNil(program, "program was nil")
	s.Require().Len(p.Errors(), 0, "parser had errors")

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", stmt)

	intLiteral, ok := stmt.Expression.(*ast.IntegerLiteral)
	s.Require().Truef(ok, "s not *ast.IntegerLiteral. got=%T", stmt)

	s.Require().Equal("5", intLiteral.TokenLiteral())
}

func (s *Suite) TestParsingPrefixExpressions() {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		p := parser.New(lex)

		program := p.ParseProgram()
		s.Require().NotNil(program, "program was nil")
		s.Require().Len(p.Errors(), 0, "parser had errors")

		s.Require().Len(program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", stmt)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		s.Require().Truef(ok, "s not *ast.PrefixExpression. got=%T", stmt)

		testLiteralExpression(s, exp.Right, tt.integerValue)
	}
}

func testIntegerLiteral(s *Suite, il ast.Expression, value int64) {
	intLiteral, ok := il.(*ast.IntegerLiteral)
	s.Require().Truef(ok, "s not *ast.IntegerLiteral. got=%T", intLiteral)

	s.Require().Equal(value, intLiteral.Value)

	s.Require().Equal(fmt.Sprintf("%d", value), intLiteral.TokenLiteral())
}

func (s *Suite) TestParsingInfixExpressions() {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		p := parser.New(lex)

		program := p.ParseProgram()
		s.Require().NotNil(program, "program was nil")
		s.Require().Len(p.Errors(), 0, "parser had errors")

		s.Require().Len(program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", stmt)

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		s.Require().Truef(ok, "s not *ast.InfixExpression. got=%T", stmt)

		testInfixExpression(s, exp, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func (s *Suite) TestOperatorPrecedenceParsing() {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := parser.New(lex)

		program := p.ParseProgram()
		s.Require().NotNil(program, "program was nil")
		s.Require().Len(p.Errors(), 0, "parser had errors")

		actual := program.String()
		s.Require().Equal(tt.expected, actual)
	}
}

func testIdentifier(s *Suite, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	s.Require().Truef(ok, "exp not *ast.Identifier. got=%T", exp)

	s.Require().Equal(value, ident.Value)

	s.Require().Equal(ident.TokenLiteral(), value)
}

func testLiteralExpression(s *Suite, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(s, exp, int64(v))
	case int64:
		testIntegerLiteral(s, exp, v)
	case string:
		testIdentifier(s, exp, v)
	case bool:
		testBooleanLiteral(s, exp, v)
	default:
		s.Require().Fail(fmt.Sprintf("type of exp not handled. got=%T", exp))
	}
}

func testInfixExpression(s *Suite, exp ast.Expression, left interface{}, operator string, right interface{}) {
	opExp, ok := exp.(*ast.InfixExpression)
	s.Require().Truef(ok, "exp not *ast.InfixExpression. got=%T", exp)

	testLiteralExpression(s, opExp.Left, left)

	s.Require().Equal(operator, opExp.Operator)

	testLiteralExpression(s, opExp.Right, right)
}

func testBooleanLiteral(s *Suite, exp ast.Expression, value bool) {
	boolLiteral, ok := exp.(*ast.BooleanLiteral)
	s.Require().Truef(ok, "exp not *ast.BooleanLiteral. got=%T", exp)

	s.Require().Equal(value, boolLiteral.Value)

	s.Require().Equal(fmt.Sprintf("%t", value), boolLiteral.TokenLiteral())
}

func (s *Suite) TestIfExpression() {
	input := `if (x < y) { x }`

	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	s.Require().Len(p.Errors(), 0)

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	s.Require().Truef(ok, "s not *ast.IfExpression. got=%T", stmt.Expression)

	testInfixExpression(s, exp.Condition, "x", "<", "y")

	s.Require().NotNil(exp.Consequence)
	s.Require().Len(exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	testIdentifier(s, consequence.Expression, "x")

	s.Require().Nil(exp.Alternative)
}

func (s *Suite) TestIfElseExpression() {
	input := `if (x < y) { x } else { y }`

	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	s.Require().Len(p.Errors(), 0)

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	s.Require().Truef(ok, "s not *ast.IfExpression. got=%T", stmt.Expression)

	testInfixExpression(s, exp.Condition, "x", "<", "y")

	s.Require().NotNil(exp.Consequence)
	s.Require().Len(exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	testIdentifier(s, consequence.Expression, "x")

	s.Require().NotNil(exp.Alternative)
	s.Require().Len(exp.Alternative.Statements, 1)

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])

	testIdentifier(s, alternative.Expression, "y")
}

func (s *Suite) TestFunctionLiteralParsing() {
	input := `fn(x, y) { x + y; }`

	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	s.Require().Len(p.Errors(), 0)

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	s.Require().Truef(ok, "s not *ast.FunctionLiteral. got=%T", stmt.Expression)

	s.Require().Len(function.Parameters, 2)

	testLiteralExpression(s, function.Parameters[0], "x")
	testLiteralExpression(s, function.Parameters[1], "y")

	s.Require().Len(function.Body.Statements, 1)

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", function.Body.Statements[0])

	testInfixExpression(s, bodyStmt.Expression, "x", "+", "y")
}

func (s *Suite) TestFunctionParameterParsing() {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := parser.New(lex)
		program := p.ParseProgram()

		s.Require().Len(p.Errors(), 0)
		s.Require().Len(program.Statements, 1)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		s.Require().Len(function.Parameters, len(tt.expectedParams))

		for i, ident := range tt.expectedParams {
			testLiteralExpression(s, function.Parameters[i], ident)
		}
	}
}

func (s *Suite) TestCallExpressionParsing() {
	input := `add(1, 2 * 3, 4 + 5);`

	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	s.Require().Len(p.Errors(), 0)

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.CallExpression)
	s.Require().Truef(ok, "s not *ast.CallExpression. got=%T", stmt.Expression)

	testIdentifier(s, exp.Function, "add")

	s.Require().Len(exp.Arguments, 3)

	testLiteralExpression(s, exp.Arguments[0], 1)
	testInfixExpression(s, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(s, exp.Arguments[2], 4, "+", 5)
}

func (s *Suite) TestStringLiteralExpression() {
	input := `"hello world";`

	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	s.Require().Len(p.Errors(), 0)

	s.Require().Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	s.Require().Truef(ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

	sl, ok := stmt.Expression.(*ast.StringLiteral)
	s.Require().Truef(ok, "s not *ast.StringLiteral. got=%T", stmt.Expression)

	s.Require().Equal("hello world", sl.Value)
}
