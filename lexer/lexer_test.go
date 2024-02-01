package lexer_test

import (
	"testing"

	"monkey/lexer"
	"monkey/token"

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

func (s *Suite) TestNextToken() {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := lexer.New(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		s.Require().Equal(tt.expectedType, tok.Type, i)
		s.Require().Equal(tt.expectedLiteral, tok.Literal, i)
	}
}
