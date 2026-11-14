package object

import "gopython/internal/ast"

type Function struct {
	Name       string
	Parameters []ast.Identifier
	Body       []ast.Statement
	Closure    *Environment
}

func (Function) Type() Type {
	return FunctionType
}

func (value Function) Display() string {
	return "<function " + value.Name + ">"
}

// UserBoundMethod keeps the original function and receiver even if their names change.
type UserBoundMethod struct {
	Receiver *Instance
	Function Function
}

func (UserBoundMethod) Type() Type { return MethodType }

func (method UserBoundMethod) Display() string {
	return "<bound method " + method.Function.Name + ">"
}
