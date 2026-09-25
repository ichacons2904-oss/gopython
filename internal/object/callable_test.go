package object

import "testing"

func TestBuiltinValue(t *testing.T) {
	builtin := Builtin{
		Name: "answer",
		Function: func([]Value) (Value, error) {
			return Integer{Value: 42}, nil
		},
	}

	if builtin.Type() != BuiltinType {
		t.Fatalf("Type() = %q, want %q", builtin.Type(), BuiltinType)
	}
	if got := builtin.Display(); got != "<built-in function answer>" {
		t.Fatalf("Display() = %q, want built-in description", got)
	}

	value, err := builtin.Function(nil)
	if err != nil {
		t.Fatal(err)
	}
	if value != (Integer{Value: 42}) {
		t.Fatalf("Invoke() = %#v, want Integer{42}", value)
	}
}
