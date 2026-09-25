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
