package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func Evaluate(node ast.Node, environment *object.Environment) (object.Value, error) {
	if environment == nil {
		environment = object.NewEnvironment(nil)
	}

	switch node := node.(type) {
	case ast.Program:
		return evaluateProgram(node, environment)
	case *ast.Program:
		return evaluateProgram(*node, environment)
	case ast.ExpressionStatement:
		return evaluateExpression(node.Expression, environment)
	case ast.Assignment:
		return evaluateAssignment(node, environment)
	default:
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported statement %T", node),
			Position: node.Position(),
		}
	}
}

func evaluateProgram(program ast.Program, environment *object.Environment) (object.Value, error) {
	var result object.Value = object.None{}
	for _, statement := range program.Statements {
		value, err := Evaluate(statement, environment)
		if err != nil {
			return nil, err
		}
		result = value
	}
	return result, nil
}

func evaluateAssignment(statement ast.Assignment, environment *object.Environment) (object.Value, error) {
	value, err := evaluateExpression(statement.Value, environment)
	if err != nil {
		return nil, err
	}

	environment.Set(statement.Name.Name, value)
	return object.None{}, nil
}
