package ast

import "gopython/internal/lexer"

type Node interface {
	Position() lexer.Token
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Token      lexer.Token
	Statements []Statement
}

func (node Program) Position() lexer.Token {
	return node.Token
}
