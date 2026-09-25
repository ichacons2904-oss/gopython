package eval

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/object"
	"gopython/internal/parser"
)

func TestEvaluateArithmeticAndAssignment(t *testing.T) {
	program := parseProgram(t, "x = 2 + 3 * 4\nx\n")
	environment := object.NewEnvironment(nil)

	result, err := Evaluate(program, environment)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 14}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}

	value, ok := environment.Get("x")
	if !ok {
		t.Fatal("environment does not contain x")
	}
	if value != want {
		t.Fatalf("x = %#v, want %#v", value, want)
	}
}

func TestEvaluateUnaryAndStringConcatenation(t *testing.T) {
	program := parseProgram(t, "x = -7\ny = \"hello\" + \" world\"\ny\n")

	result, err := Evaluate(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.String{Value: "hello world"}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateReportsUndefinedName(t *testing.T) {
	program := parseProgram(t, "missing\n")

	_, err := Evaluate(program, nil)
	if err == nil {
		t.Fatal("Evaluate() returned nil error")
	}

	evaluationError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want eval.Error", err)
	}
	if evaluationError.Kind != NameError {
		t.Fatalf("error kind = %s, want %s", evaluationError.Kind, NameError)
	}
	if evaluationError.Position.Line != 1 || evaluationError.Position.Column != 1 {
		t.Fatalf("error position = %#v, want 1:1", evaluationError.Position)
	}
}

func TestEvaluateReportsDivisionByZero(t *testing.T) {
	program := parseProgram(t, "10 / 0\n")

	_, err := Evaluate(program, nil)
	if err == nil {
		t.Fatal("Evaluate() returned nil error")
	}

	evaluationError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want eval.Error", err)
	}
	if evaluationError.Kind != ZeroDivisionError {
		t.Fatalf("error kind = %s, want %s", evaluationError.Kind, ZeroDivisionError)
	}
}

func TestEvaluateReportsIncompatibleOperands(t *testing.T) {
	program := parseProgram(t, "\"hello\" - \"world\"\n")

	_, err := Evaluate(program, nil)
	if err == nil {
		t.Fatal("Evaluate() returned nil error")
	}

	evaluationError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want eval.Error", err)
	}
	if evaluationError.Kind != TypeError {
		t.Fatalf("error kind = %s, want %s", evaluationError.Kind, TypeError)
	}
}

func parseProgram(t *testing.T, source string) *ast.Program {
	t.Helper()

	tokens, err := lexer.New(source).Lex()
	if err != nil {
		t.Fatal(err)
	}

	program, err := parser.New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	return program
}
