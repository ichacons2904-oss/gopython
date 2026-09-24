package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Error struct {
	Message string
	Line    int
	Column  int
}

func (e Error) Error() string {
	return fmt.Sprintf("lexical error at %d:%d: %s", e.Line, e.Column, e.Message)
}

type Lexer struct {
	chars  []rune
	index  int
	line   int
	column int

	atLineStart bool
	indentStack []int
	tokens      []Token
	lexed       bool
}

func New(source string) *Lexer {
	source = strings.ReplaceAll(source, "\r\n", "\n")
	return &Lexer{
		chars:       []rune(source),
		line:        1,
		column:      1,
		atLineStart: true,
		indentStack: []int{0},
	}
}

func (lexer *Lexer) Lex() ([]Token, error) {
	if lexer.lexed {
		return nil, fmt.Errorf("lexer can only be used once")
	}
	lexer.lexed = true

	for lexer.index < len(lexer.chars) {
		if lexer.atLineStart {
			if err := lexer.scanLineStart(); err != nil {
				return nil, err
			}
			if lexer.index >= len(lexer.chars) {
				break
			}
			if lexer.atLineStart {
				continue
			}
		}

		if err := lexer.scanToken(); err != nil {
			return nil, err
		}
	}

	// A final newline is optional in source files. Add one so the parser sees
	// the same statement boundary in both cases.
	if !lexer.atLineStart && len(lexer.tokens) > 0 {
		lexer.emit(Newline, "\n", lexer.line, lexer.column)
		lexer.line++
		lexer.column = 1
		lexer.atLineStart = true
	}

	for len(lexer.indentStack) > 1 {
		lexer.indentStack = lexer.indentStack[:len(lexer.indentStack)-1]
		lexer.emit(Dedent, "", lexer.line, lexer.column)
	}
	lexer.emit(EOF, "", lexer.line, lexer.column)

	return lexer.tokens, nil
}

func (lexer *Lexer) scanLineStart() error {
	indent := 0
	indentStart := lexer.index

	for lexer.index < len(lexer.chars) && lexer.chars[lexer.index] == ' ' {
		indent++
		lexer.advance()
	}

	if lexer.index < len(lexer.chars) && lexer.chars[lexer.index] == '\t' {
		return lexer.errorAt("tabs are not supported for indentation")
	}

	// Blank and comment-only lines do not affect indentation or produce a
	// NEWLINE token.
	if lexer.index >= len(lexer.chars) {
		return nil
	}
	if lexer.chars[lexer.index] == '\n' {
		lexer.advance()
		return nil
	}
	if lexer.chars[lexer.index] == '#' {
		lexer.skipComment()
		if lexer.index < len(lexer.chars) && lexer.chars[lexer.index] == '\n' {
			lexer.advance()
		}
		return nil
	}

	currentIndent := lexer.indentStack[len(lexer.indentStack)-1]
	if indent > currentIndent {
		lexer.indentStack = append(lexer.indentStack, indent)
		lexer.emit(Indent, string(lexer.chars[indentStart:lexer.index]), lexer.line, 1)
	} else if indent < currentIndent {
		for len(lexer.indentStack) > 1 && indent < lexer.indentStack[len(lexer.indentStack)-1] {
			lexer.indentStack = lexer.indentStack[:len(lexer.indentStack)-1]
			lexer.emit(Dedent, "", lexer.line, 1)
		}
		if indent != lexer.indentStack[len(lexer.indentStack)-1] {
			return lexer.errorAt("inconsistent indentation")
		}
	}

	lexer.atLineStart = false
	return nil
}

func (lexer *Lexer) scanToken() error {
	startLine, startColumn := lexer.line, lexer.column
	ch := lexer.chars[lexer.index]

	if ch == ' ' {
		lexer.advance()
		return nil
	}
	if ch == '\t' {
		return lexer.errorAt("tabs are not supported")
	}
	if ch == '\n' {
		lexer.advance()
		lexer.emit(Newline, "\n", startLine, startColumn)
		lexer.atLineStart = true
		return nil
	}
	if ch == '#' {
		lexer.skipComment()
		return nil
	}

	if isIdentifierStart(ch) {
		start := lexer.index
		for lexer.index < len(lexer.chars) && isIdentifierPart(lexer.chars[lexer.index]) {
			lexer.advance()
		}
		lexeme := string(lexer.chars[start:lexer.index])
		lexer.emit(tokenTypeForIdentifier(lexeme), lexeme, startLine, startColumn)
		return nil
	}

	if isASCIIDigit(ch) {
		start := lexer.index
		for lexer.index < len(lexer.chars) && isASCIIDigit(lexer.chars[lexer.index]) {
			lexer.advance()
		}
		lexer.emit(Integer, string(lexer.chars[start:lexer.index]), startLine, startColumn)
		return nil
	}

	if ch == '\'' || ch == '"' {
		return lexer.scanString(ch, startLine, startColumn)
	}

	if lexer.index+1 < len(lexer.chars) {
		lexeme := string(lexer.chars[lexer.index : lexer.index+2])
		if tokenType, ok := twoCharacterTokens[lexeme]; ok {
			lexer.advance()
			lexer.advance()
			lexer.emit(tokenType, lexeme, startLine, startColumn)
			return nil
		}
	}

	if tokenType, ok := oneCharacterTokens[ch]; ok {
		lexer.advance()
		lexer.emit(tokenType, string(ch), startLine, startColumn)
		return nil
	}

	return lexer.errorAt(fmt.Sprintf("unexpected character %q", ch))
}

func (lexer *Lexer) scanString(quote rune, startLine, startColumn int) error {
	start := lexer.index
	lexer.advance()

	for lexer.index < len(lexer.chars) {
		ch := lexer.chars[lexer.index]
		if ch == '\n' {
			return lexer.errorAt("unterminated string literal")
		}
		if ch == quote {
			lexer.advance()
			lexer.emit(String, string(lexer.chars[start:lexer.index]), startLine, startColumn)
			return nil
		}
		if ch == '\\' {
			lexer.advance()
			if lexer.index >= len(lexer.chars) {
				return lexer.errorAt("unterminated string literal")
			}
			if !isSupportedEscape(lexer.chars[lexer.index]) {
				return lexer.errorAt(fmt.Sprintf("unsupported escape sequence \\%c", lexer.chars[lexer.index]))
			}
		}
		lexer.advance()
	}

	return lexer.errorAt("unterminated string literal")
}

func (lexer *Lexer) emit(tokenType TokenType, lexeme string, line, column int) {
	lexer.tokens = append(lexer.tokens, Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   line,
		Column: column,
	})
}

func (lexer *Lexer) advance() {
	if lexer.index >= len(lexer.chars) {
		return
	}
	if lexer.chars[lexer.index] == '\n' {
		lexer.line++
		lexer.column = 1
	} else {
		lexer.column++
	}
	lexer.index++
}

func (lexer *Lexer) skipComment() {
	for lexer.index < len(lexer.chars) && lexer.chars[lexer.index] != '\n' {
		lexer.advance()
	}
}

func (lexer *Lexer) errorAt(message string) error {
	return Error{Message: message, Line: lexer.line, Column: lexer.column}
}

func isIdentifierStart(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

func isIdentifierPart(ch rune) bool {
	return isIdentifierStart(ch) || unicode.IsDigit(ch)
}

func isASCIIDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isSupportedEscape(ch rune) bool {
	switch ch {
	case '\\', '\'', '"', 'n', 't', 'r':
		return true
	default:
		return false
	}
}
