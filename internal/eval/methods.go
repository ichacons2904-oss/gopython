package eval

import (
	"fmt"

	"gopython/internal/object"
	"gopython/internal/source"
)

type builtinMethodFunction func(object.Value, []object.Value, source.Position) (object.Value, error)
type listMethodFunction func(*object.List, []object.Value, source.Position) (object.Value, error)

type builtinMethod struct {
	minArgs int
	maxArgs int
	call    builtinMethodFunction
}

var builtinMethodsByType = map[object.Type]map[string]builtinMethod{
	object.ListType: {
		"append": {minArgs: 1, maxArgs: 1, call: listMethod(appendList)},
		"extend": {minArgs: 1, maxArgs: 1, call: listMethod(extendList)},
		"pop":    {minArgs: 0, maxArgs: 1, call: listMethod(popList)},
		"clear":  {minArgs: 0, maxArgs: 0, call: listMethod(clearList)},
		"copy":   {minArgs: 0, maxArgs: 0, call: listMethod(copyList)},
	},
}

func lookupBuiltinMethod(receiver object.Value, name string) (builtinMethod, bool) {
	method, ok := builtinMethodsByType[receiver.Type()][name]
	return method, ok
}

func listMethod(call listMethodFunction) builtinMethodFunction {
	return func(receiver object.Value, arguments []object.Value, position source.Position) (object.Value, error) {
		list, ok := receiver.(*object.List)
		if !ok {
			return nil, Error{Kind: RuntimeError, Message: "unsupported method receiver", Position: position}
		}
		return call(list, arguments, position)
	}
}

func invokeBoundMethod(method object.BoundMethod, arguments []object.Value, position source.Position) (object.Value, error) {
	definition, ok := lookupBuiltinMethod(method.Receiver, method.Name)
	if !ok {
		return nil, Error{
			Kind:     AttributeError,
			Message:  fmt.Sprintf("%s object has no attribute %q", method.Receiver.Type(), method.Name),
			Position: position,
		}
	}
	if err := requireMethodArguments(method.Name, arguments, definition.minArgs, definition.maxArgs, position); err != nil {
		return nil, err
	}
	return definition.call(method.Receiver, arguments, position)
}

func appendList(list *object.List, arguments []object.Value, _ source.Position) (object.Value, error) {
	list.Elements = append(list.Elements, arguments[0])
	return object.None{}, nil
}

func extendList(list *object.List, arguments []object.Value, position source.Position) (object.Value, error) {
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
}

func popList(list *object.List, arguments []object.Value, position source.Position) (object.Value, error) {
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
}

func clearList(list *object.List, _ []object.Value, _ source.Position) (object.Value, error) {
	list.Elements = nil
	return object.None{}, nil
}

func copyList(list *object.List, _ []object.Value, _ source.Position) (object.Value, error) {
	return object.NewList(append([]object.Value{}, list.Elements...)), nil
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
