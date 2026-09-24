package ast

import (
	"gopython/internal/lexer"
	"gopython/internal/source"
)

type IntegerLiteral struct {
	Pos   source.Position
	Value int64
}

func (node IntegerLiteral) Position() source.Position {
	return node.Pos
}

func (IntegerLiteral) expressionNode() {}

type StringLiteral struct {
	Pos   source.Position
	Value string
}

func (node StringLiteral) Position() source.Position {
	return node.Pos
}

func (StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Pos   source.Position
	Value bool
}

func (node BooleanLiteral) Position() source.Position {
	return node.Pos
}

func (BooleanLiteral) expressionNode() {}

type NoneLiteral struct {
	Pos source.Position
}

func (node NoneLiteral) Position() source.Position {
	return node.Pos
}

func (NoneLiteral) expressionNode() {}

type Identifier struct {
	Pos  source.Position
	Name string
}

func (node Identifier) Position() source.Position {
	return node.Pos
}

func (Identifier) expressionNode() {}

type UnaryExpression struct {
	Pos      source.Position
	Operator lexer.TokenType
	Operand  Expression
}

func (node UnaryExpression) Position() source.Position {
	return node.Pos
}

func (UnaryExpression) expressionNode() {}

type BinaryExpression struct {
	Pos      source.Position
	Left     Expression
	Operator lexer.TokenType
	Right    Expression
}

func (node BinaryExpression) Position() source.Position {
	return node.Pos
}

func (BinaryExpression) expressionNode() {}

type ComparisonExpression struct {
	Pos         source.Position
	Left        Expression
	Comparisons []Comparison
}

type Comparison struct {
	Pos      source.Position
	Operator lexer.TokenType
	Right    Expression
}

func (node ComparisonExpression) Position() source.Position {
	return node.Pos
}

func (ComparisonExpression) expressionNode() {}

type CallExpression struct {
	Pos       source.Position
	Callee    Expression
	Arguments []Expression
}

func (node CallExpression) Position() source.Position {
	return node.Pos
}

func (CallExpression) expressionNode() {}
