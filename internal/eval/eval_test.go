package eval

import (
	"bytes"
	"testing"

	"gopython/internal/ast"
	"gopython/internal/builtin"
	"gopython/internal/lexer"
	"gopython/internal/object"
	"gopython/internal/parser"
)

func TestEvaluateArithmeticAndAssignment(t *testing.T) {
	program := parseProgram(t, "x = 2 + 3 * 4\nx\n")
	environment := object.NewEnvironment(nil)

	result, err := EvaluateProgram(program, environment)
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

	result, err := EvaluateProgram(program, nil)
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

	_, err := EvaluateProgram(program, nil)
	if err == nil {
		t.Fatal("EvaluateProgram() returned nil error")
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

	_, err := EvaluateProgram(program, nil)
	if err == nil {
		t.Fatal("EvaluateProgram() returned nil error")
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

	_, err := EvaluateProgram(program, nil)
	if err == nil {
		t.Fatal("EvaluateProgram() returned nil error")
	}

	evaluationError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want eval.Error", err)
	}
	if evaluationError.Kind != TypeError {
		t.Fatalf("error kind = %s, want %s", evaluationError.Kind, TypeError)
	}
}

func TestEvaluateComparisons(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{name: "less than", source: "1 < 2\n", want: true},
		{name: "equal", source: "2 == 2\n", want: true},
		{name: "not equal", source: "2 != 3\n", want: true},
		{name: "chained true", source: "1 < 2 < 3\n", want: true},
		{name: "chained false", source: "1 < 2 < 1\n", want: false},
		{name: "different types", source: "1 == \"1\"\n", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := EvaluateProgram(parseProgram(t, test.source), nil)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := result.(object.Boolean)
			if !ok {
				t.Fatalf("result = %T, want object.Boolean", result)
			}
			if got.Value != test.want {
				t.Fatalf("result = %v, want %v", got.Value, test.want)
			}
		})
	}
}

func TestEvaluateNotUsesTruthiness(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{name: "zero", source: "not 0\n", want: true},
		{name: "empty string", source: "not \"\"\n", want: true},
		{name: "none", source: "not None\n", want: true},
		{name: "nonzero integer", source: "not 1\n", want: false},
		{name: "comparison precedence", source: "not 1 == 2\n", want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := EvaluateProgram(parseProgram(t, test.source), nil)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := result.(object.Boolean)
			if !ok {
				t.Fatalf("result = %T, want object.Boolean", result)
			}
			if got.Value != test.want {
				t.Fatalf("result = %v, want %v", got.Value, test.want)
			}
		})
	}
}

func TestEvaluateLogicalOperatorsShortCircuit(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "false and", source: "False and missing\n", want: object.Boolean{Value: false}},
		{name: "true or", source: "True or missing\n", want: object.Boolean{Value: true}},
		{name: "or returns right operand", source: "0 or 5\n", want: object.Integer{Value: 5}},
		{name: "and returns right operand", source: "1 and 5\n", want: object.Integer{Value: 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := EvaluateProgram(parseProgram(t, test.source), nil)
			if err != nil {
				t.Fatal(err)
			}
			if result != test.want {
				t.Fatalf("result = %#v, want %#v", result, test.want)
			}
		})
	}
}

func TestEvaluateReportsIncompatibleComparison(t *testing.T) {
	program := parseProgram(t, "1 < \"one\"\n")

	_, err := EvaluateProgram(program, nil)
	if err == nil {
		t.Fatal("EvaluateProgram() returned nil error")
	}

	evaluationError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want eval.Error", err)
	}
	if evaluationError.Kind != TypeError {
		t.Fatalf("error kind = %s, want %s", evaluationError.Kind, TypeError)
	}
}

func TestEvaluateIfElifElse(t *testing.T) {
	program := parseProgram(t, "if value < 0:\n    result = \"negative\"\nelif value == 0:\n    result = \"zero\"\nelse:\n    result = \"positive\"\nresult\n")

	tests := []struct {
		name  string
		value int64
		want  string
	}{
		{name: "if branch", value: -1, want: "negative"},
		{name: "elif branch", value: 0, want: "zero"},
		{name: "else branch", value: 1, want: "positive"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			environment := object.NewEnvironment(nil)
			environment.Set("value", object.Integer{Value: test.value})

			result, err := EvaluateProgram(program, environment)
			if err != nil {
				t.Fatal(err)
			}

			want := object.String{Value: test.want}
			if result != want {
				t.Fatalf("result = %#v, want %#v", result, want)
			}
		})
	}
}

func TestEvaluateIfSkipsInactiveBranch(t *testing.T) {
	program := parseProgram(t, "if False:\n    missing\nelse:\n    42\n")

	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 42}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateNestedIf(t *testing.T) {
	program := parseProgram(t, "if True:\n    if False:\n        1\n    else:\n        2\n")

	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 2}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateWhile(t *testing.T) {
	program := parseProgram(t, "x = 0\nwhile x < 3:\n    x = x + 1\nx\n")

	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 3}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateWhileSkipsFalseBody(t *testing.T) {
	program := parseProgram(t, "while False:\n    missing\n42\n")

	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 42}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateBreakExitsNearestLoop(t *testing.T) {
	program := parseProgram(t, "x = 0\nwhile True:\n    x = x + 1\n    if x == 3:\n        break\nx\n")

	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 3}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateContinueStartsNextIteration(t *testing.T) {
	program := parseProgram(t, "x = 0\ntotal = 0\nwhile x < 5:\n    x = x + 1\n    if x == 3:\n        continue\n    total = total + x\ntotal\n")

	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 12}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateForOverRange(t *testing.T) {
	program := parseProgram(t, "total = 0\nfor value in range(1, 5):\n    total = total + value\ntotal\n")
	environment := object.NewEnvironment(nil)
	builtin.Register(environment, nil)

	result, err := EvaluateProgram(program, environment)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 10}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateForBreakAndContinue(t *testing.T) {
	program := parseProgram(t, "total = 0\nfor value in range(5):\n    if value == 2:\n        continue\n    if value == 4:\n        break\n    total = total + value\ntotal\n")
	environment := object.NewEnvironment(nil)
	builtin.Register(environment, nil)

	result, err := EvaluateProgram(program, environment)
	if err != nil {
		t.Fatal(err)
	}

	want := object.Integer{Value: 4}
	if result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestEvaluateForRejectsNonIterable(t *testing.T) {
	program := parseProgram(t, "for value in 42:\n    value\n")

	_, err := EvaluateProgram(program, nil)
	if err == nil {
		t.Fatal("EvaluateProgram() returned nil error")
	}

	evaluationError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want eval.Error", err)
	}
	if evaluationError.Kind != TypeError {
		t.Fatalf("error kind = %s, want %s", evaluationError.Kind, TypeError)
	}
}

func TestEvaluatePrintAndRange(t *testing.T) {
	program := parseProgram(t, "print(\"answer\", 42)\nrange(1, 5, 2)\n")
	environment := object.NewEnvironment(nil)
	var output bytes.Buffer
	builtin.Register(environment, &output)

	result, err := EvaluateProgram(program, environment)
	if err != nil {
		t.Fatal(err)
	}

	wantRange := object.Range{Start: 1, Stop: 5, Step: 2}
	if result != wantRange {
		t.Fatalf("result = %#v, want %#v", result, wantRange)
	}
	if output.String() != "answer 42\n" {
		t.Fatalf("output = %q, want %q", output.String(), "answer 42\\n")
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
