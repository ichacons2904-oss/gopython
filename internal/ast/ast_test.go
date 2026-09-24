package ast

import (
	"testing"

	"gopython/internal/source"
)

var (
	_ Node       = Program{}
	_ Expression = IntegerLiteral{}
	_ Expression = StringLiteral{}
	_ Expression = BooleanLiteral{}
	_ Expression = NoneLiteral{}
	_ Expression = Identifier{}
	_ Expression = UnaryExpression{}
	_ Expression = BinaryExpression{}
	_ Expression = ComparisonExpression{}
	_ Expression = CallExpression{}
	_ Statement  = Assignment{}
	_ Statement  = ExpressionStatement{}
	_ Statement  = IfStatement{}
	_ Statement  = WhileStatement{}
	_ Statement  = ForStatement{}
	_ Statement  = ReturnStatement{}
	_ Statement  = BreakStatement{}
	_ Statement  = ContinueStatement{}
	_ Statement  = FunctionDefinition{}
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
		Expression: Identifier{Pos: position, Name: "print"},
	}

	if got := statement.Position(); got != position {
		t.Fatalf("position = %#v, want %#v", got, position)
	}
}
