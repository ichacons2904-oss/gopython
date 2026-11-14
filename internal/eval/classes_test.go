package eval

import (
	"testing"

	"gopython/internal/builtin"
	"gopython/internal/object"
	"gopython/internal/source"
)

func TestClassDefinitionsAndAttributes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  object.Value
	}{
		{name: "class body namespace", input: "outside = 4\nclass C:\n    value = outside + 1\nC.value\n", want: object.Integer{Value: 5}},
		{name: "class name does not leak", input: "value = 1\nclass C:\n    value = 2\nvalue\n", want: object.Integer{Value: 1}},
		{name: "class body executes once", input: "count = 0\nclass C:\n    global count\n    count = count + 1\ncount\n", want: object.Integer{Value: 1}},
		{name: "class alias identity", input: "class C: pass\nalias = C\nC == alias\n", want: object.Boolean{Value: true}},
		{name: "distinct class identity", input: "class C: pass\nclass D: pass\nC != D\n", want: object.Boolean{Value: true}},
		{name: "inherited attribute", input: "class A:\n    value = 1\nclass B(A): pass\nB.value\n", want: object.Integer{Value: 1}},
		{name: "child override", input: "class A:\n    value = 1\nclass B(A):\n    value = 2\nB.value\n", want: object.Integer{Value: 2}},
		{name: "later base fallback", input: "class A: pass\nclass B:\n    value = 2\nclass C(A, B): pass\nC.value\n", want: object.Integer{Value: 2}},
		{name: "depth-first diamond", input: "class A:\n    value = 1\nclass B(A): pass\nclass C(A):\n    value = 2\nclass D(B, C): pass\nD.value\n", want: object.Integer{Value: 1}},
		{name: "base updates remain visible", input: "class A:\n    value = 1\nclass B(A): pass\nA.value = 3\nB.value\n", want: object.Integer{Value: 3}},
		{name: "method sees global not class namespace", input: "value = 3\nclass C:\n    value = 10\n    def read(self):\n        return value\nC.read(None)\n", want: object.Integer{Value: 3}},
		{name: "method captures enclosing function", input: "def make():\n    value = 7\n    class C:\n        value = 9\n        def read(self):\n            return value\n    return C\nC = make()\nC.read(None)\n", want: object.Integer{Value: 7}},
		{name: "nested class skips outer class scope", input: "value = 1\nclass Outer:\n    value = 2\n    class Inner:\n        observed = value\nOuter.Inner.observed\n", want: object.Integer{Value: 1}},
		{name: "nested class method skips both class scopes", input: "value = 1\nclass Outer:\n    value = 2\n    class Inner:\n        value = 3\n        def read(self):\n            return value\nOuter.Inner.read(None)\n", want: object.Integer{Value: 1}},
		{name: "nested class resolves base in outer class", input: "class Outer:\n    class A:\n        value = 3\n    class B(A): pass\nOuter.B.value\n", want: object.Integer{Value: 3}},
		{name: "nonlocal in class body", input: "def outer():\n    value = 1\n    class C:\n        nonlocal value\n        value = 2\n    return value\nouter()\n", want: object.Integer{Value: 2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := EvaluateProgram(parseProgram(t, test.input), nil)
			if err != nil {
				t.Fatal(err)
			}
			if !valuesEqual(result, test.want) {
				t.Fatalf("result = %#v, want %#v", result, test.want)
			}
		})
	}
}

func TestInstanceAttributesAndIdentity(t *testing.T) {
	environment := object.NewEnvironment(nil)
	_, err := EvaluateProgram(parseProgram(t, "class Base:\n    value = 10\nclass Child(Base): pass\n"), environment)
	if err != nil {
		t.Fatal(err)
	}
	value, _ := environment.Get("Child")
	class := value.(*object.Class)
	first := object.NewInstance(class)
	second := object.NewInstance(class)
	environment.Set("first", first)
	environment.Set("second", second)
	environment.Set("alias", first)

	checks := []struct {
		input string
		want  object.Value
	}{
		{input: "first.value = 20\nsecond.value\n", want: object.Integer{Value: 10}},
		{input: "Child.value = 30\nsecond.value\n", want: object.Integer{Value: 30}},
		{input: "first.value\n", want: object.Integer{Value: 20}},
		{input: "first == alias\n", want: object.Boolean{Value: true}},
		{input: "first == second\n", want: object.Boolean{Value: false}},
	}
	for _, check := range checks {
		result, err := EvaluateProgram(parseProgram(t, check.input), environment)
		if err != nil {
			t.Fatal(err)
		}
		if !valuesEqual(result, check.want) {
			t.Fatalf("%q = %#v, want %#v", check.input, result, check.want)
		}
	}
}

