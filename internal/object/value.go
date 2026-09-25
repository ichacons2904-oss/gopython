package object

import "strconv"

type Type string

const (
	IntegerType Type = "INTEGER"
	StringType  Type = "STRING"
	BooleanType Type = "BOOLEAN"
	NoneType    Type = "NONE"
	BuiltinType Type = "BUILTIN"
	RangeType   Type = "RANGE"
)

type Value interface {
	Type() Type
	Inspect() string
}

type Callable interface {
	Value
	Call(args []Value) (Value, error)
}

type Integer struct {
	Value int64
}

func (Integer) Type() Type {
	return IntegerType
}

func (value Integer) Inspect() string {
	return strconv.FormatInt(value.Value, 10)
}

type String struct {
	Value string
}

func (String) Type() Type {
	return StringType
}

func (value String) Inspect() string {
	return value.Value
}

type Boolean struct {
	Value bool
}

func (Boolean) Type() Type {
	return BooleanType
}

func (value Boolean) Inspect() string {
	if value.Value {
		return "True"
	}
	return "False"
}

type None struct{}

func (None) Type() Type {
	return NoneType
}

func (None) Inspect() string {
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

func (value Builtin) Inspect() string {
	return "<built-in function " + value.Name + ">"
}

func (value Builtin) Call(args []Value) (Value, error) {
	return value.Function(args)
}
