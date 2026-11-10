package eval

import (
	"testing"

	"gopython/internal/builtin"
	"gopython/internal/object"
)

func TestEvaluateLists(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "literal and negative index", source: "items = [1, 'x', [3]]\nitems[-1][0]\n", want: object.Integer{Value: 3}},
		{name: "mutable aliases", source: "items = [1]\nalias = items\nalias[0] = 4\nitems[0]\n", want: object.Integer{Value: 4}},
		{name: "structural equality", source: "[1, [2]] == [1.0, [2]]\n", want: object.Boolean{Value: true}},
		{name: "iteration and len", source: "items = [1, 2, 3]\ntotal = 0\nfor item in items:\n    total = total + item\nlen(items) + total\n", want: object.Integer{Value: 9}},
		{name: "empty list is falsy", source: "not []\n", want: object.Boolean{Value: true}},
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

func TestEvaluateListIndexErrors(t *testing.T) {
	tests := []struct {
		source string
		kind   ErrorKind
	}{
		{source: "[1][2]\n", kind: IndexError},
		{source: "[1]['0']\n", kind: TypeError},
		{source: "1[0]\n", kind: TypeError},
		{source: "items = [1]\nitems[-2] = 3\n", kind: IndexError},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.source), nil)
			if err == nil {
				t.Fatal("EvaluateProgram() returned nil error")
			}
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind {
				t.Fatalf("error = %#v, want %s", err, test.kind)
			}
		})
	}
}

func TestEvaluateListOperators(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "concatenation", source: "[1, 2] + [3]\n", want: object.NewList([]object.Value{object.Integer{Value: 1}, object.Integer{Value: 2}, object.Integer{Value: 3}})},
		{name: "shallow concatenation", source: "inner = [1]\nleft = [inner]\njoined = left + left\njoined[0][0] = 7\nleft[0][0]\n", want: object.Integer{Value: 7}},
		{name: "list multiplied by integer", source: "['x', 2] * 2\n", want: object.NewList([]object.Value{object.String{Value: "x"}, object.Integer{Value: 2}, object.String{Value: "x"}, object.Integer{Value: 2}})},
		{name: "integer multiplied by list", source: "2 * [3]\n", want: object.NewList([]object.Value{object.Integer{Value: 3}, object.Integer{Value: 3}})},
		{name: "negative repetition is empty", source: "[1] * -2\n", want: object.NewList([]object.Value{})},
		{name: "membership found", source: "2 in [1, 2, 3]\n", want: object.Boolean{Value: true}},
		{name: "membership not found", source: "4 in [1, 2, 3]\n", want: object.Boolean{Value: false}},
		{name: "nested membership", source: "[2] in [[1], [2]]\n", want: object.Boolean{Value: true}},
		{name: "membership precedence", source: "1 + 1 in [2]\n", want: object.Boolean{Value: true}},
		{name: "chained membership comparison", source: "1 in [1] == [1]\n", want: object.Boolean{Value: true}},
		{name: "not membership", source: "not 2 in [1]\n", want: object.Boolean{Value: true}},
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

func TestListOperatorsRejectIncompatibleOperands(t *testing.T) {
	tests := []struct {
		source string
		kind   ErrorKind
	}{
		{source: "[1] + 2\n", kind: TypeError},
		{source: "[1] * 2.0\n", kind: TypeError},
		{source: "[1] * [2]\n", kind: TypeError},
		{source: "1 in 2\n", kind: TypeError},
		{source: "[1, 2] * 9223372036854775807\n", kind: MemoryError},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.source), nil)
			if err == nil {
				t.Fatal("EvaluateProgram() returned nil error")
			}
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind {
				t.Fatalf("error = %#v, want %s", err, test.kind)
			}
		})
	}
}

func TestEvaluateListMethods(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "append and bound receiver", source: "items = [1]\nalias = items\nadd = items.append\nadd(2)\nlen(alias)\n", want: object.Integer{Value: 2}},
		{name: "append returns none", source: "items = []\nitems.append(1)\n", want: object.None{}},
		{name: "extend from list", source: "items = [1]\nitems.extend([2, 3])\nitems[2]\n", want: object.Integer{Value: 3}},
		{name: "extend returns none", source: "items = []\nitems.extend([])\n", want: object.None{}},
		{name: "extend from range", source: "items = [1]\nitems.extend(range(2))\nlen(items)\n", want: object.Integer{Value: 3}},
		{name: "extend itself", source: "items = [1, 2]\nitems.extend(items)\nlen(items)\n", want: object.Integer{Value: 4}},
		{name: "pop last", source: "items = [1, 2]\nitems.pop() + items[0]\n", want: object.Integer{Value: 3}},
		{name: "pop negative index", source: "items = [1, 2, 3]\nitems.pop(-2)\n", want: object.Integer{Value: 2}},
		{name: "clear through alias", source: "items = [1]\nalias = items\nalias.clear()\nlen(items)\n", want: object.Integer{Value: 0}},
		{name: "clear returns none", source: "items = [1]\nitems.clear()\n", want: object.None{}},
		{name: "copy has independent outer list", source: "items = [1]\ncopy = items.copy()\ncopy.append(2)\nlen(items)\n", want: object.Integer{Value: 1}},
		{name: "copy is shallow", source: "inner = [1]\nitems = [inner]\ncopy = items.copy()\ncopy[0][0] = 4\nitems[0][0]\n", want: object.Integer{Value: 4}},
		{name: "chained method calls", source: "items = []\nitems.copy().append(1)\n", want: object.None{}},
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

func TestListMethodErrors(t *testing.T) {
	tests := []struct {
		source string
		kind   ErrorKind
	}{
		{source: "items = []\nitems.nope\n", kind: AttributeError},
		{source: "items = []\nitems.append()\n", kind: TypeError},
		{source: "items = []\nitems.append(1, 2)\n", kind: TypeError},
		{source: "items = []\nitems.extend(1)\n", kind: TypeError},
		{source: "items = []\nitems.pop()\n", kind: IndexError},
		{source: "items = [1]\nitems.pop(1)\n", kind: IndexError},
		{source: "items = [1]\nitems.pop(0.0)\n", kind: TypeError},
		{source: "items = [1]\nitems.clear(1)\n", kind: TypeError},
		{source: "items = [1]\nitems.copy(1)\n", kind: TypeError},
		{source: "(1).append(2)\n", kind: AttributeError},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.source), nil)
			if err == nil {
				t.Fatal("EvaluateProgram() returned nil error")
			}
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind {
				t.Fatalf("error = %#v, want %s", err, test.kind)
			}
		})
	}
}
