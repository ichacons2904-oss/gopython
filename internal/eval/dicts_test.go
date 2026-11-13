package eval

import (
	"testing"

	"gopython/internal/builtin"
	"gopython/internal/object"
)

func evaluateDictSource(t *testing.T, source string) object.Value {
	t.Helper()
	environment := object.NewEnvironment(nil)
	builtin.Register(environment, nil)
	result, err := EvaluateProgram(parseProgram(t, source), environment)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestEvaluateDictIndexing(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "literal lookup", source: "{'a': 1, 'b': 2}['b']\n", want: object.Integer{Value: 2}},
		{name: "later duplicate wins", source: "{'a': 1, 'a': 3}['a']\n", want: object.Integer{Value: 3}},
		{name: "int and float share a key", source: "{1: 'int'}[1.0]\n", want: object.String{Value: "int"}},
		{name: "None as key", source: "{None: 0}[None]\n", want: object.Integer{Value: 0}},
		{name: "assignment adds a key", source: "ages = {}\nages['ana'] = 30\nages['ana']\n", want: object.Integer{Value: 30}},
		{name: "assignment replaces a value", source: "ages = {'ana': 30}\nages['ana'] = 31\nages['ana']\n", want: object.Integer{Value: 31}},
		{name: "mutable aliases", source: "items = {}\nalias = items\nalias[1] = 'x'\nitems[1]\n", want: object.String{Value: "x"}},
		{name: "nested containers", source: "data = {'list': [1, {'k': 2}]}\ndata['list'][1]['k']\n", want: object.Integer{Value: 2}},
		{name: "computed key", source: "key = 'a' + 'b'\n{'ab': 5}[key]\n", want: object.Integer{Value: 5}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := evaluateDictSource(t, test.source); !valuesEqual(result, test.want) {
				t.Fatalf("result = %#v, want %#v", result, test.want)
			}
		})
	}
}

func TestEvaluateDictDisplay(t *testing.T) {
	result := evaluateDictSource(t, "{'a': [1, 'b'], 2: None, 'c': {}}\n")
	if got := object.Repr(result); got != "{'a': [1, 'b'], 2: None, 'c': {}}" {
		t.Fatalf("Repr() = %q", got)
	}
}

func TestEvaluateDictErrors(t *testing.T) {
	tests := []struct {
		source  string
		kind    ErrorKind
		message string
	}{
		{source: "{'a': 1}['b']\n", kind: KeyError, message: "'b'"},
		{source: "{1: 1}[True]\n", kind: KeyError, message: "True"},
		{source: "{[1]: 2}\n", kind: TypeError, message: "unhashable type: LIST"},
		{source: "{}[{}]\n", kind: TypeError, message: "unhashable type: DICT"},
		{source: "items = {}\nitems[[1]] = 2\n", kind: TypeError, message: "unhashable type: LIST"},
		{source: "items = 1\nitems['a'] = 2\n", kind: TypeError, message: "INTEGER object is not subscriptable"},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.source), nil)
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind || evaluationError.Message != test.message {
				t.Fatalf("error = %#v, want %s: %s", err, test.kind, test.message)
			}
		})
	}
}
