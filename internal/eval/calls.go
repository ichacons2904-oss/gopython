package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
	"gopython/internal/source"
)

func evaluateCall(expression *ast.CallExpression, environment *object.Environment) (object.Value, error) {
	callee, err := evaluateExpression(expression.Callee, environment)
	if err != nil {
		return nil, err
	}

	arguments := make([]object.Value, len(expression.Arguments))
	for index, argument := range expression.Arguments {
		value, err := evaluateExpression(argument, environment)
		if err != nil {
			return nil, err
		}
		arguments[index] = value
	}

	switch callable := callee.(type) {
	case object.Builtin:
		return invokeBuiltin(callable, arguments, expression.Position())
	case object.Function:
		return invokeFunction(callable, arguments, expression.Position())
	default:
		return nil, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("%s object is not callable", callee.Type()),
			Position: expression.Position(),
		}
	}
}

func invokeBuiltin(builtin object.Builtin, arguments []object.Value, position source.Position) (object.Value, error) {
	value, err := builtin.Function(arguments)
	if err != nil {
		return nil, Error{
			Kind:     RuntimeError,
			Message:  err.Error(),
			Position: position,
		}
	}
	return value, nil
}

func invokeFunction(function object.Function, arguments []object.Value, position source.Position) (object.Value, error) {
	return nil, Error{
		Kind:     RuntimeError,
		Message:  fmt.Sprintf("function %q is not executable yet", function.Name),
		Position: position,
	}
}
