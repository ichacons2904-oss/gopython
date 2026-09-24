package parser

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{Token: parser.current()}

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
	nameToken, _ := parser.expect(lexer.Identifier)
	name := ast.Identifier{Token: nameToken, Name: nameToken.Lexeme}

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
		Token: assignmentToken,
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

	return ast.ReturnStatement{Token: returnToken, Value: value}, nil
}

func (parser *Parser) parseBreakStatement() (ast.Statement, error) {
	breakToken := parser.advance()
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}
	return ast.BreakStatement{Token: breakToken}, nil
}

func (parser *Parser) parseContinueStatement() (ast.Statement, error) {
	continueToken := parser.advance()
	if _, err := parser.expect(lexer.Newline); err != nil {
		return nil, err
	}
	return ast.ContinueStatement{Token: continueToken}, nil
}

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
			Token: parameterToken,
			Name:  parameterToken.Lexeme,
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

	parser.functionDepth++
	body, bodyErr := parser.parseBlock()
	parser.functionDepth--
	if bodyErr != nil {
		return nil, bodyErr
	}

	return ast.FunctionDefinition{
		Token: defToken,
		Name: ast.Identifier{
			Token: nameToken,
			Name:  nameToken.Lexeme,
		},
		Parameters: parameters,
		Body:       body,
	}, nil
}

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
			Token:     elifToken,
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
		Token:     ifToken,
		Condition: condition,
		Body:      body,
		Elif:      branches,
		Else:      elseBody,
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
		Token:     whileToken,
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
		Token:    forToken,
		Target:   ast.Identifier{Token: targetToken, Name: targetToken.Lexeme},
		Iterable: iterable,
		Body:     body,
	}, nil
}

func (parser *Parser) parseLoopBlock() ([]ast.Statement, error) {
	parser.loopDepth++
	defer func() {
		parser.loopDepth--
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

func expressionError(token lexer.Token) error {
	return Error{
		Token:   token,
		Message: fmt.Sprintf("expected expression, got %s", token.Type),
	}
}
