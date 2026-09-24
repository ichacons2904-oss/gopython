package parser

import (
	"testing"

	"gopython/internal/lexer"
)

func TestCursorConsumesTokens(t *testing.T) {
	tokens, err := lexer.New("if x > 0:\n    print(x)\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	parser := New(tokens)
	if got := parser.current().Type; got != lexer.If {
		t.Fatalf("current token = %s, want %s", got, lexer.If)
	}
	if got := parser.peek().Type; got != lexer.Identifier {
		t.Fatalf("peek token = %s, want %s", got, lexer.Identifier)
	}

	parser.advance()
	if got := parser.current().Lexeme; got != "x" {
		t.Fatalf("current lexeme = %q, want %q", got, "x")
	}
}

func TestMatch(t *testing.T) {
	tokens, err := lexer.New("if x\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	parser := New(tokens)
	if !parser.match(lexer.If) {
		t.Fatal("match(If) = false, want true")
	}
	if parser.match(lexer.Else, lexer.While) {
		t.Fatal("match(Else, While) = true, want false")
	}
	if got := parser.current().Lexeme; got != "x" {
		t.Fatalf("current lexeme = %q, want %q", got, "x")
	}
}

func TestExpect(t *testing.T) {
	tokens, err := lexer.New("if x\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	parser := New(tokens)
	if _, err := parser.expect(lexer.If); err != nil {
		t.Fatalf("expect(If) returned error: %v", err)
	}

	_, err = parser.expect(lexer.Colon)
	if err == nil {
		t.Fatal("expect(Colon) returned nil error")
	}
	parserError, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want parser.Error", err)
	}
	if parserError.Token.Line != 1 || parserError.Token.Column != 4 {
		t.Fatalf("error position = %d:%d, want 1:4", parserError.Token.Line, parserError.Token.Column)
	}
}

func TestCursorStopsAtEOF(t *testing.T) {
	tokens, err := lexer.New("x\n").Lex()
	if err != nil {
		t.Fatal(err)
	}

	parser := New(tokens)
	parser.advance()
	parser.advance()
	parser.advance()

	if got := parser.current().Type; got != lexer.EOF {
		t.Fatalf("current token = %s, want %s", got, lexer.EOF)
	}
	if got := parser.peek().Type; got != lexer.EOF {
		t.Fatalf("peek token = %s, want %s", got, lexer.EOF)
	}
}
