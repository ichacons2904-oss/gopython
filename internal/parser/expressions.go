package parser

import (
	"fmt"
	"strings"

	"gopython/internal/ast"
	"gopython/internal/lexer"
)

func (parser *Parser) parseExpression() (ast.Expression, error) {
	return parser.parseOr()
}

func (parser *Parser) parseOr() (ast.Expression, error) {
	left, err := parser.parseAnd()
	if err != nil {
		return nil, err
	}

	for parser.check(lexer.Or) {
		operator := parser.advance()
		right, err := parser.parseAnd()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{Token: operator, Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseAnd() (ast.Expression, error) {
	left, err := parser.parseComparison()
	if err != nil {
		return nil, err
	}

	for parser.check(lexer.And) {
		operator := parser.advance()
		right, err := parser.parseComparison()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{Token: operator, Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseComparison() (ast.Expression, error) {
	left, err := parser.parseAdditive()
	if err != nil {
		return nil, err
	}

	for isComparisonOperator(parser.current().Type) {
		operator := parser.advance()
		right, err := parser.parseAdditive()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{Token: operator, Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseAdditive() (ast.Expression, error) {
	left, err := parser.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for parser.check(lexer.Plus) || parser.check(lexer.Minus) {
		operator := parser.advance()
		right, err := parser.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{Token: operator, Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseMultiplicative() (ast.Expression, error) {
	left, err := parser.parseUnary()
	if err != nil {
		return nil, err
	}

	for parser.check(lexer.Asterisk) || parser.check(lexer.Slash) || parser.check(lexer.Percent) {
		operator := parser.advance()
		right, err := parser.parseUnary()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpression{Token: operator, Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseUnary() (ast.Expression, error) {
	if parser.check(lexer.Not) || parser.check(lexer.Plus) || parser.check(lexer.Minus) {
		operator := parser.advance()
		operand, err := parser.parseUnary()
		if err != nil {
			return nil, err
		}
		return ast.UnaryExpression{Token: operator, Operator: operator.Type, Operand: operand}, nil
	}
	return parser.parsePrimary()
}

func (parser *Parser) parsePrimary() (ast.Expression, error) {
	token := parser.current()
	var expression ast.Expression

	switch token.Type {
	case lexer.Integer:
		parser.advance()
		value, err := parseInteger(token)
		if err != nil {
			return nil, err
		}
		expression = ast.IntegerLiteral{Token: token, Value: value}
	case lexer.String:
		parser.advance()
		value, err := decodeString(token)
		if err != nil {
			return nil, err
		}
		expression = ast.StringLiteral{Token: token, Value: value}
	case lexer.True, lexer.False:
		parser.advance()
		expression = ast.BooleanLiteral{Token: token, Value: token.Type == lexer.True}
	case lexer.None:
		parser.advance()
		expression = ast.NoneLiteral{Token: token}
	case lexer.Identifier:
		parser.advance()
		expression = ast.Identifier{Token: token, Name: token.Lexeme}
	case lexer.LeftParenthesis:
		parser.advance()
		var err error
		expression, err = parser.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := parser.expect(lexer.RightParenthesis); err != nil {
			return nil, err
		}
	default:
		return nil, expressionError(token)
	}

	for parser.check(lexer.LeftParenthesis) {
		var err error
		expression, err = parser.parseCall(expression)
		if err != nil {
			return nil, err
		}
	}
	return expression, nil
}

func (parser *Parser) parseCall(callee ast.Expression) (ast.Expression, error) {
	openParen := parser.advance()
	arguments := []ast.Expression{}

	for !parser.check(lexer.RightParenthesis) {
		argument, err := parser.parseExpression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)

		if !parser.match(lexer.Comma) {
			break
		}
	}

	if _, err := parser.expect(lexer.RightParenthesis); err != nil {
		return nil, err
	}
	return ast.CallExpression{Token: openParen, Callee: callee, Arguments: arguments}, nil
}

func isComparisonOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.Equal, lexer.NotEqual, lexer.LessThan, lexer.LessEqual, lexer.GreaterThan, lexer.GreaterEqual:
		return true
	default:
		return false
	}
}

func parseInteger(token lexer.Token) (int64, error) {
	var value int64
	for _, digit := range token.Lexeme {
		value = value*10 + int64(digit-'0')
		if value < 0 {
			return 0, Error{Token: token, Message: "integer literal is too large"}
		}
	}
	return value, nil
}

func decodeString(token lexer.Token) (string, error) {
	raw := token.Lexeme
	if len(raw) < 2 {
		return "", Error{Token: token, Message: "invalid string literal"}
	}

	var value strings.Builder
	for index := 1; index < len(raw)-1; index++ {
		if raw[index] != '\\' {
			value.WriteByte(raw[index])
			continue
		}

		index++
		switch raw[index] {
		case '\\':
			value.WriteByte('\\')
		case '\'':
			value.WriteByte('\'')
		case '"':
			value.WriteByte('"')
		case 'n':
			value.WriteByte('\n')
		case 't':
			value.WriteByte('\t')
		case 'r':
			value.WriteByte('\r')
		default:
			return "", Error{
				Token:   token,
				Message: fmt.Sprintf("unsupported escape sequence \\%c", raw[index]),
			}
		}
	}
	return value.String(), nil
}
