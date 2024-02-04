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
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	lex := lexer.New(input)
	p := parser.New(lex)

	program := p.ParseProgram()
	s.Require().NotNil(program, "program was nil")
	s.Require().Len(p.Errors(), 0, "parser had errors")

	s.Require().Len(program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(s, stmt, tt.expectedIdentifier)
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
	input := `
return 5;	
return 10;
return 993322;
`
	lex := lexer.New(input)
	p := parser.New(lex)

	program := p.ParseProgram()
	s.Require().NotNil(program, "program was nil")
	s.Require().Len(p.Errors(), 0, "parser had errors")

	s.Require().Len(program.Statements, 3)

	for _, stmt := range program.Statements {
		letStmt, ok := stmt.(*ast.ReturnStatement)
		s.Require().Truef(ok, "s not *ast.ReturnStatement. got=%T", stmt)

		s.Require().Equal("return", letStmt.TokenLiteral())
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
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
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

		testIntegerLiteral(s, exp.Right, tt.integerValue)
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
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
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

		testIntegerLiteral(s, exp.Left, tt.leftValue)

		s.Require().Equal(tt.operator, exp.Operator)

		testIntegerLiteral(s, exp.Right, tt.rightValue)
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
