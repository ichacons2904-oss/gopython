package eval

import (
	"fmt"

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

func evaluateNonlocalStatement(statement *ast.NonlocalStatement, environment *object.Environment) (completion, error) {
	for _, name := range statement.Names {
		if !environment.DeclareNonlocal(name.Name) {
			return completion{}, Error{
				Kind:     SyntaxError,
				Message:  fmt.Sprintf("no binding for nonlocal %q found", name.Name),
				Position: name.Pos,
			}
		}
	}
	return completion{value: object.None{}}, nil
}
