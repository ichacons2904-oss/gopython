package builtin

import (
	"fmt"

	"gopython/internal/object"
)

func dictFunction(arguments []object.Value) (object.Value, error) {
	if len(arguments) > 1 {
		return nil, object.BuiltinTypeError{Message: fmt.Sprintf("dict expected at most 1 argument, got %d", len(arguments))}
	}
	if len(arguments) == 0 {
		return object.NewDict(), nil
	}
	source, ok := arguments[0].(*object.Dict)
	if !ok {
		return nil, object.BuiltinTypeError{Message: fmt.Sprintf("dict() requires a dict, got %s", arguments[0].Type())}
	}
	return source.Copy(), nil
}
