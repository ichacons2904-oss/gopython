package ast

import (
	"testing"

	"gopython/internal/source"
)

var (
	_ Node       = (*Program)(nil)
	_ Expression = (*IntegerLiteral)(nil)
	_ Expression = (*StringLiteral)(nil)
	_ Expression = (*BooleanLiteral)(nil)
	_ Expression = (*NoneLiteral)(nil)
	_ Expression = (*Identifier)(nil)
	_ Expression = (*UnaryExpression)(nil)
	_ Expression = (*BinaryExpression)(nil)
	_ Expression = (*ComparisonExpression)(nil)
	_ Expression = (*CallExpression)(nil)
	_ Statement  = (*Assignment)(nil)
	_ Statement  = (*ExpressionStatement)(nil)
	_ Statement  = (*IfStatement)(nil)
	_ Statement  = (*WhileStatement)(nil)
	_ Statement  = (*ForStatement)(nil)
	_ Statement  = (*ReturnStatement)(nil)
	_ Statement  = (*BreakStatement)(nil)
	_ Statement  = (*ContinueStatement)(nil)
	_ Statement  = (*FunctionDefinition)(nil)
)

func TestASTNodesPreserveSourcePosition(t *testing.T) {
	position := source.Position{Line: 3, Column: 7}
	node := IntegerLiteral{Pos: position, Value: 42}

	if got := node.Position(); got != position {
		t.Fatalf("position = %#v, want %#v", got, position)
	}
}

func TestExpressionStatementUsesExpressionPosition(t *testing.T) {
	position := source.Position{Line: 2, Column: 1}
	statement := ExpressionStatement{
		Expression: &Identifier{Pos: position, Name: "print"},
	}

	if got := statement.Position(); got != position {
		t.Fatalf("position = %#v, want %#v", got, position)
	}
}
