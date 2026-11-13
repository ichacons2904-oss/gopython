package eval

import (
	"testing"

	"gopython/internal/builtin"
	"gopython/internal/object"
)

func TestListBasics(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "not in missing", source: "2 not in [1]\n", want: object.Boolean{Value: true}},
		{name: "not in present", source: "1 not in [1]\n", want: object.Boolean{Value: false}},
		{name: "not in empty", source: "1 not in []\n", want: object.Boolean{Value: true}},
		{name: "nested equality", source: "[1] not in [[1]]\n", want: object.Boolean{Value: false}},
		{name: "not in chain", source: "2 not in [1] == [2]\n", want: object.Boolean{Value: false}},
		{name: "unary not on chain", source: "not 2 in [1] == [2]\n", want: object.Boolean{Value: true}},
		{name: "comparison short circuit", source: "1 not in [1] == missing\n", want: object.Boolean{Value: false}},
		{name: "comparison evaluates middle once", source: "calls = 0\ndef middle():\n    global calls\n    calls = calls + 1\n    return [1]\nresult = 2 not in middle() == [1]\ncalls\n", want: object.Integer{Value: 1}},
		{name: "empty constructor", source: "list()\n", want: object.NewList(nil)},
		{name: "range constructor", source: "list(range(3))\n", want: object.NewList([]object.Value{object.Integer{Value: 0}, object.Integer{Value: 1}, object.Integer{Value: 2}})},
		{name: "constructor copies outer list", source: "a = [1]\nb = list(a)\nb.append(2)\nlen(a)\n", want: object.Integer{Value: 1}},
		{name: "constructor copy is shallow", source: "a = [[1]]\nb = list(a)\nb[0][0] = 7\na[0][0]\n", want: object.Integer{Value: 7}},
		{name: "insert at negative index", source: "a = [1, 2, 3]\na.insert(-1, 9)\na\n", want: object.NewList([]object.Value{object.Integer{Value: 1}, object.Integer{Value: 2}, object.Integer{Value: 9}, object.Integer{Value: 3}})},
		{name: "insert clamps negative", source: "a = [1]\na.insert(-9223372036854775807, 2)\na\n", want: object.NewList([]object.Value{object.Integer{Value: 2}, object.Integer{Value: 1}})},
		{name: "insert clamps positive", source: "a = [1]\na.insert(9223372036854775807, 2)\na\n", want: object.NewList([]object.Value{object.Integer{Value: 1}, object.Integer{Value: 2}})},
		{name: "insert returns none", source: "a = []\na.insert(0, 1)\n", want: object.None{}},
		{name: "remove first equal value", source: "a = [1, 2, 1]\na.remove(1)\na\n", want: object.NewList([]object.Value{object.Integer{Value: 2}, object.Integer{Value: 1}})},
		{name: "remove nested value", source: "a = [[1], [2]]\na.remove([1])\na[0][0]\n", want: object.Integer{Value: 2}},
		{name: "remove returns none", source: "a = [1]\na.remove(1)\n", want: object.None{}},
		{name: "bound insert keeps receiver", source: "a = []\noriginal = a\ninsert = a.insert\na = []\ninsert(0, 3)\noriginal[0]\n", want: object.Integer{Value: 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			environment := object.NewEnvironment(nil)
			builtin.Register(environment, nil)
			result, err := EvaluateProgram(parseProgram(t, test.source), environment)
			if err != nil {
				t.Fatal(err)
			}
			if !valuesEqual(result, test.want) {
				t.Fatalf("result = %#v, want %#v", result, test.want)
			}
		})
	}
}

func TestListBasicsErrors(t *testing.T) {
	tests := []struct {
		source string
		kind   ErrorKind
	}{
		{source: "1 not in 2\n", kind: TypeError},
		{source: "list(1)\n", kind: TypeError},
		{source: "list('abc')\n", kind: TypeError},
		{source: "list([], [])\n", kind: TypeError},
		{source: "a = []\na.insert(1.0, 2)\n", kind: TypeError},
		{source: "a = []\na.insert(1)\n", kind: TypeError},
		{source: "a = []\na.remove()\n", kind: TypeError},
		{source: "a = [1]\na.remove(2)\n", kind: ValueError},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			environment := object.NewEnvironment(nil)
			builtin.Register(environment, nil)
			_, err := EvaluateProgram(parseProgram(t, test.source), environment)
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind {
				t.Fatalf("error = %#v, want %s", err, test.kind)
			}
		})
	}
}

func TestRemoveMissingValueLeavesListUnchanged(t *testing.T) {
	environment := object.NewEnvironment(nil)
	_, err := EvaluateProgram(parseProgram(t, "items = [1, 2]\nitems.remove(3)\n"), environment)
	if evaluationError, ok := err.(Error); !ok || evaluationError.Kind != ValueError {
		t.Fatalf("error = %#v, want ValueError", err)
	}
	value, ok := environment.Get("items")
	if !ok || object.Repr(value) != "[1, 2]" {
		t.Fatalf("items = %#v, want [1, 2]", value)
	}
}
