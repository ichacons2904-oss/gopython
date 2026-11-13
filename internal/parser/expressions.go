package parser

import (
	"fmt"
	"strconv"
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
		left = &ast.BinaryExpression{Pos: positionOf(operator), Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseAnd() (ast.Expression, error) {
	left, err := parser.parseNot()
	if err != nil {
		return nil, err
	}

	for parser.check(lexer.And) {
		operator := parser.advance()
		right, err := parser.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{Pos: positionOf(operator), Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseNot() (ast.Expression, error) {
	if parser.check(lexer.Not) {
		operator := parser.advance()
		operand, err := parser.parseNot()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpression{Pos: positionOf(operator), Operator: operator.Type, Operand: operand}, nil
	}
	return parser.parseComparison()
}

func (parser *Parser) parseComparison() (ast.Expression, error) {
	left, err := parser.parseAdditive()
	if err != nil {
		return nil, err
	}

	if !parser.isComparisonStart() {
		return left, nil
	}

	comparison := &ast.ComparisonExpression{
		Pos:  left.Position(),
		Left: left,
	}
	for parser.isComparisonStart() {
		operator := parser.advance()
		if operator.Type == lexer.Not {
			parser.advance() // IN, confirmed by isComparisonStart.
			operator.Type = lexer.NotIn
		}
		right, err := parser.parseAdditive()
		if err != nil {
			return nil, err
		}
		comparison.Comparisons = append(comparison.Comparisons, ast.Comparison{
			Pos:      positionOf(operator),
			Operator: operator.Type,
			Right:    right,
		})
	}
	return comparison, nil
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
		left = &ast.BinaryExpression{Pos: positionOf(operator), Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseMultiplicative() (ast.Expression, error) {
	left, err := parser.parseUnary()
	if err != nil {
		return nil, err
	}

	for parser.check(lexer.Asterisk) || parser.check(lexer.Slash) || parser.check(lexer.DoubleSlash) || parser.check(lexer.Percent) {
		operator := parser.advance()
		right, err := parser.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpression{Pos: positionOf(operator), Left: left, Operator: operator.Type, Right: right}
	}
	return left, nil
}

func (parser *Parser) parseUnary() (ast.Expression, error) {
	if parser.check(lexer.Plus) || parser.check(lexer.Minus) {
		operator := parser.advance()
		operand, err := parser.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpression{Pos: positionOf(operator), Operator: operator.Type, Operand: operand}, nil
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
		expression = &ast.IntegerLiteral{Pos: positionOf(token), Value: value}
	case lexer.Float:
		parser.advance()
		value, err := parseFloat(token)
		if err != nil {
			return nil, err
		}
		expression = &ast.FloatLiteral{Pos: positionOf(token), Value: value}
	case lexer.String:
		parser.advance()
		value, err := decodeString(token)
		if err != nil {
			return nil, err
		}
		expression = &ast.StringLiteral{Pos: positionOf(token), Value: value}
	case lexer.True, lexer.False:
		parser.advance()
		expression = &ast.BooleanLiteral{Pos: positionOf(token), Value: token.Type == lexer.True}
	case lexer.None:
		parser.advance()
		expression = &ast.NoneLiteral{Pos: positionOf(token)}
	case lexer.Identifier:
		parser.advance()
		expression = &ast.Identifier{Pos: positionOf(token), Name: token.Lexeme}
	case lexer.LeftBracket:
		var err error
		expression, err = parser.parseListLiteral()
		if err != nil {
			return nil, err
		}
	case lexer.LeftBrace:
		var err error
		expression, err = parser.parseDictLiteral()
		if err != nil {
			return nil, err
		}
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

	for {
		switch parser.current().Type {
		case lexer.LeftParenthesis:
			var err error
			expression, err = parser.parseCall(expression)
			if err != nil {
				return nil, err
			}
		case lexer.LeftBracket:
			open := parser.advance()
			index, err := parser.parseExpression()
			if err != nil {
				return nil, err
			}
			if _, err := parser.expect(lexer.RightBracket); err != nil {
				return nil, err
			}
			expression = &ast.IndexExpression{Pos: positionOf(open), Collection: expression, Index: index}
		case lexer.Dot:
			dot := parser.advance()
			name, err := parser.expect(lexer.Identifier)
			if err != nil {
				return nil, err
			}
			expression = &ast.AttributeExpression{Pos: positionOf(dot), Object: expression, Name: name.Lexeme}
		default:
			return expression, nil
		}
	}
}

func (parser *Parser) parseListLiteral() (ast.Expression, error) {
	open := parser.advance()
	elements := []ast.Expression{}
	for !parser.check(lexer.RightBracket) {
		element, err := parser.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
		if !parser.match(lexer.Comma) {
			break
		}
	}
	if _, err := parser.expect(lexer.RightBracket); err != nil {
		return nil, err
	}
	return &ast.ListLiteral{Pos: positionOf(open), Elements: elements}, nil
}

func (parser *Parser) parseDictLiteral() (ast.Expression, error) {
	open := parser.advance()
	entries := []ast.DictEntry{}
	for !parser.check(lexer.RightBrace) {
		key, err := parser.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := parser.expect(lexer.Colon); err != nil {
			return nil, err
		}
		value, err := parser.parseExpression()
		if err != nil {
			return nil, err
		}
		entries = append(entries, ast.DictEntry{Key: key, Value: value})
		if !parser.match(lexer.Comma) {
			break
		}
	}
	if _, err := parser.expect(lexer.RightBrace); err != nil {
		return nil, err
	}
	return &ast.DictLiteral{Pos: positionOf(open), Entries: entries}, nil
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
	return &ast.CallExpression{Pos: positionOf(openParen), Callee: callee, Arguments: arguments}, nil
}

func isComparisonOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.In, lexer.Equal, lexer.NotEqual, lexer.LessThan, lexer.LessEqual, lexer.GreaterThan, lexer.GreaterEqual:
		return true
	default:
		return false
	}
}

func (parser *Parser) isComparisonStart() bool {
	return isComparisonOperator(parser.current().Type) ||
		(parser.check(lexer.Not) && parser.peek().Type == lexer.In)
}

func parseFloat(token lexer.Token) (float64, error) {
	value, err := strconv.ParseFloat(token.Lexeme, 64)
	if err != nil {
		if numberError, ok := err.(*strconv.NumError); ok && numberError.Err == strconv.ErrRange {
			return value, nil
		}
		return 0, Error{Token: token, Message: "invalid float literal"}
	}
	return value, nil
}

func parseInteger(token lexer.Token) (int64, error) {
	value, err := strconv.ParseInt(token.Lexeme, 10, 64)
	if err != nil {
		if numberError, ok := err.(*strconv.NumError); ok && numberError.Err == strconv.ErrRange {
			return 0, Error{Token: token, Message: "integer literal is too large"}
		}
		return 0, Error{Token: token, Message: "invalid integer literal"}
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

func expressionError(token lexer.Token) error {
	return Error{
		Token:   token,
		Message: fmt.Sprintf("expected expression, got %s", token.Type),
	}
}
