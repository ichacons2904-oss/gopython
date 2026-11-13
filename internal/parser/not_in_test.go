package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseNotInAsChainedComparison(t *testing.T) {
	tokens, err := lexer.New("value not in items == other\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[1].Type != lexer.Not || tokens[2].Type != lexer.In {
		t.Fatalf("not in tokens = %s, %s; want NOT, IN", tokens[1].Type, tokens[2].Type)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	statement := program.Statements[0].(*ast.ExpressionStatement)
	comparison, ok := statement.Expression.(*ast.ComparisonExpression)
	if !ok || len(comparison.Comparisons) != 2 {
		t.Fatalf("expression = %#v, want a comparison chain", statement.Expression)
	}
	if comparison.Comparisons[0].Operator != lexer.NotIn || comparison.Comparisons[1].Operator != lexer.Equal {
		t.Fatalf("operators = %#v, want NOT_IN then EQUAL", comparison.Comparisons)
	}
}

func TestParseRejectsIncompleteNotIn(t *testing.T) {
	for _, input := range []string{"1 not [1]\n", "1 not in\n"} {
		t.Run(input, func(t *testing.T) {
			tokens, err := lexer.New(input).Lex()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := New(tokens).Parse(); err == nil {
				t.Fatal("Parse() accepted an incomplete comparison")
			}
		})
	}
}
