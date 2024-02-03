package parser_test

import (
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
