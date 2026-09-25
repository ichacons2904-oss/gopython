package eval

import (
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
		default:
			result = body
		}
	}
}
