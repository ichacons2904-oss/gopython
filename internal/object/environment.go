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
	for current := environment; current != nil; current = current.parent {
		if value, ok := current.values[name]; ok {
			return value, true
		}
	}

	return nil, false
}

func (environment *Environment) Set(name string, value Value) {
	environment.values[name] = value
}
