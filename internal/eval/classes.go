package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateClassDefinition(statement *ast.ClassDefinition, environment *object.Environment) (completion, error) {
	bases := make([]*object.Class, 0, len(statement.Bases))
	seen := make(map[*object.Class]bool)
	for _, name := range statement.Bases {
		value, ok := environment.Get(name.Name)
		if !ok {
			return completion{}, Error{Kind: NameError, Message: fmt.Sprintf("name %q is not defined", name.Name), Position: name.Position()}
		}
		base, ok := value.(*object.Class)
		if !ok {
			return completion{}, Error{Kind: TypeError, Message: fmt.Sprintf("base %q must be a class, got %s", name.Name, value.Type()), Position: name.Position()}
		}
		if seen[base] {
			return completion{}, Error{Kind: TypeError, Message: fmt.Sprintf("duplicate base class %q", base.Name), Position: name.Position()}
		}
		seen[base] = true
		bases = append(bases, base)
	}

	classEnvironment := object.NewClassEnvironment(environment)
	outcome, err := evaluateStatements(statement.Body, classEnvironment)
	if err != nil {
		return completion{}, err
	}
	if outcome.kind != normalCompletion {
		return completion{}, Error{Kind: RuntimeError, Message: fmt.Sprintf("unexpected %s in class body", outcome.kind), Position: statement.Position()}
	}
	class := object.NewClass(statement.Name.Name, bases, classEnvironment.LocalBindings())
	environment.Set(statement.Name.Name, class)
	return completion{value: object.None{}}, nil
}

func evaluateAttributeAssignment(statement *ast.AttributeAssignment, environment *object.Environment) (completion, error) {
	value, err := evaluateExpression(statement.Value, environment)
	if err != nil {
		return completion{}, err
	}
	receiver, err := evaluateExpression(statement.Target.Object, environment)
	if err != nil {
		return completion{}, err
	}
	switch receiver := receiver.(type) {
	case *object.Class:
		receiver.SetAttribute(statement.Target.Name, value)
	case *object.Instance:
		receiver.SetAttribute(statement.Target.Name, value)
	default:
		return completion{}, Error{Kind: AttributeError, Message: fmt.Sprintf("%s object has no writable attribute %q", receiver.Type(), statement.Target.Name), Position: statement.Position()}
	}
	return completion{value: object.None{}}, nil
}
