package eval

import (
	"fmt"

	"gopython/internal/source"
)

type ErrorKind string

const (
	NameError         ErrorKind = "NameError"
	TypeError         ErrorKind = "TypeError"
	ZeroDivisionError ErrorKind = "ZeroDivisionError"
	RuntimeError      ErrorKind = "RuntimeError"
)

type Error struct {
	Kind     ErrorKind
	Message  string
	Position source.Position
}

func (err Error) Error() string {
	return fmt.Sprintf(
		"%s at %d:%d: %s",
		err.Kind,
		err.Position.Line,
		err.Position.Column,
		err.Message,
	)
}
