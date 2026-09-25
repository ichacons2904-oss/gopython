package eval

import (
	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/object"
)

func isTruthy(value object.Value) bool {
	switch value := value.(type) {
	case object.Boolean:
		return value.Value
	case object.Integer:
		return value.Value != 0
	case object.String:
		return value.Value != ""
	case object.None:
		return false
	default:
		return true
	}
}

func evaluateLogical(expression *ast.BinaryExpression, environment *object.Environment) (object.Value, error) {
	left, err := evaluateExpression(expression.Left, environment)
	if err != nil {
		return nil, err
	}

	switch expression.Operator {
	case lexer.And:
		if !isTruthy(left) {
			return left, nil
		}
	case lexer.Or:
		if isTruthy(left) {
			return left, nil
		}
	}

	return evaluateExpression(expression.Right, environment)
}
