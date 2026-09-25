package eval

import (
	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateIfStatement(statement *ast.IfStatement, environment *object.Environment) (completion, error) {
	condition, err := evaluateExpression(statement.Condition, environment)
	if err != nil {
		return completion{}, err
	}
	if isTruthy(condition) {
		return evaluateStatements(statement.Body, environment)
	}

	for _, branch := range statement.ElifBranches {
		condition, err := evaluateExpression(branch.Condition, environment)
		if err != nil {
			return completion{}, err
		}
		if isTruthy(condition) {
			return evaluateStatements(branch.Body, environment)
		}
	}

	if statement.ElseBody != nil {
		return evaluateStatements(statement.ElseBody, environment)
	}

	return completion{value: object.None{}}, nil
}
