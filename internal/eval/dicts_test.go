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

func TestEvaluateDictOperators(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "in checks keys", source: "'a' in {'a': 1}\n", want: object.Boolean{Value: true}},
		{name: "in ignores values", source: "1 in {'a': 1}\n", want: object.Boolean{Value: false}},
		{name: "not in", source: "'b' not in {'a': 1}\n", want: object.Boolean{Value: true}},
		{name: "in with numeric promotion", source: "1.0 in {1: 'x'}\n", want: object.Boolean{Value: true}},
		{name: "equality ignores order", source: "{'a': 1, 'b': 2} == {'b': 2, 'a': 1}\n", want: object.Boolean{Value: true}},
		{name: "equality compares values", source: "{'a': 1} == {'a': 2}\n", want: object.Boolean{Value: false}},
		{name: "nested equality", source: "{'a': [1, {'b': 2}]} == {'a': [1.0, {'b': 2}]}\n", want: object.Boolean{Value: true}},
		{name: "dict is not a list", source: "{} != []\n", want: object.Boolean{Value: true}},
		{name: "self-referential equality", source: "items = {}\nitems['self'] = items\nitems == items\n", want: object.Boolean{Value: true}},
		{name: "empty dict is falsy", source: "not {}\n", want: object.Boolean{Value: true}},
		{name: "non-empty dict is truthy", source: "result = 0\nif {0: 0}:\n    result = 1\nresult\n", want: object.Integer{Value: 1}},
		{name: "len", source: "len({'a': 1, 'b': 2, 'a': 3})\n", want: object.Integer{Value: 2}},
		{name: "for iterates keys in insertion order", source: "text = ''\nfor key in {'b': 1, 'a': 2, 'c': 3}:\n    text = text + key\ntext\n", want: object.String{Value: "bac"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := evaluateDictSource(t, test.source); !valuesEqual(result, test.want) {
				t.Fatalf("result = %#v, want %#v", result, test.want)
			}
		})
	}
}

func TestEvaluateDictOperatorErrors(t *testing.T) {
	tests := []struct {
		source string
		kind   ErrorKind
	}{
		{source: "[1] in {}\n", kind: TypeError},
		{source: "{} < {}\n", kind: TypeError},
		{source: "{} + {}\n", kind: TypeError},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.source), nil)
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind {
				t.Fatalf("error = %#v, want %s", err, test.kind)
			}
		})
	}
}

func TestEvaluateDictMethods(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   object.Value
	}{
		{name: "get existing key", source: "{'a': 1}.get('a')\n", want: object.Integer{Value: 1}},
		{name: "get missing key", source: "{'a': 1}.get('b')\n", want: object.None{}},
		{name: "get with default", source: "{'a': 1}.get('b', 0)\n", want: object.Integer{Value: 0}},
		{name: "keys in insertion order", source: "{'b': 1, 'a': 2}.keys()\n", want: object.NewList([]object.Value{object.String{Value: "b"}, object.String{Value: "a"}})},
		{name: "values in insertion order", source: "{'b': 1, 'a': 2}.values()\n", want: object.NewList([]object.Value{object.Integer{Value: 1}, object.Integer{Value: 2}})},
		{name: "pop returns the value", source: "items = {'a': 1, 'b': 2}\nitems.pop('a') + len(items)\n", want: object.Integer{Value: 2}},
		{name: "pop with default", source: "{}.pop('a', 7)\n", want: object.Integer{Value: 7}},
		{name: "pop then reinsert moves key to the end", source: "items = {'a': 1, 'b': 2}\nitems['a'] = items.pop('a')\nitems.keys()\n", want: object.NewList([]object.Value{object.String{Value: "b"}, object.String{Value: "a"}})},
		{name: "update merges and overwrites", source: "items = {'a': 1}\nitems.update({'a': 3, 'b': 2})\nitems == {'a': 3, 'b': 2}\n", want: object.Boolean{Value: true}},
		{name: "update returns None", source: "{}.update({})\n", want: object.None{}},
		{name: "clear empties the dict", source: "items = {'a': 1}\nalias = items\nitems.clear()\nlen(alias)\n", want: object.Integer{Value: 0}},
		{name: "copy is shallow and independent", source: "inner = [1]\nitems = {'a': inner}\ncopied = items.copy()\ncopied['b'] = 2\ninner.append(2)\n[len(items), copied['a']]\n", want: object.NewList([]object.Value{object.Integer{Value: 1}, object.NewList([]object.Value{object.Integer{Value: 1}, object.Integer{Value: 2}})})},
		{name: "keys returns a new list", source: "items = {'a': 1}\nkeys = items.keys()\nitems['b'] = 2\nlen(keys)\n", want: object.Integer{Value: 1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := evaluateDictSource(t, test.source); !valuesEqual(result, test.want) {
				t.Fatalf("result = %s, want %s", object.Repr(result), object.Repr(test.want))
			}
		})
	}
}

func TestEvaluateDictMethodErrors(t *testing.T) {
	tests := []struct {
		source string
		kind   ErrorKind
	}{
		{source: "{}.pop('a')\n", kind: KeyError},
		{source: "{}.get([1])\n", kind: TypeError},
		{source: "{}.pop([1], 0)\n", kind: TypeError},
		{source: "{}.update([1])\n", kind: TypeError},
		{source: "{}.get()\n", kind: TypeError},
		{source: "{}.keys(1)\n", kind: TypeError},
		{source: "{}.append(1)\n", kind: AttributeError},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.source), nil)
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind {
				t.Fatalf("error = %#v, want %s", err, test.kind)
			}
		})
	}
}

func TestEvaluateDictBuiltin(t *testing.T) {
	result := evaluateDictSource(t, "items = {'a': 1}\ncopied = dict(items)\ncopied['a'] = 2\n[items['a'], copied['a'], dict() == {}]\n")
	want := object.NewList([]object.Value{object.Integer{Value: 1}, object.Integer{Value: 2}, object.Boolean{Value: true}})
	if !valuesEqual(result, want) {
		t.Fatalf("result = %s, want %s", object.Repr(result), object.Repr(want))
	}

	environment := object.NewEnvironment(nil)
	builtin.Register(environment, nil)
	_, err := EvaluateProgram(parseProgram(t, "dict([1])\n"), environment)
	if evaluationError, ok := err.(Error); !ok || evaluationError.Kind != TypeError {
		t.Fatalf("error = %#v, want TypeError", err)
	}
}
