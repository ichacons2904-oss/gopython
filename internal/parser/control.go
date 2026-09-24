package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) parseIfStatement() (ast.Statement, error) {
	ifToken := parser.advance()
	condition, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Colon); err != nil {
		return nil, err
	}
	body, err := parser.parseBlock()
	if err != nil {
		return nil, err
	}

	branches := []ast.ElifBranch{}
	for parser.check(lexer.Elif) {
		elifToken := parser.advance()
		elifCondition, err := parser.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := parser.expect(lexer.Colon); err != nil {
			return nil, err
		}
		elifBody, err := parser.parseBlock()
		if err != nil {
			return nil, err
		}
		branches = append(branches, ast.ElifBranch{
			Pos:       positionOf(elifToken),
			Condition: elifCondition,
			Body:      elifBody,
		})
	}

	var elseBody []ast.Statement
	if parser.match(lexer.Else) {
		if _, err := parser.expect(lexer.Colon); err != nil {
			return nil, err
		}
		elseBody, err = parser.parseBlock()
		if err != nil {
			return nil, err
		}
	}

	return ast.IfStatement{
		Pos:          positionOf(ifToken),
		Condition:    condition,
		Body:         body,
		ElifBranches: branches,
		ElseBody:     elseBody,
	}, nil
}

func (parser *Parser) parseWhileStatement() (ast.Statement, error) {
	whileToken := parser.advance()
	condition, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Colon); err != nil {
		return nil, err
	}

	body, err := parser.parseLoopBlock()
	if err != nil {
		return nil, err
	}
	return ast.WhileStatement{
		Pos:       positionOf(whileToken),
		Condition: condition,
		Body:      body,
	}, nil
}

func (parser *Parser) parseForStatement() (ast.Statement, error) {
	forToken := parser.advance()
	targetToken, err := parser.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.In); err != nil {
		return nil, err
	}

	iterable, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Colon); err != nil {
		return nil, err
	}

	body, err := parser.parseLoopBlock()
	if err != nil {
		return nil, err
	}
	return ast.ForStatement{
		Pos:      positionOf(forToken),
		Target:   ast.Identifier{Pos: positionOf(targetToken), Name: targetToken.Lexeme},
		Iterable: iterable,
		Body:     body,
	}, nil
}
