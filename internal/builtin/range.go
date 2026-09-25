package builtin

import (
	"fmt"

	"gopython/internal/object"
)

func rangeFunction(arguments []object.Value) (object.Value, error) {
	if len(arguments) < 1 || len(arguments) > 3 {
		return nil, fmt.Errorf("range expected 1 to 3 arguments, got %d", len(arguments))
	}

	values := make([]int64, len(arguments))
	for index, argument := range arguments {
		integer, ok := argument.(object.Integer)
		if !ok {
			return nil, fmt.Errorf("range arguments must be integers, got %s", argument.Type())
		}
		values[index] = integer.Value
	}

	start, stop, step := int64(0), values[0], int64(1)
	if len(values) >= 2 {
		start, stop = values[0], values[1]
	}
	if len(values) == 3 {
		step = values[2]
	}

	rangeValue, err := object.NewRange(start, stop, step)
	if err != nil {
		return nil, err
	}
	return rangeValue, nil
}
