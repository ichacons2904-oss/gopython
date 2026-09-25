package eval

import "gopython/internal/object"

type Flow int

const (
	NoFlow Flow = iota
	BreakFlow
	ContinueFlow
	ReturnFlow
)

type Result struct {
	Value object.Value
	Flow  Flow
}

func (flow Flow) String() string {
	switch flow {
	case BreakFlow:
		return "break"
	case ContinueFlow:
		return "continue"
	case ReturnFlow:
		return "return"
	default:
		return "normal execution"
	}
}
