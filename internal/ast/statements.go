package ast

import "gopython/internal/source"

type Assignment struct {
	Pos   source.Position
	Name  Identifier
	Value Expression
}

func (node Assignment) Position() source.Position {
	return node.Pos
}

func (*Assignment) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (node ExpressionStatement) Position() source.Position {
	return node.Expression.Position()
}

func (*ExpressionStatement) statementNode() {}

type FunctionDefinition struct {
	Pos        source.Position
	Name       Identifier
	Parameters []Identifier
	Body       []Statement
}

func (node FunctionDefinition) Position() source.Position {
	return node.Pos
}

func (*FunctionDefinition) statementNode() {}

type NonlocalStatement struct {
	Pos   source.Position
	Names []Identifier
}

type GlobalStatement struct {
	Pos   source.Position
	Names []Identifier
}

func (node GlobalStatement) Position() source.Position {
	return node.Pos
}

func (*GlobalStatement) statementNode() {}

func (node NonlocalStatement) Position() source.Position {
	return node.Pos
}

func (*NonlocalStatement) statementNode() {}
