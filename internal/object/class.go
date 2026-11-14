package object

type Class struct {
	Name       string
	bases      []*Class
	attributes map[string]Value
}

func NewClass(name string, bases []*Class, attributes map[string]Value) *Class {
	copiedAttributes := make(map[string]Value, len(attributes))
	for key, value := range attributes {
		copiedAttributes[key] = value
	}
	return &Class{
		Name:       name,
		bases:      append([]*Class(nil), bases...),
		attributes: copiedAttributes,
	}
}

func (*Class) Type() Type { return ClassType }

func (class *Class) Display() string { return "<class '" + class.Name + "'>" }

func (class *Class) SetAttribute(name string, value Value) {
	class.attributes[name] = value
}

func (class *Class) LookupAttribute(name string) (Value, bool) {
	return class.lookupAttribute(name, make(map[*Class]bool))
}

func (class *Class) lookupAttribute(name string, visited map[*Class]bool) (Value, bool) {
	if visited[class] {
		return nil, false
	}
	visited[class] = true
	if value, ok := class.attributes[name]; ok {
		return value, true
	}
	for _, base := range class.bases {
		if value, ok := base.lookupAttribute(name, visited); ok {
			return value, true
		}
	}
	return nil, false
}

type Instance struct {
	class      *Class
	attributes map[string]Value
}

func NewInstance(class *Class) *Instance {
	return &Instance{class: class, attributes: make(map[string]Value)}
}

func (*Instance) Type() Type { return InstanceType }

func (instance *Instance) Display() string { return "<" + instance.class.Name + " instance>" }

func (instance *Instance) Class() *Class { return instance.class }

func (instance *Instance) SetAttribute(name string, value Value) {
	instance.attributes[name] = value
}

func (instance *Instance) OwnAttribute(name string) (Value, bool) {
	value, ok := instance.attributes[name]
	return value, ok
}
