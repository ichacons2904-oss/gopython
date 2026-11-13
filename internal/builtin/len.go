package builtin

import (
	"fmt"
	"unicode/utf8"

	"gopython/internal/object"
)

func lenFunction(arguments []object.Value) (object.Value, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("len() takes exactly one argument (%d given)", len(arguments))
	}

	switch argument := arguments[0].(type) {
	case object.String:
		return object.Integer{Value: int64(utf8.RuneCountInString(argument.Value))}, nil
	case object.Range:
		return object.Integer{Value: argument.Len()}, nil
	case *object.List:
		return object.Integer{Value: int64(len(argument.Elements))}, nil
	case *object.Dict:
		return object.Integer{Value: int64(argument.Len())}, nil
	default:
		return nil, fmt.Errorf("%s object has no len()", argument.Type())
	}
}
