package eval

import (
	"fmt"

	"gopython/internal/object"
	"gopython/internal/source"
)

func invokeBoundMethod(method object.BoundMethod, arguments []object.Value, position source.Position) (object.Value, error) {
	list, ok := method.Receiver.(*object.List)
	if !ok {
		return nil, Error{Kind: RuntimeError, Message: "unsupported method receiver", Position: position}
	}

	switch method.Name {
	case "append":
		if err := requireMethodArguments(method.Name, arguments, 1, 1, position); err != nil {
			return nil, err
		}
		list.Elements = append(list.Elements, arguments[0])
		return object.None{}, nil
	case "extend":
		if err := requireMethodArguments(method.Name, arguments, 1, 1, position); err != nil {
			return nil, err
		}
		iterable, ok := arguments[0].(object.Iterable)
		if !ok {
			return nil, Error{Kind: TypeError, Message: fmt.Sprintf("%s object is not iterable", arguments[0].Type()), Position: position}
		}
		extended := make([]object.Value, 0)
		iterator := iterable.Iterator()
		for {
			value, ok := iterator.Next()
			if !ok {
				break
			}
			extended = append(extended, value)
		}
		maxInt := int(^uint(0) >> 1)
		if len(list.Elements) > maxInt-len(extended) {
			return nil, Error{Kind: MemoryError, Message: "list extension result is too large", Position: position}
		}
		list.Elements = append(list.Elements, extended...)
		return object.None{}, nil
	case "pop":
		if err := requireMethodArguments(method.Name, arguments, 0, 1, position); err != nil {
			return nil, err
		}
		index := int64(-1)
		if len(arguments) == 1 {
			integer, ok := arguments[0].(object.Integer)
			if !ok {
				return nil, Error{Kind: TypeError, Message: fmt.Sprintf("list indices must be integers, not %s", arguments[0].Type()), Position: position}
			}
			index = integer.Value
		}
		value, ok := list.Pop(index)
		if !ok {
			message := "pop index out of range"
			if len(list.Elements) == 0 {
				message = "pop from empty list"
			}
			return nil, Error{Kind: IndexError, Message: message, Position: position}
		}
		return value, nil
	case "clear":
		if err := requireMethodArguments(method.Name, arguments, 0, 0, position); err != nil {
			return nil, err
		}
		list.Elements = nil
		return object.None{}, nil
	case "copy":
		if err := requireMethodArguments(method.Name, arguments, 0, 0, position); err != nil {
			return nil, err
		}
		return object.NewList(append([]object.Value{}, list.Elements...)), nil
	default:
		return nil, Error{Kind: AttributeError, Message: fmt.Sprintf("list object has no method %q", method.Name), Position: position}
	}
}

func requireMethodArguments(name string, arguments []object.Value, minimum, maximum int, position source.Position) error {
	if len(arguments) >= minimum && len(arguments) <= maximum {
		return nil
	}
	message := fmt.Sprintf("%s() takes %d argument(s), got %d", name, minimum, len(arguments))
	if minimum != maximum {
		message = fmt.Sprintf("%s() takes from %d to %d arguments, got %d", name, minimum, maximum, len(arguments))
	}
	return Error{Kind: TypeError, Message: message, Position: position}
}
