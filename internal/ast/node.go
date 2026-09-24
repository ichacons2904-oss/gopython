package ast

import "gopython/internal/source"

type Node interface {
	Position() source.Position
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
	Pos        source.Position
	Statements []Statement
}

func (node Program) Position() source.Position {
	return node.Pos
}
