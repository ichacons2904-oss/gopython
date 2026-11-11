package eval

import (
	"testing"

	"gopython/internal/object"
	"gopython/internal/source"
)

func TestMethodRegistrySupportsOtherRuntimeTypes(t *testing.T) {
	previous, existed := builtinMethodsByType[object.StringType]
	builtinMethodsByType[object.StringType] = map[string]builtinMethod{
		"echo": {
			minArgs: 0,
			maxArgs: 0,
			call: func(receiver object.Value, _ []object.Value, _ source.Position) (object.Value, error) {
				return receiver, nil
			},
		},
	}
	t.Cleanup(func() {
		if existed {
			builtinMethodsByType[object.StringType] = previous
		} else {
			delete(builtinMethodsByType, object.StringType)
		}
	})

	result, err := EvaluateProgram(parseProgram(t, "method = 'hello'.echo\nmethod()\n"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != (object.String{Value: "hello"}) {
		t.Fatalf("result = %#v, want the original receiver", result)
	}
}

func TestBoundMethodKeepsReceiverAfterReassignment(t *testing.T) {
	program := parseProgram(t, "items = [1]\noriginal = items\nadd = items.append\nitems = [99]\nadd(2)\noriginal[1]\n")
	result, err := EvaluateProgram(program, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result != (object.Integer{Value: 2}) {
		t.Fatalf("result = %#v, want 2", result)
	}
}

func TestUnknownAttributeFailsBeforeEvaluatingArguments(t *testing.T) {
	program := parseProgram(t, "items = []\nitems.missing(undefined)\n")
	_, err := EvaluateProgram(program, nil)
	evaluationError, ok := err.(Error)
	if !ok || evaluationError.Kind != AttributeError {
		t.Fatalf("error = %#v, want AttributeError", err)
	}
	if evaluationError.Position != (source.Position{Line: 2, Column: 6}) {
		t.Fatalf("position = %#v, want 2:6", evaluationError.Position)
	}
}
