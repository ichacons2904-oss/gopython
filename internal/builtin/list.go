package builtin

import (
	"fmt"

	"gopython/internal/object"
)

func listFunction(arguments []object.Value) (object.Value, error) {
	if len(arguments) > 1 {
		return nil, object.BuiltinTypeError{Message: fmt.Sprintf("list expected at most 1 argument, got %d", len(arguments))}
	}
	elements := []object.Value{}
	if len(arguments) == 0 {
		return object.NewList(elements), nil
	}
	iterable, ok := arguments[0].(object.Iterable)
	if !ok {
		return nil, object.BuiltinTypeError{Message: fmt.Sprintf("%s object is not iterable", arguments[0].Type())}
	}
	iterator := iterable.Iterator()
	for {
		value, ok := iterator.Next()
		if !ok {
			break
		}
		elements = append(elements, value)
	}
	return object.NewList(elements), nil
}
