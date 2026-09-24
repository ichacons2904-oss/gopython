package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseFunctionDefinition(t *testing.T) {
	source := "def square(value):\n    return value * value\nprint(square(4))\n"
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

	function, ok := program.Statements[0].(ast.FunctionDefinition)
	if !ok {
		t.Fatalf("first statement = %T, want ast.FunctionDefinition", program.Statements[0])
	}
	if function.Name.Name != "square" {
		t.Fatalf("function name = %q, want %q", function.Name.Name, "square")
	}
	if len(function.Parameters) != 1 || function.Parameters[0].Name != "value" {
		t.Fatalf("parameters = %#v, want one parameter named value", function.Parameters)
	}
	if len(function.Body) != 1 {
		t.Fatalf("function body length = %d, want 1", len(function.Body))
	}
	returnStatement, ok := function.Body[0].(ast.ReturnStatement)
	if !ok || returnStatement.Value == nil {
		t.Fatalf("function body statement = %#v, want return with value", function.Body[0])
	}
}

func TestRejectsNestedFunctionDefinitions(t *testing.T) {
	source := "def outer():\n    def inner():\n        return 1\n"
	tokens, err := lexer.New(source).Lex()
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
	if parserError.Message != "nested function definitions are not supported" {
		t.Fatalf("error message = %q", parserError.Message)
	}
}

func TestRejectsFunctionDefinitionInsideConditional(t *testing.T) {
	source := "if condition:\n    def inner():\n        return 1\n"
	tokens, err := lexer.New(source).Lex()
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
	if parserError.Message != "nested function definitions are not supported" {
		t.Fatalf("error message = %q", parserError.Message)
	}
}

func TestRejectsReturnOutsideFunction(t *testing.T) {
	tokens, err := lexer.New("if x:\n    return x\n").Lex()
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
	if parserError.Message != "return outside function" {
		t.Fatalf("error message = %q", parserError.Message)
	}
}
