package ast

import "monkey/token"

// every node must implement this interface
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode() // not strictly necessary, added to distinguish from other types of nodes
}

type Expression interface {
	Node
	expressionNode() // not strictly necessary, added to distinguish from other types of nodes
}

// Program: every program is a list of statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
