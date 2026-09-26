package eval

import (
	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateFunctionDefinition(statement *ast.FunctionDefinition, environment *object.Environment) (completion, error) {
	function := object.Function{
		Name:       statement.Name.Name,
		Parameters: statement.Parameters,
		Body:       statement.Body,
		Closure:    environment,
	}
	environment.Set(statement.Name.Name, function)
	return completion{value: object.None{}}, nil
}

func evaluateReturnStatement(statement *ast.ReturnStatement, environment *object.Environment) (completion, error) {
	if statement.Value == nil {
		return completion{value: object.None{}, kind: returnCompletion}, nil
	}

	value, err := evaluateExpression(statement.Value, environment)
	if err != nil {
		return completion{}, err
	}
	return completion{value: value, kind: returnCompletion}, nil
}
