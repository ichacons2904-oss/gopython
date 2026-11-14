package parser

import (
	"reflect"
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/source"
)

func parseClassSource(t *testing.T, input string) (*ast.Program, error) {
	t.Helper()
	tokens, err := lexer.New(input).Lex()
	if err != nil {
		t.Fatalf("Lex() = %v", err)
	}
	return New(tokens).Parse()
}

func TestParseClassHeaders(t *testing.T) {
	tests := []struct {
		input string
		bases []string
	}{
		{input: "class C:\n    pass\n", bases: []string{}},
		{input: "class C(): pass\n", bases: []string{}},
		{input: "class C(Base): pass\n", bases: []string{"Base"}},
		{input: "class C(Left, Right,): pass\n", bases: []string{"Left", "Right"}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program, err := parseClassSource(t, test.input)
			if err != nil {
				t.Fatal(err)
			}
			class, ok := program.Statements[0].(*ast.ClassDefinition)
			if !ok {
				t.Fatalf("statement = %T, want *ast.ClassDefinition", program.Statements[0])
			}
			if class.Name.Name != "C" || class.Position() != (source.Position{Line: 1, Column: 1}) {
				t.Fatalf("class = %#v, want C at 1:1", class)
			}
			got := make([]string, len(class.Bases))
			for index, base := range class.Bases {
				got[index] = base.Name
			}
			if !reflect.DeepEqual(got, test.bases) {
				t.Fatalf("bases = %v, want %v", got, test.bases)
			}
			if len(class.Body) != 1 {
				t.Fatalf("body length = %d, want 1", len(class.Body))
			}
			if _, ok := class.Body[0].(*ast.PassStatement); !ok {
				t.Fatalf("body statement = %T, want *ast.PassStatement", class.Body[0])
			}
		})
	}
}

func TestParseClassMethodsAndAttributeAssignment(t *testing.T) {
	input := "class Child(Left, Right):\n    def __init__(self, value):\n        self.value = value\n"
	program, err := parseClassSource(t, input)
	if err != nil {
		t.Fatal(err)
	}
	class := program.Statements[0].(*ast.ClassDefinition)
	if class.Name.Position() != (source.Position{Line: 1, Column: 7}) ||
		class.Bases[0].Position() != (source.Position{Line: 1, Column: 13}) ||
		class.Bases[1].Position() != (source.Position{Line: 1, Column: 19}) {
		t.Fatalf("name/base positions = %#v, %#v, %#v", class.Name.Position(), class.Bases[0].Position(), class.Bases[1].Position())
	}
	method, ok := class.Body[0].(*ast.FunctionDefinition)
	if !ok || method.Name.Name != "__init__" || len(method.Body) != 1 {
		t.Fatalf("class body = %#v, want __init__ method", class.Body)
	}
	assignment, ok := method.Body[0].(*ast.AttributeAssignment)
	if !ok {
		t.Fatalf("method body = %T, want *ast.AttributeAssignment", method.Body[0])
	}
	receiver, ok := assignment.Target.Object.(*ast.Identifier)
	if !ok || receiver.Name != "self" || assignment.Target.Name != "value" || assignment.Position() != (source.Position{Line: 3, Column: 13}) {
		t.Fatalf("assignment target = %#v, want self.value at 3:13", assignment.Target)
	}
	if value, ok := assignment.Value.(*ast.Identifier); !ok || value.Name != "value" {
		t.Fatalf("assignment value = %#v, want value identifier", assignment.Value)
	}
}

func TestParseAttributeAssignmentTargets(t *testing.T) {
	program, err := parseClassSource(t, "factory().item = 1\nnodes[0].next = 2\nobject.inner.value = 3\n")
	if err != nil {
		t.Fatal(err)
	}
	first := program.Statements[0].(*ast.AttributeAssignment)
	if _, ok := first.Target.Object.(*ast.CallExpression); !ok || first.Target.Name != "item" {
		t.Fatalf("first target = %#v", first.Target)
	}
	second := program.Statements[1].(*ast.AttributeAssignment)
	if _, ok := second.Target.Object.(*ast.IndexExpression); !ok || second.Target.Name != "next" {
		t.Fatalf("second target = %#v", second.Target)
	}
	third := program.Statements[2].(*ast.AttributeAssignment)
	if _, ok := third.Target.Object.(*ast.AttributeExpression); !ok || third.Target.Name != "value" {
		t.Fatalf("third target = %#v", third.Target)
	}
}

func TestParseDictionaryInClassMethod(t *testing.T) {
	program, err := parseClassSource(t, "class C:\n    def __init__(self):\n        self.data = {'count': 1}\n")
	if err != nil {
		t.Fatal(err)
	}
	class, ok := program.Statements[0].(*ast.ClassDefinition)
	if !ok {
		t.Fatalf("statement = %T, want class", program.Statements[0])
	}
	method, ok := class.Body[0].(*ast.FunctionDefinition)
	if !ok {
		t.Fatalf("class body = %T, want method", class.Body[0])
	}
	assignment, ok := method.Body[0].(*ast.AttributeAssignment)
	if !ok || assignment.Target.Name != "data" {
		t.Fatalf("method body = %#v, want self.data assignment", method.Body[0])
	}
	dict, ok := assignment.Value.(*ast.DictLiteral)
	if !ok || len(dict.Entries) != 1 {
		t.Fatalf("assignment value = %#v, want dictionary with one entry", assignment.Value)
	}
}

func TestClassBodiesHaveTheirOwnControlFlowContext(t *testing.T) {
	valid := []string{
		"class Outer:\n    class Inner: pass\n",
		"def outer():\n    value = 1\n    class C:\n        nonlocal value\n        def method(self):\n            return value\n    return C\n",
		"while True:\n    class C:\n        while True:\n            break\n    break\n",
	}
	for _, input := range valid {
		t.Run(input, func(t *testing.T) {
			if _, err := parseClassSource(t, input); err != nil {
				t.Fatalf("Parse() = %v", err)
			}
		})
	}
	invalid := []struct {
		input    string
		message  string
		position source.Position
	}{
		{input: "def outer():\n    class C:\n        return 1\n", message: "return outside function", position: source.Position{Line: 3, Column: 9}},
		{input: "while True:\n    class C:\n        break\n", message: "break outside loop", position: source.Position{Line: 3, Column: 9}},
		{input: "while True:\n    class C:\n        continue\n", message: "continue outside loop", position: source.Position{Line: 3, Column: 9}},
		{input: "class C:\n    nonlocal value\n", message: "nonlocal declaration not allowed at module level", position: source.Position{Line: 2, Column: 5}},
	}
	for _, test := range invalid {
		t.Run(test.input, func(t *testing.T) {
			_, err := parseClassSource(t, test.input)
			parseError, ok := err.(Error)
			if !ok || parseError.Message != test.message || positionOf(parseError.Token) != test.position {
				t.Fatalf("error = %#v, want %q at %#v", err, test.message, test.position)
			}
		})
	}
}

func TestRejectInvalidClassHeadersAndAssignmentTargets(t *testing.T) {
	invalid := []string{
		"class : pass\n",
		"class C(A B): pass\n",
		"class C(A,, B): pass\n",
		"class C(A.x): pass\n",
		"class C(factory()): pass\n",
		"class C\n    pass\n",
		"class C:\npass\n",
		"factory() = 1\n",
	}
	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			_, err := parseClassSource(t, input)
			if parseError, ok := err.(Error); !ok || parseError.Token.Line < 1 || parseError.Token.Column < 1 {
				t.Fatalf("error = %#v, want parser.Error", err)
			}
		})
	}
}
