package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateWhileStatement(statement *ast.WhileStatement, environment *object.Environment) (completion, error) {
	result := completion{value: object.None{}}

	for {
		condition, err := evaluateExpression(statement.Condition, environment)
		if err != nil {
			return completion{}, err
		}
		if !isTruthy(condition) {
			return result, nil
		}

		body, err := evaluateStatements(statement.Body, environment)
		if err != nil {
			return completion{}, err
		}

		switch body.kind {
		case breakCompletion:
			return result, nil
		case continueCompletion:
			continue
		case returnCompletion:
			return body, nil
		default:
			result = body
		}
	}
}

func evaluateForStatement(statement *ast.ForStatement, environment *object.Environment) (completion, error) {
	iterable, err := evaluateExpression(statement.Iterable, environment)
	if err != nil {
		return completion{}, err
	}

	iterableValue, ok := iterable.(object.Iterable)
	if !ok {
		return completion{}, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("%s object is not iterable", iterable.Type()),
			Position: statement.Iterable.Position(),
		}
	}

	result := completion{value: object.None{}}
	iterator := iterableValue.Iterator()
	for {
		value, ok := iterator.Next()
		if !ok {
			return result, nil
		}

		environment.Set(statement.Target.Name, value)
		body, err := evaluateStatements(statement.Body, environment)
		if err != nil {
			return completion{}, err
		}

		switch body.kind {
		case breakCompletion:
			return result, nil
		case continueCompletion:
			continue
		case returnCompletion:
			return body, nil
		default:
			result = body
		}
	}
}
