package parser

import (
	"fmt"

	"gopython/internal/lexer"
)

type Parser struct {
	tokens        []lexer.Token
	position      int
	loopDepth     int
	functionDepth int
	blockDepth    int
}

type Error struct {
	Message string
	Token   lexer.Token
}

func (e Error) Error() string {
	return fmt.Sprintf("syntax error at %d:%d: %s", e.Token.Line, e.Token.Column, e.Message)
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (parser *Parser) current() lexer.Token {
	if parser.position < len(parser.tokens) {
		return parser.tokens[parser.position]
	}
	if len(parser.tokens) == 0 {
		return lexer.Token{Type: lexer.EOF, Line: 1, Column: 1}
	}
	last := parser.tokens[len(parser.tokens)-1]
	return lexer.Token{Type: lexer.EOF, Line: last.Line, Column: last.Column}
}

func (parser *Parser) peek() lexer.Token {
	if parser.position+1 < len(parser.tokens) {
		return parser.tokens[parser.position+1]
	}
	return parser.current()
}

func (parser *Parser) advance() lexer.Token {
	token := parser.current()
	if token.Type != lexer.EOF {
		parser.position++
	}
	return token
}

func (parser *Parser) check(tokenType lexer.TokenType) bool {
	return parser.current().Type == tokenType
}

func (parser *Parser) match(tokenTypes ...lexer.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser *Parser) expect(tokenType lexer.TokenType) (lexer.Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}

	current := parser.current()
	return current, Error{
		Token: current,
		Message: fmt.Sprintf(
			"expected %s, got %s",
			tokenType,
			current.Type,
		),
	}
}
