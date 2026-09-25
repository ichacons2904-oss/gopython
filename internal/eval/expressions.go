package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/object"
)

func evaluateExpression(expression ast.Expression, environment *object.Environment) (object.Value, error) {
	switch expression := expression.(type) {
	case ast.IntegerLiteral:
		return object.Integer{Value: expression.Value}, nil
	case ast.StringLiteral:
		return object.String{Value: expression.Value}, nil
	case ast.BooleanLiteral:
		return object.Boolean{Value: expression.Value}, nil
	case ast.NoneLiteral:
		return object.None{}, nil
	case ast.Identifier:
		return evaluateIdentifier(expression, environment)
	case ast.UnaryExpression:
		return evaluateUnary(expression, environment)
	case ast.BinaryExpression:
		return evaluateBinary(expression, environment)
	case ast.ComparisonExpression:
		return evaluateComparison(expression, environment)
	case ast.CallExpression:
		return evaluateCall(expression, environment)
	default:
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported expression %T", expression),
			Position: expression.Position(),
		}
	}
}

func evaluateIdentifier(expression ast.Identifier, environment *object.Environment) (object.Value, error) {
	value, ok := environment.Get(expression.Name)
	if !ok {
		return nil, Error{
			Kind:     NameError,
			Message:  fmt.Sprintf("name %q is not defined", expression.Name),
			Position: expression.Position(),
		}
	}
	return value, nil
}

func evaluateUnary(expression ast.UnaryExpression, environment *object.Environment) (object.Value, error) {
	operand, err := evaluateExpression(expression.Operand, environment)
	if err != nil {
		return nil, err
	}

	if expression.Operator == lexer.Not {
		return object.Boolean{Value: !isTruthy(operand)}, nil
	}

	integer, ok := operand.(object.Integer)
	if !ok {
		return nil, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("unary operator %s requires an integer, got %s", expression.Operator, operand.Type()),
			Position: expression.Position(),
		}
	}

	switch expression.Operator {
	case lexer.Plus:
		return integer, nil
	case lexer.Minus:
		return object.Integer{Value: -integer.Value}, nil
	default:
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported unary operator %s", expression.Operator),
			Position: expression.Position(),
		}
	}
}

func evaluateBinary(expression ast.BinaryExpression, environment *object.Environment) (object.Value, error) {
	if expression.Operator == lexer.And || expression.Operator == lexer.Or {
		return evaluateLogical(expression, environment)
	}

	left, err := evaluateExpression(expression.Left, environment)
	if err != nil {
		return nil, err
	}

	right, err := evaluateExpression(expression.Right, environment)
	if err != nil {
		return nil, err
	}

	if expression.Operator == lexer.Plus {
		if leftString, ok := left.(object.String); ok {
			if rightString, ok := right.(object.String); ok {
				return object.String{Value: leftString.Value + rightString.Value}, nil
			}
		}
	}

	leftInteger, leftIsInteger := left.(object.Integer)
	rightInteger, rightIsInteger := right.(object.Integer)
	if !leftIsInteger || !rightIsInteger {
		return nil, Error{
			Kind: TypeError,
			Message: fmt.Sprintf(
				"operator %s requires integers, got %s and %s",
				expression.Operator,
				left.Type(),
				right.Type(),
			),
			Position: expression.Position(),
		}
	}

	switch expression.Operator {
	case lexer.Plus:
		return object.Integer{Value: leftInteger.Value + rightInteger.Value}, nil
	case lexer.Minus:
		return object.Integer{Value: leftInteger.Value - rightInteger.Value}, nil
	case lexer.Asterisk:
		return object.Integer{Value: leftInteger.Value * rightInteger.Value}, nil
	case lexer.Slash:
		if rightInteger.Value == 0 {
			return nil, Error{
				Kind:     ZeroDivisionError,
				Message:  "division by zero",
				Position: expression.Position(),
			}
		}
		return object.Integer{Value: leftInteger.Value / rightInteger.Value}, nil
	case lexer.Percent:
		if rightInteger.Value == 0 {
			return nil, Error{
				Kind:     ZeroDivisionError,
				Message:  "integer modulo by zero",
				Position: expression.Position(),
			}
		}
		return object.Integer{Value: leftInteger.Value % rightInteger.Value}, nil
	default:
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported binary operator %s", expression.Operator),
			Position: expression.Position(),
		}
	}
}
