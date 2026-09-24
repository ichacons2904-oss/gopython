package lexer

import (
	"reflect"
	"strings"
	"testing"
)

func TestIfBlock(t *testing.T) {
	source := "if x > 0:\n    print(x)\n"

	tokens, err := New(source).Lex()
	if err != nil {
		t.Fatalf("Lex() returned error: %v", err)
	}

	want := []Token{
		{Type: If, Lexeme: "if", Line: 1, Column: 1},
		{Type: Identifier, Lexeme: "x", Line: 1, Column: 4},
		{Type: GreaterThan, Lexeme: ">", Line: 1, Column: 6},
		{Type: Integer, Lexeme: "0", Line: 1, Column: 8},
		{Type: Colon, Lexeme: ":", Line: 1, Column: 9},
		{Type: Newline, Lexeme: "\n", Line: 1, Column: 10},
		{Type: Indent, Lexeme: "    ", Line: 2, Column: 1},
		{Type: Identifier, Lexeme: "print", Line: 2, Column: 5},
		{Type: LeftParenthesis, Lexeme: "(", Line: 2, Column: 10},
		{Type: Identifier, Lexeme: "x", Line: 2, Column: 11},
		{Type: RightParenthesis, Lexeme: ")", Line: 2, Column: 12},
		{Type: Newline, Lexeme: "\n", Line: 2, Column: 13},
		{Type: Dedent, Lexeme: "", Line: 3, Column: 1},
		{Type: EOF, Lexeme: "", Line: 3, Column: 1},
	}

	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("tokens mismatch\n got: %#v\nwant: %#v", tokens, want)
	}
}

func TestLexAddsFinalNewline(t *testing.T) {
	tokens, err := New("x = 1").Lex()
	if err != nil {
		t.Fatalf("Lex() returned error: %v", err)
	}

	if got := tokens[len(tokens)-2].Type; got != Newline {
		t.Fatalf("last statement token = %s, want %s", got, Newline)
	}
}

func TestLexClosesIndentedBlockWithoutTrailingNewline(t *testing.T) {
	tokens, err := New("if x:\n    print(x)").Lex()
	if err != nil {
		t.Fatalf("Lex() returned error: %v", err)
	}

	if got := tokens[len(tokens)-2]; got != (Token{Type: Dedent, Line: 3, Column: 1}) {
		t.Fatalf("dedent = %#v, want line 3, column 1", got)
	}
}

func TestLexRejectsInconsistentIndentation(t *testing.T) {
	_, err := New("if x:\n    print(x)\n  print(x)\n").Lex()
	if err == nil {
		t.Fatal("Lex() returned nil error for inconsistent indentation")
	}

	lexError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want lexer.Error", err)
	}
	if lexError.Line != 3 || lexError.Column != 3 {
		t.Fatalf("error position = %d:%d, want 3:3", lexError.Line, lexError.Column)
	}
}

func TestLexCanOnlyBeCalledOnce(t *testing.T) {
	lexer := New("x = 1")
	if _, err := lexer.Lex(); err != nil {
		t.Fatalf("first Lex() returned error: %v", err)
	}

	if _, err := lexer.Lex(); err == nil {
		t.Fatal("second Lex() returned nil error")
	}
}

func TestLexSupportsStringsAndEscapes(t *testing.T) {
	tokens, err := New(`print("hello\nworld")`).Lex()
	if err != nil {
		t.Fatalf("Lex() returned error: %v", err)
	}

	if got := tokens[2]; got != (Token{Type: String, Lexeme: `"hello\nworld"`, Line: 1, Column: 7}) {
		t.Fatalf("string token = %#v", got)
	}
}

func TestLexRecognizesKeywordsAndOperators(t *testing.T) {
	source := "if else elif while for in def return break continue and or not True False None = + - * / % == != < <= > >= : ( ) ,\n"
	tokens, err := New(source).Lex()
	if err != nil {
		t.Fatalf("Lex() returned error: %v", err)
	}

	want := []TokenType{
		If, Else, Elif, While, For, In, Def, Return, Break, Continue,
		And, Or, Not, True, False, None,
		Assign, Plus, Minus, Asterisk, Slash, Percent,
		Equal, NotEqual, LessThan, LessEqual, GreaterThan, GreaterEqual,
		Colon, LeftParenthesis, RightParenthesis, Comma, Newline, EOF,
	}

	if got := tokenTypes(tokens); !reflect.DeepEqual(got, want) {
		t.Fatalf("token types = %v, want %v", got, want)
	}
}

func TestLexHandlesCommentsBlankLinesAndNestedIndentation(t *testing.T) {
	source := "if outer:\n    if inner:\n        print(\"x\") # comment\n    print(\"y\")\n\nprint(\"z\")\n"
	tokens, err := New(source).Lex()
	if err != nil {
		t.Fatalf("Lex() returned error: %v", err)
	}

	want := []TokenType{
		If, Identifier, Colon, Newline,
		Indent,
		If, Identifier, Colon, Newline,
		Indent,
		Identifier, LeftParenthesis, String, RightParenthesis, Newline,
		Dedent,
		Identifier, LeftParenthesis, String, RightParenthesis, Newline,
		Dedent,
		Identifier, LeftParenthesis, String, RightParenthesis, Newline,
		EOF,
	}

	if got := tokenTypes(tokens); !reflect.DeepEqual(got, want) {
		t.Fatalf("token types = %v, want %v", got, want)
	}
}

func TestLexReportsErrors(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{name: "invalid character", source: "x @ 1", want: "unexpected character"},
		{name: "tab indentation", source: "if x:\n\tprint(x)", want: "tabs are not supported"},
		{name: "unterminated string", source: `"hello`, want: "unterminated string literal"},
		{name: "unsupported escape", source: `"hello\q"`, want: "unsupported escape sequence"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.source).Lex()
			if err == nil {
				t.Fatal("Lex() returned nil error")
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %q, want it to contain %q", err, test.want)
			}
		})
	}
}

func tokenTypes(tokens []Token) []TokenType {
	types := make([]TokenType, len(tokens))
	for index, token := range tokens {
		types[index] = token.Type
	}
	return types
}
