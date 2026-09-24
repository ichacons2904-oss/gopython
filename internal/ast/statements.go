package ast

import "gopython/internal/lexer"

type Assignment struct {
	Token lexer.Token
	Name  Identifier
	Value Expression
}

func (node Assignment) Position() lexer.Token {
	return node.Token
}

func (Assignment) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (node ExpressionStatement) Position() lexer.Token {
	return node.Expression.Position()
}

func (ExpressionStatement) statementNode() {}

type FunctionDefinition struct {
	Token      lexer.Token
	Name       Identifier
	Parameters []Identifier
	Body       []Statement
}

func (node FunctionDefinition) Position() lexer.Token {
	return node.Token
}

func (FunctionDefinition) statementNode() {}
