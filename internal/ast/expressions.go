package ast

import "gopython/internal/lexer"

type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (node IntegerLiteral) Position() lexer.Token {
	return node.Token
}

func (IntegerLiteral) expressionNode() {}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (node StringLiteral) Position() lexer.Token {
	return node.Token
}

func (StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (node BooleanLiteral) Position() lexer.Token {
	return node.Token
}

func (BooleanLiteral) expressionNode() {}

type NoneLiteral struct {
	Token lexer.Token
}

func (node NoneLiteral) Position() lexer.Token {
	return node.Token
}

func (NoneLiteral) expressionNode() {}

type Identifier struct {
	Token lexer.Token
	Name  string
}

func (node Identifier) Position() lexer.Token {
	return node.Token
}

func (Identifier) expressionNode() {}

type UnaryExpression struct {
	Token    lexer.Token
	Operator lexer.TokenType
	Operand  Expression
}

func (node UnaryExpression) Position() lexer.Token {
	return node.Token
}

func (UnaryExpression) expressionNode() {}

type BinaryExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator lexer.TokenType
	Right    Expression
}

func (node BinaryExpression) Position() lexer.Token {
	return node.Token
}

func (BinaryExpression) expressionNode() {}

type CallExpression struct {
	Token     lexer.Token
	Callee    Expression
	Arguments []Expression
}

func (node CallExpression) Position() lexer.Token {
	return node.Token
}

func (CallExpression) expressionNode() {}
