package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateWhileStatement(statement ast.WhileStatement, environment *object.Environment) (Result, error) {
	result := Result{Value: object.None{}}

	for {
		condition, err := evaluateExpression(statement.Condition, environment)
		if err != nil {
			return Result{}, err
		}
		if !isTruthy(condition) {
			return result, nil
		}

		body, err := evaluateStatements(statement.Body, environment)
		if err != nil {
			return Result{}, err
		}

		switch body.Flow {
		case BreakFlow:
			return result, nil
		case ContinueFlow:
			continue
		case ReturnFlow:
			return body, nil
		default:
			result = body
		}
	}
}

func evaluateForStatement(statement ast.ForStatement, environment *object.Environment) (Result, error) {
	iterable, err := evaluateExpression(statement.Iterable, environment)
	if err != nil {
		return Result{}, err
	}

	iterableValue, ok := iterable.(object.Iterable)
	if !ok {
		return Result{}, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("%s object is not iterable", iterable.Type()),
			Position: statement.Iterable.Position(),
		}
	}

	result := Result{Value: object.None{}}
	iterator := iterableValue.Iterator()
	for {
		value, ok := iterator.Next()
		if !ok {
			return result, nil
		}

		environment.Set(statement.Target.Name, value)
		body, err := evaluateStatements(statement.Body, environment)
		if err != nil {
			return Result{}, err
		}

		switch body.Flow {
		case BreakFlow:
			return result, nil
		case ContinueFlow:
			continue
		case ReturnFlow:
			return body, nil
		default:
			result = body
		}
	}
}
