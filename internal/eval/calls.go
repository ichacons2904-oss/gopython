package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateCall(expression ast.CallExpression, environment *object.Environment) (object.Value, error) {
	callee, err := evaluateExpression(expression.Callee, environment)
	if err != nil {
		return nil, err
	}

	callable, ok := callee.(object.Callable)
	if !ok {
		return nil, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("%s object is not callable", callee.Type()),
			Position: expression.Position(),
		}
	}

	arguments := make([]object.Value, len(expression.Arguments))
	for index, argument := range expression.Arguments {
		value, err := evaluateExpression(argument, environment)
		if err != nil {
			return nil, err
		}
		arguments[index] = value
	}

	value, err := callable.Call(arguments)
	if err != nil {
		return nil, Error{
			Kind:     RuntimeError,
			Message:  err.Error(),
			Position: expression.Position(),
		}
	}
	return value, nil
}
