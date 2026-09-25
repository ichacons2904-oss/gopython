package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/object"
	"gopython/internal/source"
)

func evaluateComparison(expression ast.ComparisonExpression, environment *object.Environment) (object.Value, error) {
	if len(expression.Comparisons) == 0 {
		return nil, Error{
			Kind:     RuntimeError,
			Message:  "comparison has no operators",
			Position: expression.Position(),
		}
	}

	left, err := evaluateExpression(expression.Left, environment)
	if err != nil {
		return nil, err
	}

	for _, comparison := range expression.Comparisons {
		right, err := evaluateExpression(comparison.Right, environment)
		if err != nil {
			return nil, err
		}

		result, err := compareValues(left, right, comparison.Operator, comparison.Pos)
		if err != nil {
			return nil, err
		}
		if !result {
			return object.Boolean{Value: false}, nil
		}

		left = right
	}

	return object.Boolean{Value: true}, nil
}

func compareValues(left, right object.Value, operator lexer.TokenType, position source.Position) (bool, error) {
	switch operator {
	case lexer.Equal:
		return valuesEqual(left, right), nil
	case lexer.NotEqual:
		return !valuesEqual(left, right), nil
	}

	if leftInteger, ok := left.(object.Integer); ok {
		if rightInteger, ok := right.(object.Integer); ok {
			return compareIntegers(leftInteger.Value, rightInteger.Value, operator), nil
		}
	}

	if leftString, ok := left.(object.String); ok {
		if rightString, ok := right.(object.String); ok {
			return compareStrings(leftString.Value, rightString.Value, operator), nil
		}
	}

	return false, Error{
		Kind: TypeError,
		Message: fmt.Sprintf(
			"operator %s is not supported between %s and %s",
			operator,
			left.Type(),
			right.Type(),
		),
		Position: position,
	}
}

func valuesEqual(left, right object.Value) bool {
	switch left := left.(type) {
	case object.Integer:
		right, ok := right.(object.Integer)
		return ok && left.Value == right.Value
	case object.String:
		right, ok := right.(object.String)
		return ok && left.Value == right.Value
	case object.Boolean:
		right, ok := right.(object.Boolean)
		return ok && left.Value == right.Value
	case object.None:
		_, ok := right.(object.None)
		return ok
	default:
		return false
	}
}

func compareIntegers(left, right int64, operator lexer.TokenType) bool {
	switch operator {
	case lexer.LessThan:
		return left < right
	case lexer.LessEqual:
		return left <= right
	case lexer.GreaterThan:
		return left > right
	case lexer.GreaterEqual:
		return left >= right
	default:
		return false
	}
}

func compareStrings(left, right string, operator lexer.TokenType) bool {
	switch operator {
	case lexer.LessThan:
		return left < right
	case lexer.LessEqual:
		return left <= right
	case lexer.GreaterThan:
		return left > right
	case lexer.GreaterEqual:
		return left >= right
	default:
		return false
	}
}
