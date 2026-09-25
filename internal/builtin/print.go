package builtin

import (
	"fmt"
	"io"
	"strings"

	"gopython/internal/object"
)

func printFunction(output io.Writer) object.BuiltinFunction {
	return func(arguments []object.Value) (object.Value, error) {
		values := make([]string, len(arguments))
		for index, argument := range arguments {
			values[index] = argument.Display()
		}

		if _, err := fmt.Fprintln(output, strings.Join(values, " ")); err != nil {
			return nil, err
		}
		return object.None{}, nil
	}
}
