package builtin

import (
	"fmt"
	"gopython/internal/object"
)

func tupleFunction(arguments []object.Value) (object.Value, error) {
	if len(arguments) > 1 {
		return nil, object.BuiltinTypeError{Message: fmt.Sprintf("tuple expected at most 1 argument, got %d", len(arguments))}
	}
	if len(arguments) == 0 {
		return object.NewTuple(nil), nil
	}
	iterable, ok := arguments[0].(object.Iterable)
	if !ok {
		return nil, object.BuiltinTypeError{Message: fmt.Sprintf("%s object is not iterable", arguments[0].Type())}
	}
	var elements []object.Value
	iterator := iterable.Iterator()
	for {
		value, ok := iterator.Next()
		if !ok {
			break
		}
		elements = append(elements, value)
	}
	return object.NewTuple(elements), nil
}
