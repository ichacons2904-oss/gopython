package eval

import (
	"fmt"
	"gopython/internal/ast"
	"gopython/internal/object"
)

func evaluateUnpackAssignment(statement *ast.UnpackAssignment, environment *object.Environment) (completion, error) {
	value, err := evaluateExpression(statement.Value, environment)
	if err != nil {
		return completion{}, err
	}
	iterable, ok := value.(object.Iterable)
	if !ok {
		return completion{}, Error{Kind: TypeError, Message: fmt.Sprintf("cannot unpack non-iterable %s object", value.Type()), Position: statement.Position()}
	}
	iterator := iterable.Iterator()
	values := make([]object.Value, len(statement.Targets))
	for i := range values {
		v, ok := iterator.Next()
		if !ok {
			return completion{}, Error{Kind: ValueError, Message: fmt.Sprintf("not enough values to unpack (expected %d, got %d)", len(values), i), Position: statement.Position()}
		}
		values[i] = v
	}
	if _, ok := iterator.Next(); ok {
		return completion{}, Error{Kind: ValueError, Message: fmt.Sprintf("too many values to unpack (expected %d)", len(values)), Position: statement.Position()}
	}
	for i, target := range statement.Targets {
		switch target := target.(type) {
		case *ast.Identifier:
			environment.Set(target.Name, values[i])
		case *ast.IndexExpression:
			if _, err := assignIndex(target, values[i], environment); err != nil {
				return completion{}, err
			}
		case *ast.AttributeExpression:
			if _, err := assignAttribute(target, values[i], environment); err != nil {
				return completion{}, err
			}
		}
	}
	return completion{value: object.None{}}, nil
}

func evaluateSlice(expression *ast.SliceExpression, environment *object.Environment) (object.Value, error) {
	collection, err := evaluateExpression(expression.Collection, environment)
	if err != nil {
		return nil, err
	}
	bounds := [3]ast.Expression{expression.Start, expression.Stop, expression.Step}
	var values [3]object.Value
	for i, bound := range bounds {
		if bound != nil {
			values[i], err = evaluateExpression(bound, environment)
			if err != nil {
				return nil, err
			}
		}
	}
	var numbers [3]*int64
	for i, value := range values {
		if value == nil {
			continue
		}
		if _, ok := value.(object.None); ok {
			continue
		}
		integer, ok := value.(object.Integer)
		if !ok {
			return nil, Error{Kind: TypeError, Message: "slice indices must be integers or None", Position: bounds[i].Position()}
		}
		n := integer.Value
		numbers[i] = &n
	}
	step := int64(1)
	if numbers[2] != nil {
		step = *numbers[2]
	}
	if step == 0 {
		return nil, Error{Kind: ValueError, Message: "slice step cannot be zero", Position: expression.Position()}
	}
	var length int
	var at func(int64) object.Value
	switch sequence := collection.(type) {
	case *object.List:
		length = len(sequence.Elements)
		at = func(i int64) object.Value { return sequence.Elements[i] }
	case *object.Tuple:
		length = sequence.Len()
		at = func(i int64) object.Value { v, _ := sequence.At(i); return v }
	default:
		return nil, Error{Kind: TypeError, Message: fmt.Sprintf("%s object does not support slicing", collection.Type()), Position: expression.Position()}
	}
	size := int64(length)
	start, stop := int64(0), size
	if step < 0 {
		start, stop = size-1, -1
	}
	if numbers[0] != nil {
		start = sliceBound(*numbers[0], size, step)
	}
	if numbers[1] != nil {
		stop = sliceBound(*numbers[1], size, step)
	}
	elements := []object.Value{}
	for i := start; (step > 0 && i < stop) || (step < 0 && i > stop); {
		elements = append(elements, at(i))
		// Check the remaining distance before adding, so an extreme step cannot overflow.
		if (step > 0 && step >= stop-i) || (step < 0 && step <= stop-i) {
			break
		}
		i += step
	}
	if _, ok := collection.(*object.Tuple); ok {
		return object.NewTuple(elements), nil
	}
	return object.NewList(elements), nil
}

func sliceBound(bound, size, step int64) int64 {
	if bound < 0 {
		bound += size
	}
	if step > 0 {
		if bound < 0 {
			return 0
		}
		if bound > size {
			return size
		}
	} else {
		if bound < 0 {
			return -1
		}
		if bound >= size {
			return size - 1
		}
	}
	return bound
}
