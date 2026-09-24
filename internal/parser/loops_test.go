package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseWhileAndFor(t *testing.T) {
	source := "while x < 3:\n    if x == 1:\n        continue\n    x = x + 1\n    if x > 2:\n        break\nfor item in items:\n    print(item)\n"
	tokens, err := lexer.New(source).Lex()
	if err != nil {
		t.Fatal(err)
	}

	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("statement count = %d, want 2", len(program.Statements))
	}

	whileStatement, ok := program.Statements[0].(ast.WhileStatement)
	if !ok {
		t.Fatalf("first statement = %T, want ast.WhileStatement", program.Statements[0])
	}
	if len(whileStatement.Body) != 3 {
		t.Fatalf("while body length = %d, want 3", len(whileStatement.Body))
	}

	firstIf, ok := whileStatement.Body[0].(ast.IfStatement)
	if !ok {
		t.Fatalf("first while statement = %T, want ast.IfStatement", whileStatement.Body[0])
	}
	if _, ok := firstIf.Body[0].(ast.ContinueStatement); !ok {
		t.Fatalf("nested statement = %T, want ast.ContinueStatement", firstIf.Body[0])
	}

	secondIf, ok := whileStatement.Body[2].(ast.IfStatement)
	if !ok {
		t.Fatalf("third while statement = %T, want ast.IfStatement", whileStatement.Body[2])
	}
	if _, ok := secondIf.Body[0].(ast.BreakStatement); !ok {
		t.Fatalf("nested statement = %T, want ast.BreakStatement", secondIf.Body[0])
	}

	forStatement, ok := program.Statements[1].(ast.ForStatement)
	if !ok {
		t.Fatalf("second statement = %T, want ast.ForStatement", program.Statements[1])
	}
	if forStatement.Target.Name != "item" {
		t.Fatalf("for target = %q, want %q", forStatement.Target.Name, "item")
	}
	if len(forStatement.Body) != 1 {
		t.Fatalf("for body length = %d, want 1", len(forStatement.Body))
	}
}

func TestRejectsLoopControlOutsideLoop(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		message string
	}{
		{name: "break", source: "break\n", message: "break outside loop"},
		{name: "continue", source: "continue\n", message: "continue outside loop"},
		{name: "break in if", source: "if x:\n    break\n", message: "break outside loop"},
		{name: "continue in if", source: "if x:\n    continue\n", message: "continue outside loop"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := lexer.New(test.source).Lex()
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
			if parserError.Message != test.message {
				t.Fatalf("error message = %q, want %q", parserError.Message, test.message)
			}
		})
	}
}
