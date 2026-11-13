package parser

import (
	"testing"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func TestParseDictLiterals(t *testing.T) {
	tests := []struct {
		source  string
		entries int
	}{
		{source: "{}\n", entries: 0},
		{source: "{'a': 1}\n", entries: 1},
		{source: "{'a': 1, 2: [3], 'b': {'c': 4},}\n", entries: 3},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			tokens, err := lexer.New(test.source).Lex()
			if err != nil {
				t.Fatal(err)
			}
			program, err := New(tokens).Parse()
			if err != nil {
				t.Fatal(err)
			}
			statement := program.Statements[0].(*ast.ExpressionStatement)
			dict, ok := statement.Expression.(*ast.DictLiteral)
			if !ok {
				t.Fatalf("expression = %T, want *ast.DictLiteral", statement.Expression)
			}
			if len(dict.Entries) != test.entries {
				t.Fatalf("entry count = %d, want %d", len(dict.Entries), test.entries)
			}
		})
	}
}

func TestParseDictIndexAssignment(t *testing.T) {
	tokens, err := lexer.New("ages = {}\nages['ana'] = 30\nages['ana']\n").Lex()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := program.Statements[1].(*ast.IndexAssignment); !ok {
		t.Fatalf("second statement = %T, want *ast.IndexAssignment", program.Statements[1])
	}
}

func TestParseRejectsMalformedDicts(t *testing.T) {
	for _, input := range []string{"{1}\n", "{1: }\n", "{1: 2\n", "{1: 2 3: 4}\n", "{, }\n"} {
		t.Run(input, func(t *testing.T) {
			tokens, err := lexer.New(input).Lex()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := New(tokens).Parse(); err == nil {
				t.Fatal("Parse() accepted a malformed dict literal")
			}
		})
	}
}
