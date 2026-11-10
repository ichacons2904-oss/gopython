package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseListsAndIndexedAssignment(t *testing.T) {
	tokens, err := lexer.New("items = [1, [2, 3],]\nitems[0] = 4\nitems[1][0]\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 3 {
		t.Fatalf("statement count = %d, want 3", len(program.Statements))
	}
	assignment, ok := program.Statements[0].(*ast.Assignment)
	if !ok {
		t.Fatalf("first statement = %T, want *ast.Assignment", program.Statements[0])
	}
	if _, ok := assignment.Value.(*ast.ListLiteral); !ok {
		t.Fatalf("assignment value = %T, want *ast.ListLiteral", assignment.Value)
	}
	if _, ok := program.Statements[1].(*ast.IndexAssignment); !ok {
		t.Fatalf("second statement = %T, want *ast.IndexAssignment", program.Statements[1])
	}
	last, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("third statement = %T, want *ast.ExpressionStatement", program.Statements[2])
	}
	if _, ok := last.Expression.(*ast.IndexExpression); !ok {
		t.Fatalf("last expression = %T, want *ast.IndexExpression", last.Expression)
	}
}

func TestParseInAsComparisonOperator(t *testing.T) {
	tokens, err := lexer.New("value + 1 in [2]\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	statement := program.Statements[0].(*ast.ExpressionStatement)
	comparison, ok := statement.Expression.(*ast.ComparisonExpression)
	if !ok {
		t.Fatalf("expression = %T, want *ast.ComparisonExpression", statement.Expression)
	}
	if len(comparison.Comparisons) != 1 || comparison.Comparisons[0].Operator != lexer.In {
		t.Fatalf("comparison operators = %#v, want one IN operator", comparison.Comparisons)
	}
}

func TestParseAttributeAccessAndCall(t *testing.T) {
	tokens, err := lexer.New("items.copy().append(1)\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	statement := program.Statements[0].(*ast.ExpressionStatement)
	outerCall, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression = %T, want *ast.CallExpression", statement.Expression)
	}
	attribute, ok := outerCall.Callee.(*ast.AttributeExpression)
	if !ok || attribute.Name != "append" {
		t.Fatalf("callee = %#v, want append attribute", outerCall.Callee)
	}
	innerCall, ok := attribute.Object.(*ast.CallExpression)
	if !ok {
		t.Fatalf("attribute receiver = %T, want *ast.CallExpression", attribute.Object)
	}
	innerAttribute, ok := innerCall.Callee.(*ast.AttributeExpression)
	if !ok || innerAttribute.Name != "copy" {
		t.Fatalf("inner callee = %#v, want copy attribute", innerCall.Callee)
	}
}
