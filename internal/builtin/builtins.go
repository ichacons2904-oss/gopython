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
	environment.Set("dict", object.Builtin{
		Name:     "dict",
		Function: dictFunction,
	})
	environment.Set("len", object.Builtin{
		Name:     "len",
		Function: lenFunction,
	})
	environment.Set("list", object.Builtin{
		Name:     "list",
		Function: listFunction,
	})
	environment.Set("range", object.Builtin{
		Name:     "range",
		Function: rangeFunction,
	})
}
