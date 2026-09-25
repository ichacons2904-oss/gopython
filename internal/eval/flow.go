package eval

import "gopython/internal/object"

type completionKind uint8

const (
	normalCompletion completionKind = iota
	breakCompletion
	continueCompletion
	returnCompletion
)

type completion struct {
	value object.Value
	kind  completionKind
}

func (kind completionKind) String() string {
	switch kind {
	case breakCompletion:
		return "break"
	case continueCompletion:
		return "continue"
	case returnCompletion:
		return "return"
	default:
		return "normal execution"
	}
}
