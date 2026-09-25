package object

import "strconv"

type Type string

const (
	IntegerType  Type = "INTEGER"
	StringType   Type = "STRING"
	BooleanType  Type = "BOOLEAN"
	NoneType     Type = "NONE"
	BuiltinType  Type = "BUILTIN"
	FunctionType Type = "FUNCTION"
	RangeType    Type = "RANGE"
)

type Value interface {
	Type() Type
	Display() string
}

type Integer struct {
	Value int64
}

func (Integer) Type() Type {
	return IntegerType
}

func (value Integer) Display() string {
	return strconv.FormatInt(value.Value, 10)
}

type String struct {
	Value string
}

func (String) Type() Type {
	return StringType
}

func (value String) Display() string {
	return value.Value
}

type Boolean struct {
	Value bool
}

func (Boolean) Type() Type {
	return BooleanType
}

func (value Boolean) Display() string {
	if value.Value {
		return "True"
	}
	return "False"
}

type None struct{}

func (None) Type() Type {
	return NoneType
}

func (None) Display() string {
	return "None"
}

type BuiltinFunction func(args []Value) (Value, error)

type Builtin struct {
	Name     string
	Function BuiltinFunction
}

func (Builtin) Type() Type {
	return BuiltinType
}

func (value Builtin) Display() string {
	return "<built-in function " + value.Name + ">"
}
