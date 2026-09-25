package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseSimpleProgram(t *testing.T) {
	tokens, err := lexer.New("x = 2 + 3 * 4\nprint(x)\n").Lex()
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

	assignment, ok := program.Statements[0].(*ast.Assignment)
	if !ok {
		t.Fatalf("first statement = %T, want ast.Assignment", program.Statements[0])
	}
	if assignment.Name.Name != "x" {
		t.Fatalf("assignment name = %q, want %q", assignment.Name.Name, "x")
	}

	addition, ok := assignment.Value.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("assignment value = %T, want ast.BinaryExpression", assignment.Value)
	}
	if addition.Operator != lexer.Plus {
		t.Fatalf("outer operator = %s, want %s", addition.Operator, lexer.Plus)
	}

	multiplication, ok := addition.Right.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("right expression = %T, want ast.BinaryExpression", addition.Right)
	}
	if multiplication.Operator != lexer.Asterisk {
		t.Fatalf("inner operator = %s, want %s", multiplication.Operator, lexer.Asterisk)
	}

	expressionStatement, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("second statement = %T, want ast.ExpressionStatement", program.Statements[1])
	}
	call, ok := expressionStatement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression = %T, want ast.CallExpression", expressionStatement.Expression)
	}
	if len(call.Arguments) != 1 {
		t.Fatalf("argument count = %d, want 1", len(call.Arguments))
	}
}

func TestParseSimpleStatements(t *testing.T) {
	tokens, err := lexer.New("def f():\n    return\n").Lex()
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

	function := program.Statements[0].(*ast.FunctionDefinition)
	returnWithoutValue := function.Body[0].(*ast.ReturnStatement)
	if returnWithoutValue.Value != nil {
		t.Fatal("return without expression has a value")
	}
}

func TestParseReportsExpressionErrors(t *testing.T) {
	tokens, err := lexer.New("x =\n").Lex()
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
	if parserError.Token.Type != lexer.Newline {
		t.Fatalf("error token = %s, want %s", parserError.Token.Type, lexer.Newline)
	}
}
