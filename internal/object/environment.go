package object

type Environment struct {
	values    map[string]Value
	parent    *Environment
	redirects map[string]*Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]Value),
		parent: parent,
	}
}

func (environment *Environment) Get(name string) (Value, bool) {
	if target, ok := environment.redirects[name]; ok {
		value, found := target.values[name]
		return value, found
	}
	for current := environment; current != nil; current = current.parent {
		if value, ok := current.values[name]; ok {
			return value, true
		}
	}

	return nil, false
}

func (environment *Environment) Set(name string, value Value) {
	if target, ok := environment.redirects[name]; ok {
		target.values[name] = value
		return
	}
	environment.values[name] = value
}

func (environment *Environment) DeclareNonlocal(name string) bool {
	for current := environment.parent; current != nil && current.parent != nil; current = current.parent {
		if _, ok := current.values[name]; ok {
			environment.redirect(name, current)
			return true
		}
	}
	return false
}

func (environment *Environment) DeclareGlobal(name string) {
	global := environment
	for global.parent != nil {
		global = global.parent
	}
	if global != environment {
		environment.redirect(name, global)
	}
}

func (environment *Environment) redirect(name string, target *Environment) {
	if environment.redirects == nil {
		environment.redirects = make(map[string]*Environment)
	}
	environment.redirects[name] = target
}
