package eval

import (
	"fmt"
	"math"

	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/object"
)

func evaluateExpression(expression ast.Expression, environment *object.Environment) (object.Value, error) {
	switch expression := expression.(type) {
	case *ast.IntegerLiteral:
		return object.Integer{Value: expression.Value}, nil
	case *ast.FloatLiteral:
		return object.Float{Value: expression.Value}, nil
	case *ast.StringLiteral:
		return object.String{Value: expression.Value}, nil
	case *ast.BooleanLiteral:
		return object.Boolean{Value: expression.Value}, nil
	case *ast.NoneLiteral:
		return object.None{}, nil
	case *ast.Identifier:
		return evaluateIdentifier(expression, environment)
	case *ast.UnaryExpression:
		return evaluateUnary(expression, environment)
	case *ast.BinaryExpression:
		return evaluateBinary(expression, environment)
	case *ast.ComparisonExpression:
		return evaluateComparison(expression, environment)
	case *ast.CallExpression:
		return evaluateCall(expression, environment)
	default:
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported expression %T", expression),
			Position: expression.Position(),
		}
	}
}

func evaluateIdentifier(expression *ast.Identifier, environment *object.Environment) (object.Value, error) {
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

func evaluateUnary(expression *ast.UnaryExpression, environment *object.Environment) (object.Value, error) {
	operand, err := evaluateExpression(expression.Operand, environment)
	if err != nil {
		return nil, err
	}

	if expression.Operator == lexer.Not {
		return object.Boolean{Value: !isTruthy(operand)}, nil
	}
	if expression.Operator != lexer.Plus && expression.Operator != lexer.Minus {
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported unary operator %s", expression.Operator),
			Position: expression.Position(),
		}
	}

	switch operand := operand.(type) {
	case object.Integer:
		if expression.Operator == lexer.Minus {
			return object.Integer{Value: -operand.Value}, nil
		}
		return operand, nil
	case object.Float:
		if expression.Operator == lexer.Minus {
			return object.Float{Value: -operand.Value}, nil
		}
		return operand, nil
	default:
		return nil, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("unary operator %s requires a number, got %s", expression.Operator, operand.Type()),
			Position: expression.Position(),
		}
	}
}

func evaluateBinary(expression *ast.BinaryExpression, environment *object.Environment) (object.Value, error) {
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

	if leftInteger, ok := left.(object.Integer); ok {
		if rightInteger, ok := right.(object.Integer); ok {
			return evaluateIntegerArithmetic(expression, leftInteger.Value, rightInteger.Value)
		}
	}

	leftNumber, leftIsNumber := toFloat(left)
	rightNumber, rightIsNumber := toFloat(right)
	if !leftIsNumber || !rightIsNumber {
		return nil, Error{
			Kind: TypeError,
			Message: fmt.Sprintf(
				"operator %s requires numbers, got %s and %s",
				expression.Operator,
				left.Type(),
				right.Type(),
			),
			Position: expression.Position(),
		}
	}
	return evaluateFloatArithmetic(expression, leftNumber, rightNumber)
}

func evaluateIntegerArithmetic(expression *ast.BinaryExpression, left, right int64) (object.Value, error) {
	switch expression.Operator {
	case lexer.Plus:
		return object.Integer{Value: left + right}, nil
	case lexer.Minus:
		return object.Integer{Value: left - right}, nil
	case lexer.Asterisk:
		return object.Integer{Value: left * right}, nil
	case lexer.Slash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Integer{Value: left / right}, nil
	case lexer.DoubleSlash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Integer{Value: floorDivide(left, right)}, nil
	case lexer.Percent:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Integer{Value: floorModulo(left, right)}, nil
	default:
		return nil, unsupportedBinaryOperator(expression)
	}
}

func evaluateFloatArithmetic(expression *ast.BinaryExpression, left, right float64) (object.Value, error) {
	switch expression.Operator {
	case lexer.Plus:
		return object.Float{Value: left + right}, nil
	case lexer.Minus:
		return object.Float{Value: left - right}, nil
	case lexer.Asterisk:
		return object.Float{Value: left * right}, nil
	case lexer.Slash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: left / right}, nil
	case lexer.DoubleSlash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: floorDivideFloat(left, right)}, nil
	case lexer.Percent:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: floorModuloFloat(left, right)}, nil
	default:
		return nil, unsupportedBinaryOperator(expression)
	}
}

func toFloat(value object.Value) (float64, bool) {
	switch value := value.(type) {
	case object.Integer:
		return float64(value.Value), true
	case object.Float:
		return value.Value, true
	default:
		return 0, false
	}
}

func divisionByZero(expression *ast.BinaryExpression) error {
	return Error{
		Kind:     ZeroDivisionError,
		Message:  "division by zero",
		Position: expression.Position(),
	}
}

func unsupportedBinaryOperator(expression *ast.BinaryExpression) error {
	return Error{
		Kind:     RuntimeError,
		Message:  fmt.Sprintf("unsupported binary operator %s", expression.Operator),
		Position: expression.Position(),
	}
}

func floorModulo(dividend, divisor int64) int64 {
	remainder := dividend % divisor
	if remainder != 0 && (remainder < 0) != (divisor < 0) {
		remainder += divisor
	}
	return remainder
}

func floorModuloFloat(dividend, divisor float64) float64 {
	remainder := math.Mod(dividend, divisor)
	if remainder == 0 {
		return math.Copysign(0, divisor)
	}
	if (remainder < 0) != (divisor < 0) {
		remainder += divisor
	}
	return remainder
}

func floorDivide(dividend, divisor int64) int64 {
	quotient := dividend / divisor
	if dividend%divisor != 0 && (dividend < 0) != (divisor < 0) {
		quotient--
	}
	return quotient
}

func floorDivideFloat(dividend, divisor float64) float64 {
	remainder := math.Mod(dividend, divisor)
	quotient := (dividend - remainder) / divisor
	if remainder != 0 && (remainder < 0) != (divisor < 0) {
		quotient--
	}
	if quotient == 0 {
		return math.Copysign(0, dividend/divisor)
	}
	floored := math.Floor(quotient)
	if quotient-floored > 0.5 {
		floored++
	}
	return floored
}
