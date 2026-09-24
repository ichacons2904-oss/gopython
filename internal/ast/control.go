package ast

import "gopython/internal/source"

type IfStatement struct {
	Pos          source.Position
	Condition    Expression
	Body         []Statement
	ElifBranches []ElifBranch
	ElseBody     []Statement
}

func (node IfStatement) Position() source.Position {
	return node.Pos
}

func (IfStatement) statementNode() {}

type ElifBranch struct {
	Pos       source.Position
	Condition Expression
	Body      []Statement
}

func (node ElifBranch) Position() source.Position {
	return node.Pos
}

type WhileStatement struct {
	Pos       source.Position
	Condition Expression
	Body      []Statement
}

func (node WhileStatement) Position() source.Position {
	return node.Pos
}

func (WhileStatement) statementNode() {}

type ForStatement struct {
	Pos      source.Position
	Target   Identifier
	Iterable Expression
	Body     []Statement
}

func (node ForStatement) Position() source.Position {
	return node.Pos
}

func (ForStatement) statementNode() {}

type ReturnStatement struct {
	Pos   source.Position
	Value Expression
}

func (node ReturnStatement) Position() source.Position {
	return node.Pos
}

func (ReturnStatement) statementNode() {}

type BreakStatement struct {
	Pos source.Position
}

func (node BreakStatement) Position() source.Position {
	return node.Pos
}

func (BreakStatement) statementNode() {}

type ContinueStatement struct {
	Pos source.Position
}

func (node ContinueStatement) Position() source.Position {
	return node.Pos
}

func (ContinueStatement) statementNode() {}
