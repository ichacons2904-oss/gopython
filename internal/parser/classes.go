package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) parseClassDefinition() (ast.Statement, error) {
	classToken := parser.advance()
	nameToken, err := parser.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	bases := []ast.Identifier{}
	if parser.match(lexer.LeftParenthesis) {
		for !parser.check(lexer.RightParenthesis) {
			base, err := parser.expect(lexer.Identifier)
			if err != nil {
				return nil, err
			}
			bases = append(bases, ast.Identifier{Pos: positionOf(base), Name: base.Lexeme})
			if !parser.match(lexer.Comma) {
				break
			}
		}
		if _, err := parser.expect(lexer.RightParenthesis); err != nil {
			return nil, err
		}
	}
	if _, err := parser.expect(lexer.Colon); err != nil {
		return nil, err
	}
	body, err := parser.parseClassBlock()
	if err != nil {
		return nil, err
	}
	return &ast.ClassDefinition{
		Pos:   positionOf(classToken),
		Name:  ast.Identifier{Pos: positionOf(nameToken), Name: nameToken.Lexeme},
		Bases: bases,
		Body:  body,
	}, nil
}
