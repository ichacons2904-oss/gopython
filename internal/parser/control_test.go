package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseIfElifElse(t *testing.T) {
	source := "if x > 0:\n    print(x)\nelif x == 0:\n    print(0)\nelse:\n    print(-x)\n"
	tokens, err := lexer.New(source).Lex()
	if err != nil {
		t.Fatal(err)
	}

	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("statement count = %d, want 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("statement = %T, want ast.IfStatement", program.Statements[0])
	}
	if len(statement.Body) != 1 {
		t.Fatalf("if body length = %d, want 1", len(statement.Body))
	}
	if len(statement.ElifBranches) != 1 {
		t.Fatalf("elif count = %d, want 1", len(statement.ElifBranches))
	}
	if len(statement.ElifBranches[0].Body) != 1 {
		t.Fatalf("elif body length = %d, want 1", len(statement.ElifBranches[0].Body))
	}
	if len(statement.ElseBody) != 1 {
		t.Fatalf("else body length = %d, want 1", len(statement.ElseBody))
	}
}

func TestParseNestedIfBlock(t *testing.T) {
	source := "if outer:\n    if inner:\n        print(inner)\n    print(outer)\n"
	tokens, err := lexer.New(source).Lex()
	if err != nil {
		t.Fatal(err)
	}

	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}

	outer := program.Statements[0].(*ast.IfStatement)
	inner, ok := outer.Body[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("nested statement = %T, want ast.IfStatement", outer.Body[0])
	}
	if len(inner.Body) != 1 || len(outer.Body) != 2 {
		t.Fatalf("nested body lengths = %d and %d, want 1 and 2", len(inner.Body), len(outer.Body))
	}
}

func TestParseRejectsEmptyBlock(t *testing.T) {
	tokens, err := lexer.New("if x:\n").Lex()
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
	if parserError.Message != "expected INDENT, got EOF" {
		t.Fatalf("error message = %q", parserError.Message)
	}
}
