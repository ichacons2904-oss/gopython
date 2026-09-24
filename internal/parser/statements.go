package parser

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{Pos: positionOf(parser.current())}

	for !parser.check(lexer.EOF) {
		if parser.match(lexer.Newline) {
			continue
		}

		statement, err := parser.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, statement)
	}

	return program, nil
}

func (parser *Parser) parseStatement() (ast.Statement, error) {
	switch parser.current().Type {
	case lexer.If:
		return parser.parseIfStatement()
	case lexer.While:
		return parser.parseWhileStatement()
	case lexer.For:
		return parser.parseForStatement()
	case lexer.Def:
		if parser.blockDepth > 0 {
			return nil, Error{Token: parser.current(), Message: "nested function definitions are not supported"}
		}
		return parser.parseFunctionDefinition()
	case lexer.Return:
		if parser.functionDepth == 0 {
			return nil, Error{Token: parser.current(), Message: "return outside function"}
		}
		return parser.parseReturnStatement()
	case lexer.Break:
		if parser.loopDepth == 0 {
			return nil, Error{Token: parser.current(), Message: "break outside loop"}
		}
		return parser.parseBreakStatement()
	case lexer.Continue:
		if parser.loopDepth == 0 {
			return nil, Error{Token: parser.current(), Message: "continue outside loop"}
		}
		return parser.parseContinueStatement()
	}

	if parser.check(lexer.Identifier) && parser.peek().Type == lexer.Assign {
		return parser.parseAssignment()
	}

	expression, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}

	return ast.ExpressionStatement{Expression: expression}, nil
}

func (parser *Parser) parseAssignment() (ast.Statement, error) {
	nameToken, err := parser.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}
	name := ast.Identifier{Pos: positionOf(nameToken), Name: nameToken.Lexeme}

	assignmentToken, err := parser.expect(lexer.Assign)
	if err != nil {
		return nil, err
	}

	value, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}

	return ast.Assignment{
		Pos:   positionOf(assignmentToken),
		Name:  name,
		Value: value,
	}, nil
}

func (parser *Parser) parseReturnStatement() (ast.Statement, error) {
	returnToken := parser.advance()
	var value ast.Expression

	if !parser.check(lexer.Newline) {
		var err error
		value, err = parser.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}

	return ast.ReturnStatement{Pos: positionOf(returnToken), Value: value}, nil
}

func (parser *Parser) parseBreakStatement() (ast.Statement, error) {
	breakToken := parser.advance()
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}
	return ast.BreakStatement{Pos: positionOf(breakToken)}, nil
}

func (parser *Parser) parseContinueStatement() (ast.Statement, error) {
	continueToken := parser.advance()
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}
	return ast.ContinueStatement{Pos: positionOf(continueToken)}, nil
}
