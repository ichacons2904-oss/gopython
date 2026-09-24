package ast

import (
	"testing"

	"gopython/internal/lexer"
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
	token := lexer.Token{
		Type:   lexer.Integer,
		Lexeme: "42",
		Line:   3,
		Column: 7,
	}
	node := IntegerLiteral{Token: token, Value: 42}

	if got := node.Position(); got != token {
		t.Fatalf("position = %#v, want %#v", got, token)
	}
}

func TestExpressionStatementUsesExpressionPosition(t *testing.T) {
	token := lexer.Token{
		Type:   lexer.Identifier,
		Lexeme: "print",
		Line:   2,
		Column: 1,
	}
	statement := ExpressionStatement{
		Expression: Identifier{Token: token, Name: "print"},
	}

	if got := statement.Position(); got != token {
		t.Fatalf("position = %#v, want %#v", got, token)
	}
}
