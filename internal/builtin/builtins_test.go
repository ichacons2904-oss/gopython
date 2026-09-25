package builtin

import (
	"bytes"
	"testing"

	"gopython/internal/object"
)

func TestPrintBuiltin(t *testing.T) {
	var output bytes.Buffer
	environment := object.NewEnvironment(nil)
	Register(environment, &output)

	value, ok := environment.Get("print")
	if !ok {
		t.Fatal("print was not registered")
	}

	print, ok := value.(object.Builtin)
	if !ok {
		t.Fatalf("print = %T, want object.Builtin", value)
	}
	if _, err := print.Call([]object.Value{
		object.String{Value: "answer"},
		object.Integer{Value: 42},
	}); err != nil {
		t.Fatal(err)
	}

	if output.String() != "answer 42\n" {
		t.Fatalf("output = %q, want %q", output.String(), "answer 42\\n")
	}
}

func TestRangeBuiltin(t *testing.T) {
	environment := object.NewEnvironment(nil)
	Register(environment, nil)

	value, ok := environment.Get("range")
	if !ok {
		t.Fatal("range was not registered")
	}
	rangeBuiltin, ok := value.(object.Builtin)
	if !ok {
		t.Fatalf("range = %T, want object.Builtin", value)
	}

	value, err := rangeBuiltin.Call([]object.Value{object.Integer{Value: 3}})
	if err != nil {
		t.Fatal(err)
	}

	rangeValue, ok := value.(object.Range)
	if !ok {
		t.Fatalf("range result = %T, want object.Range", value)
	}
	if rangeValue != (object.Range{Start: 0, Stop: 3, Step: 1}) {
		t.Fatalf("range result = %#v, want range(0, 3, 1)", rangeValue)
	}
}
