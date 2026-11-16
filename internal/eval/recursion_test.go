package eval

import (
	"testing"

	"gopython/internal/object"
)

func TestInfiniteRecursionRaisesRecursionError(t *testing.T) {
	tests := []string{
		"def f():\n    return f()\nf()\n",
		"def even(n):\n    return odd(n)\ndef odd(n):\n    return even(n)\neven(1)\n",
		"class A:\n    def __init__(self):\n        A()\nA()\n",
		"class A:\n    def again(self):\n        return self.again()\nA().again()\n",
	}
	for _, source := range tests {
		t.Run(source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, source), object.NewEnvironment(nil))
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != RecursionError {
				t.Fatalf("error = %#v, want RecursionError", err)
			}
		})
	}
}

func TestRecursionBelowLimitSucceeds(t *testing.T) {
	source := "def depth(n):\n    if n == 0:\n        return 0\n    return 1 + depth(n - 1)\ndepth(990)\n"
	result, err := EvaluateProgram(parseProgram(t, source), object.NewEnvironment(nil))
	if err != nil {
		t.Fatal(err)
	}
	if want := (object.Integer{Value: 990}); result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestRecursionErrorResetsCallDepth(t *testing.T) {
	environment := object.NewEnvironment(nil)
	definition := "def depth(n):\n    if n == 0:\n        return 0\n    return 1 + depth(n - 1)\n"
	if _, err := EvaluateProgram(parseProgram(t, definition), environment); err != nil {
		t.Fatal(err)
	}

	if _, err := EvaluateProgram(parseProgram(t, "depth(5000)\n"), environment); err == nil {
		t.Fatal("depth(5000) succeeded, want RecursionError")
	}

	result, err := EvaluateProgram(parseProgram(t, "depth(990)\n"), environment)
	if err != nil {
		t.Fatalf("depth(990) after RecursionError: %v", err)
	}
	if want := (object.Integer{Value: 990}); result != want {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}