func TestAttributeAssignmentEvaluatesValueBeforeReceiver(t *testing.T) {
	environment := object.NewEnvironment(nil)
	environment.Set("item", object.NewInstance(object.NewClass("C", nil, nil)))
	input := "events = 0\ndef value():\n    global events\n    events = events * 10 + 1\n    return 7\ndef target():\n    global events\n    events = events * 10 + 2\n    return item\ntarget().field = value()\nevents\n"
	result, err := EvaluateProgram(parseProgram(t, input), environment)
	if err != nil {
		t.Fatal(err)
	}
	if result != (object.Integer{Value: 12}) {
		t.Fatalf("events = %#v, want 12", result)
	}
}

func TestClassErrors(t *testing.T) {
	tests := []struct {
		input    string
		kind     ErrorKind
		position source.Position
	}{
		{input: "class C(Missing): pass\n", kind: NameError, position: source.Position{Line: 1, Column: 9}},
		{input: "value = 1\nclass C(value): pass\n", kind: TypeError, position: source.Position{Line: 2, Column: 9}},
		{input: "class A: pass\nalias = A\nclass C(A, alias): pass\n", kind: TypeError, position: source.Position{Line: 3, Column: 12}},
		{input: "class C: pass\nC.missing\n", kind: AttributeError, position: source.Position{Line: 2, Column: 2}},
		{input: "value = 1\nclass C:\n    global value\n    value = 2\nC.value\n", kind: AttributeError, position: source.Position{Line: 5, Column: 2}},
		{input: "(1).field = 2\n", kind: AttributeError, position: source.Position{Line: 1, Column: 4}},
		{input: "class C: pass\nC(1)\n", kind: TypeError, position: source.Position{Line: 2, Column: 2}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, test.input), nil)
			evaluationError, ok := err.(Error)
			if !ok || evaluationError.Kind != test.kind || evaluationError.Position != test.position {
				t.Fatalf("error = %#v, want %s at %#v", err, test.kind, test.position)
			}
		})
	}
}

func TestClassBodyErrorDoesNotReplaceExistingName(t *testing.T) {
	environment := object.NewEnvironment(nil)
	environment.Set("C", object.Integer{Value: 5})
	_, err := EvaluateProgram(parseProgram(t, "class C:\n    value = missing\n"), environment)
	if evaluationError, ok := err.(Error); !ok || evaluationError.Kind != NameError {
		t.Fatalf("error = %#v, want NameError", err)
	}
	value, _ := environment.Get("C")
	if value != (object.Integer{Value: 5}) {
		t.Fatalf("C = %#v, want original value 5", value)
	}
}

