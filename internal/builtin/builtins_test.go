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
	if _, err := print.Function([]object.Value{
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

	value, err := rangeBuiltin.Function([]object.Value{object.Integer{Value: 3}})
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

func TestLenBuiltin(t *testing.T) {
	tests := []struct {
		argument object.Value
		want     int64
	}{
		{argument: object.String{Value: ""}, want: 0},
		{argument: object.String{Value: "ñandú"}, want: 5},
		{argument: object.Range{Start: 0, Stop: 5, Step: 1}, want: 5},
		{argument: object.Range{Start: 1, Stop: 10, Step: 3}, want: 3},
		{argument: object.Range{Start: 5, Stop: 0, Step: -2}, want: 3},
		{argument: object.Range{Start: 5, Stop: 0, Step: 1}, want: 0},
		{argument: object.NewDict(), want: 0},
		{argument: dictOf(object.String{Value: "a"}, object.String{Value: "b"}), want: 2},
	}

	for _, test := range tests {
		got, err := lenFunction([]object.Value{test.argument})
		if err != nil {
			t.Fatal(err)
		}
		if got != (object.Integer{Value: test.want}) {
			t.Fatalf("len(%s) = %#v, want %d", test.argument.Display(), got, test.want)
		}
	}

	if _, err := lenFunction([]object.Value{object.Integer{Value: 1}}); err == nil {
		t.Fatal("len(1) returned nil error")
	}
}

func dictOf(keys ...object.Value) *object.Dict {
	dict := object.NewDict()
	for _, key := range keys {
		dict.Set(key, object.None{})
	}
	return dict
}

func TestDictBuiltin(t *testing.T) {
	empty, err := dictFunction(nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := empty.Display(); got != "{}" {
		t.Fatalf("dict() = %s, want {}", got)
	}

	source := dictOf(object.String{Value: "a"})
	copied, err := dictFunction([]object.Value{source})
	if err != nil {
		t.Fatal(err)
	}
	copied.(*object.Dict).Set(object.String{Value: "b"}, object.None{})
	if source.Len() != 1 || copied.(*object.Dict).Len() != 2 {
		t.Fatal("dict(source) did not return an independent copy")
	}

	invalid := [][]object.Value{
		{object.NewList(nil)},
		{object.NewDict(), object.NewDict()},
	}
	for _, arguments := range invalid {
		if _, err := dictFunction(arguments); err == nil {
			t.Fatalf("dict(%d arguments) returned nil error", len(arguments))
		}
	}
}
