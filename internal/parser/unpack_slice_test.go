package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
	"testing"
)

func TestParseUnpackingAndSlice(t *testing.T) {
	tokens, err := lexer.New("a, = [1]\n(a, b[0], c.x) = (1, 2, 3)\nitems[1:5:2]\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if a, ok := program.Statements[0].(*ast.UnpackAssignment); !ok || len(a.Targets) != 1 {
		t.Fatal("expected singleton unpacking")
	}
	if a, ok := program.Statements[1].(*ast.UnpackAssignment); !ok || len(a.Targets) != 3 {
		t.Fatal("expected three targets")
	}
	s, ok := program.Statements[2].(*ast.ExpressionStatement).Expression.(*ast.SliceExpression)
	if !ok || s.Start == nil || s.Stop == nil || s.Step == nil {
		t.Fatal("expected slice with three bounds")
	}
}

func TestRejectUnsupportedAssignmentTargets(t *testing.T) {
	for _, input := range []string{"a[1:2] = [3]\n", "a, (b, c) = [1, (2, 3)]\n", "a, [b, c] = [1, [2, 3]]\n", "a, *b = [1, 2]\n", "a, 1 = [1, 2]\n"} {
		t.Run(input, func(t *testing.T) {
			tokens, err := lexer.New(input).Lex()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := New(tokens).Parse(); err == nil {
				t.Fatal("expected parser error")
			}
		})
	}
}
