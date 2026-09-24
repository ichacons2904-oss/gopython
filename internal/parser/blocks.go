package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) parseLoopBlock() ([]ast.Statement, error) {
	parser.loopDepth++
	defer func() {
		parser.loopDepth--
	}()
	return parser.parseBlock()
}

func (parser *Parser) parseFunctionBlock() ([]ast.Statement, error) {
	previousLoopDepth := parser.loopDepth
	parser.loopDepth = 0
	parser.functionDepth++
	defer func() {
		parser.functionDepth--
		parser.loopDepth = previousLoopDepth
	}()
	return parser.parseBlock()
}

func (parser *Parser) parseBlock() ([]ast.Statement, error) {
	parser.blockDepth++
	defer func() {
		parser.blockDepth--
	}()

	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Indent); err != nil {
		return nil, err
	}

	statements := []ast.Statement{}
	for !parser.check(lexer.Dedent) && !parser.check(lexer.EOF) {
		if parser.match(lexer.Newline) {
			continue
		}

		statement, err := parser.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	if len(statements) == 0 {
		return nil, Error{
			Token:   parser.current(),
			Message: "expected statement in block",
		}
	}
	if _, err := parser.expect(lexer.Dedent); err != nil {
		return nil, err
	}
	return statements, nil
}
