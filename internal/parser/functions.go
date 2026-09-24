package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) parseFunctionDefinition() (ast.Statement, error) {
	defToken := parser.advance()
	nameToken, err := parser.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.LeftParenthesis); err != nil {
		return nil, err
	}

	parameters := []ast.Identifier{}
	for !parser.check(lexer.RightParenthesis) {
		parameterToken, err := parser.expect(lexer.Identifier)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, ast.Identifier{
			Pos:  positionOf(parameterToken),
			Name: parameterToken.Lexeme,
		})

		if !parser.match(lexer.Comma) {
			break
		}
	}

	if _, err := parser.expect(lexer.RightParenthesis); err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Colon); err != nil {
		return nil, err
	}

	body, bodyErr := parser.parseFunctionBlock()
	if bodyErr != nil {
		return nil, bodyErr
	}

	return ast.FunctionDefinition{
		Pos: positionOf(defToken),
		Name: ast.Identifier{
			Pos:  positionOf(nameToken),
			Name: nameToken.Lexeme,
		},
		Parameters: parameters,
		Body:       body,
	}, nil
}
