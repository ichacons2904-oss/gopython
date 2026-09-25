package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseComparisonChain(t *testing.T) {
	tokens, err := lexer.New("a < b < c\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	expressionStatement := program.Statements[0].(*ast.ExpressionStatement)
	comparison, ok := expressionStatement.Expression.(*ast.ComparisonExpression)
	if !ok {
		t.Fatalf("expression = %T, want ast.ComparisonExpression", expressionStatement.Expression)
	}
	if len(comparison.Comparisons) != 2 {
		t.Fatalf("comparison count = %d, want 2", len(comparison.Comparisons))
	}
	if comparison.Comparisons[0].Operator != lexer.LessThan || comparison.Comparisons[1].Operator != lexer.LessThan {
		t.Fatalf("comparison operators = %#v, want two < operators", comparison.Comparisons)
	}
	if comparison.Comparisons[0].Right.(*ast.Identifier).Name != "b" {
		t.Fatalf("first comparison right side = %#v, want b", comparison.Comparisons[0].Right)
	}
	if comparison.Comparisons[1].Right.(*ast.Identifier).Name != "c" {
		t.Fatalf("second comparison right side = %#v, want c", comparison.Comparisons[1].Right)
	}
}

func TestNotBindsLooserThanComparison(t *testing.T) {
	tokens, err := lexer.New("not a == b\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	expressionStatement := program.Statements[0].(*ast.ExpressionStatement)
	notExpression, ok := expressionStatement.Expression.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("expression = %T, want ast.UnaryExpression", expressionStatement.Expression)
	}
	if notExpression.Operator != lexer.Not {
		t.Fatalf("unary operator = %s, want %s", notExpression.Operator, lexer.Not)
	}
	if _, ok := notExpression.Operand.(*ast.ComparisonExpression); !ok {
		t.Fatalf("not operand = %T, want ast.ComparisonExpression", notExpression.Operand)
	}
}

func TestParseRejectsIntegerOverflow(t *testing.T) {
	tokens, err := lexer.New("x = 9223372036854775808\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	_, err = New(tokens).Parse()
	if err == nil {
		t.Fatal("Parse() returned nil error")
	}
	parserError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want parser.Error", err)
	}
	if parserError.Message != "integer literal is too large" {
		t.Fatalf("error message = %q", parserError.Message)
	}
}
