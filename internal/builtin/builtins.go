package builtin

import (
	"io"

	"gopython/internal/object"
)

func Register(environment *object.Environment, output io.Writer) {
	if output == nil {
		output = io.Discard
	}

	environment.Set("print", object.Builtin{
		Name:     "print",
		Function: printFunction(output),
	})
	environment.Set("range", object.Builtin{
		Name:     "range",
		Function: rangeFunction,
	})
}