func TestConstructionAndUserMethods(t *testing.T) {
	tests := []struct {
		name, input string
		want        object.Value
	}{
		{"empty construction", "class C: pass\na = C()\nb = C()\na != b\n", object.Boolean{Value: true}},
		{"initializer and isolated instances", "class C:\n    def __init__(self, value):\n        self.value = value\na = C(1)\nb = C(2)\nalias = a\nalias.value = 3\na.value + b.value\n", object.Integer{Value: 5}},
		{"inherited initializer and method", "class A:\n    def __init__(self, x):\n        self.x = x\n    def read(self):\n        return self.x\nclass B(A): pass\nB(7).read()\n", object.Integer{Value: 7}},
		{"override initializer", "class A:\n    def __init__(self):\n        self.x = 1\nclass B(A):\n    def __init__(self):\n        self.x = 2\nB().x\n", object.Integer{Value: 2}},
		{"saved method preserves function and receiver", "class C:\n    def get(self):\n        return self.x\na = C()\na.x = 4\nsaved = a.get\ndef replacement(self):\n    return 99\nC.get = replacement\na = C()\nsaved()\n", object.Integer{Value: 4}},
		{"instance function is not bound", "def f():\n    return 8\nclass C:\n    def f(self):\n        return 1\na = C()\na.f = f\na.f()\n", object.Integer{Value: 8}},
		{"explicit class method", "class C:\n    def get(self):\n        return self.x\na = C()\na.x = 9\nC.get(a)\n", object.Integer{Value: 9}},
		{"builtin methods with instances", "class C: pass\na = C()\na.items = []\na.items.append(4)\na.data = {}\na.data.update({1: 3})\na.items[0] + a.data[1]\n", object.Integer{Value: 7}},
		{"inherited method follows superfather", "class A:\n    def value(self):\n        return 1\nclass B(A): pass\nclass C(A):\n    def value(self):\n        return 2\nclass D(B, C): pass\nD().value()\n", object.Integer{Value: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateProgram(parseProgram(t, tt.input), nil)
			if err != nil {
				t.Fatal(err)
			}
			if !valuesEqual(got, tt.want) {
				t.Fatalf("got %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestConstructionAndMethodErrors(t *testing.T) {
	tests := []struct {
		name, input string
		kind        ErrorKind
	}{
		{"missing initializer argument", "class C:\n    def __init__(self, x): pass\nC()\n", TypeError},
		{"extra initializer argument", "class C:\n    def __init__(self): pass\nC(1)\n", TypeError},
		{"initializer returns value", "class C:\n    def __init__(self):\n        return 1\nC()\n", TypeError},
		{"initializer error propagates", "class C:\n    def __init__(self):\n        self.x = missing\nC()\n", NameError},
		{"noncallable initializer", "class C:\n    __init__ = 1\nC()\n", TypeError},
		{"method missing argument", "class C:\n    def f(self, x): pass\nC().f()\n", TypeError},
		{"instance shadows method", "class C:\n    def f(self): pass\na = C()\na.f = 2\na.f()\n", TypeError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EvaluateProgram(parseProgram(t, tt.input), nil)
			e, ok := err.(Error)
			if !ok || e.Kind != tt.kind {
				t.Fatalf("got %#v; want %s", err, tt.kind)
			}
		})
	}
}

func TestBuiltinConstructorsCannotBeBases(t *testing.T) {
	tests := []struct {
		name, input string
		position    source.Position
	}{
		{"list", "class C(list): pass\n", source.Position{Line: 1, Column: 9}},
		{"dict", "class C(dict): pass\n", source.Position{Line: 1, Column: 9}},
		{"aliased list", "base = list\nclass C(base): pass\n", source.Position{Line: 2, Column: 9}},
		{"mixed bases", "class A: pass\nclass C(A, dict): pass\n", source.Position{Line: 2, Column: 12}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := object.NewEnvironment(nil)
			builtin.Register(environment, nil)
			_, err := EvaluateProgram(parseProgram(t, tt.input), environment)
			e, ok := err.(Error)
			if !ok || e.Kind != TypeError || e.Position != tt.position {
				t.Fatalf("error = %#v, want TypeError at %#v", err, tt.position)
			}
		})
	}
}

func TestBuiltinBaseRejectedBeforeClassBody(t *testing.T) {
	environment := object.NewEnvironment(nil)
	builtin.Register(environment, nil)
	_, err := EvaluateProgram(parseProgram(t, "events = []\nclass C(list):\n    events.append(1)\n"), environment)
	if e, ok := err.(Error); !ok || e.Kind != TypeError {
		t.Fatalf("error = %#v, want TypeError", err)
	}
	result, err := EvaluateProgram(parseProgram(t, "len(events)\n"), environment)
	if err != nil || result != (object.Integer{Value: 0}) {
		t.Fatalf("events length = %#v, %v; want 0", result, err)
	}
}

func TestBuiltinNameCanBeReboundToClassBase(t *testing.T) {
	result, err := EvaluateProgram(parseProgram(t, "class Parent:\n    value = 7\nlist = Parent\nclass Child(list): pass\nChild().value\n"), nil)
	if err != nil || result != (object.Integer{Value: 7}) {
		t.Fatalf("result = %#v, %v; want 7", result, err)
	}
}
