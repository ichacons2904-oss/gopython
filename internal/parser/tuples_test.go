package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
	"testing"
)

func TestTupleAndCallCommaParsing(t *testing.T) {
	tokens, err := lexer.New("a = (1)\nb = (1,)\nc = 1, 2\nf(1, 2, (3, 4))\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := program.Statements[0].(*ast.Assignment).Value.(*ast.IntegerLiteral); !ok {
		t.Fatal("grouped value must not be tuple")
	}
	for _, index := range []int{1, 2} {
		if _, ok := program.Statements[index].(*ast.Assignment).Value.(*ast.TupleLiteral); !ok {
			t.Fatalf("assignment %d is not a tuple", index)
		}
	}
	call := program.Statements[3].(*ast.ExpressionStatement).Expression.(*ast.CallExpression)
	if len(call.Arguments) != 3 {
		t.Fatalf("got %d arguments, want 3", len(call.Arguments))
	}
	if _, ok := call.Arguments[2].(*ast.TupleLiteral); !ok {
		t.Fatal("third argument must be tuple")
	}
}
