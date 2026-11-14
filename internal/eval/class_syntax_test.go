package eval

import (
	"testing"

	"gopython/internal/source"
)

func TestClassSyntaxReportsExecutionNotImplemented(t *testing.T) {
	tests := []struct {
		input    string
		message  string
		position source.Position
	}{
		{input: "class C: pass\n", message: "class definitions are not supported yet", position: source.Position{Line: 1, Column: 1}},
		{input: "item.value = 1\n", message: "attribute assignment is not supported yet", position: source.Position{Line: 1, Column: 5}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.input), nil)
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != RuntimeError || evaluationError.Message != test.message || evaluationError.Position != test.position {
				t.Fatalf("error = %#v, want RuntimeError %q at %#v", err, test.message, test.position)
			}
		})
	}
}
