package object

type Environment struct {
	values map[string]Value
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]Value),
		parent: parent,
	}
}

func (environment *Environment) Get(name string) (Value, bool) {
	value, ok := environment.values[name]
	if ok {
		return value, true
	}

	if environment.parent != nil {
		return environment.parent.Get(name)
	}

	return nil, false
}

func (environment *Environment) Set(name string, value Value) {
	environment.values[name] = value
}
