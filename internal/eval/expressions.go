package eval

import (
	"fmt"
	"math"

	"gopython/internal/ast"
	"gopython/internal/lexer"
	"gopython/internal/object"
)

func evaluateExpression(expression ast.Expression, environment *object.Environment) (object.Value, error) {
	switch expression := expression.(type) {
	case *ast.IntegerLiteral:
		return object.Integer{Value: expression.Value}, nil
	case *ast.FloatLiteral:
		return object.Float{Value: expression.Value}, nil
	case *ast.StringLiteral:
		return object.String{Value: expression.Value}, nil
	case *ast.BooleanLiteral:
		return object.Boolean{Value: expression.Value}, nil
	case *ast.NoneLiteral:
		return object.None{}, nil
	case *ast.ListLiteral:
		elements := make([]object.Value, len(expression.Elements))
		for index, element := range expression.Elements {
			value, err := evaluateExpression(element, environment)
			if err != nil {
				return nil, err
			}
			elements[index] = value
		}
		return object.NewList(elements), nil
	case *ast.IndexExpression:
		list, index, err := evaluateListIndex(expression, environment)
		if err != nil {
			return nil, err
		}
		value, ok := list.At(index)
		if !ok {
			return nil, Error{Kind: IndexError, Message: "list index out of range", Position: expression.Position()}
		}
		return value, nil
	case *ast.AttributeExpression:
		return evaluateAttribute(expression, environment)
	case *ast.Identifier:
		return evaluateIdentifier(expression, environment)
	case *ast.UnaryExpression:
		return evaluateUnary(expression, environment)
	case *ast.BinaryExpression:
		return evaluateBinary(expression, environment)
	case *ast.ComparisonExpression:
		return evaluateComparison(expression, environment)
	case *ast.CallExpression:
		return evaluateCall(expression, environment)
	default:
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported expression %T", expression),
			Position: expression.Position(),
		}
	}
}

func evaluateAttribute(expression *ast.AttributeExpression, environment *object.Environment) (object.Value, error) {
	receiver, err := evaluateExpression(expression.Object, environment)
	if err != nil {
		return nil, err
	}
	list, ok := receiver.(*object.List)
	if !ok || !isListMethod(expression.Name) {
		return nil, Error{
			Kind:     AttributeError,
			Message:  fmt.Sprintf("%s object has no attribute %q", receiver.Type(), expression.Name),
			Position: expression.Position(),
		}
	}
	return object.BoundMethod{Receiver: list, Name: expression.Name}, nil
}

func isListMethod(name string) bool {
	switch name {
	case "extend", "append", "pop", "clear", "copy":
		return true
	default:
		return false
	}
}

func evaluateListIndex(expression *ast.IndexExpression, environment *object.Environment) (*object.List, int64, error) {
	collection, err := evaluateExpression(expression.Collection, environment)
	if err != nil {
		return nil, 0, err
	}
	indexValue, err := evaluateExpression(expression.Index, environment)
	if err != nil {
		return nil, 0, err
	}
	list, ok := collection.(*object.List)
	if !ok {
		return nil, 0, Error{Kind: TypeError, Message: fmt.Sprintf("%s object is not subscriptable", collection.Type()), Position: expression.Position()}
	}
	index, ok := indexValue.(object.Integer)
	if !ok {
		return nil, 0, Error{Kind: TypeError, Message: fmt.Sprintf("list indices must be integers, not %s", indexValue.Type()), Position: expression.Index.Position()}
	}
	return list, index.Value, nil
}

func evaluateIdentifier(expression *ast.Identifier, environment *object.Environment) (object.Value, error) {
	value, ok := environment.Get(expression.Name)
	if !ok {
		return nil, Error{
			Kind:     NameError,
			Message:  fmt.Sprintf("name %q is not defined", expression.Name),
			Position: expression.Position(),
		}
	}
	return value, nil
}

func evaluateUnary(expression *ast.UnaryExpression, environment *object.Environment) (object.Value, error) {
	operand, err := evaluateExpression(expression.Operand, environment)
	if err != nil {
		return nil, err
	}

	if expression.Operator == lexer.Not {
		return object.Boolean{Value: !isTruthy(operand)}, nil
	}
	if expression.Operator != lexer.Plus && expression.Operator != lexer.Minus {
		return nil, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported unary operator %s", expression.Operator),
			Position: expression.Position(),
		}
	}

	switch operand := operand.(type) {
	case object.Integer:
		if expression.Operator == lexer.Minus {
			return object.Integer{Value: -operand.Value}, nil
		}
		return operand, nil
	case object.Float:
		if expression.Operator == lexer.Minus {
			return object.Float{Value: -operand.Value}, nil
		}
		return operand, nil
	default:
		return nil, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("unary operator %s requires a number, got %s", expression.Operator, operand.Type()),
			Position: expression.Position(),
		}
	}
}

func evaluateBinary(expression *ast.BinaryExpression, environment *object.Environment) (object.Value, error) {
	if expression.Operator == lexer.And || expression.Operator == lexer.Or {
		return evaluateLogical(expression, environment)
	}

	left, err := evaluateExpression(expression.Left, environment)
	if err != nil {
		return nil, err
	}

	right, err := evaluateExpression(expression.Right, environment)
	if err != nil {
		return nil, err
	}

	if expression.Operator == lexer.Plus {
		leftList, leftIsList := left.(*object.List)
		rightList, rightIsList := right.(*object.List)
		if leftIsList && rightIsList {
			maxInt := int(^uint(0) >> 1)
			if len(leftList.Elements) > maxInt-len(rightList.Elements) {
				return nil, Error{Kind: MemoryError, Message: "list concatenation result is too large", Position: expression.Position()}
			}
			elements := make([]object.Value, 0, len(leftList.Elements)+len(rightList.Elements))
			elements = append(elements, leftList.Elements...)
			elements = append(elements, rightList.Elements...)
			return object.NewList(elements), nil
		}
		if leftIsList || rightIsList {
			return nil, Error{Kind: TypeError, Message: fmt.Sprintf("operator PLUS requires two lists, got %s and %s", left.Type(), right.Type()), Position: expression.Position()}
		}
		if leftString, ok := left.(object.String); ok {
			if rightString, ok := right.(object.String); ok {
				return object.String{Value: leftString.Value + rightString.Value}, nil
			}
		}
	}
	if expression.Operator == lexer.Asterisk {
		if list, ok := left.(*object.List); ok {
			return repeatList(expression, list, right)
		}
		if list, ok := right.(*object.List); ok {
			return repeatList(expression, list, left)
		}
	}

	if leftInteger, ok := left.(object.Integer); ok {
		if rightInteger, ok := right.(object.Integer); ok {
			return evaluateIntegerArithmetic(expression, leftInteger.Value, rightInteger.Value)
		}
	}

	leftNumber, leftIsNumber := toFloat(left)
	rightNumber, rightIsNumber := toFloat(right)
	if !leftIsNumber || !rightIsNumber {
		return nil, Error{
			Kind: TypeError,
			Message: fmt.Sprintf(
				"operator %s requires numbers, got %s and %s",
				expression.Operator,
				left.Type(),
				right.Type(),
			),
			Position: expression.Position(),
		}
	}
	return evaluateFloatArithmetic(expression, leftNumber, rightNumber)
}

func repeatList(expression *ast.BinaryExpression, list *object.List, countValue object.Value) (object.Value, error) {
	count, ok := countValue.(object.Integer)
	if !ok {
		return nil, Error{
			Kind:     TypeError,
			Message:  fmt.Sprintf("can't multiply list by %s", countValue.Type()),
			Position: expression.Position(),
		}
	}
	if count.Value <= 0 || len(list.Elements) == 0 {
		return object.NewList([]object.Value{}), nil
	}
	maxInt := int64(int(^uint(0) >> 1))
	if int64(len(list.Elements)) > maxInt/count.Value {
		return nil, Error{Kind: MemoryError, Message: "list repetition result is too large", Position: expression.Position()}
	}
	repeated := make([]object.Value, 0, int64(len(list.Elements))*count.Value)
	for i := int64(0); i < count.Value; i++ {
		repeated = append(repeated, list.Elements...)
	}
	return object.NewList(repeated), nil
}

func evaluateIntegerArithmetic(expression *ast.BinaryExpression, left, right int64) (object.Value, error) {
	switch expression.Operator {
	case lexer.Plus:
		return object.Integer{Value: left + right}, nil
	case lexer.Minus:
		return object.Integer{Value: left - right}, nil
	case lexer.Asterisk:
		return object.Integer{Value: left * right}, nil
	case lexer.Slash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: float64(left) / float64(right)}, nil
	case lexer.DoubleSlash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Integer{Value: floorDivide(left, right)}, nil
	case lexer.Percent:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Integer{Value: floorModulo(left, right)}, nil
	default:
		return nil, unsupportedBinaryOperator(expression)
	}
}

func evaluateFloatArithmetic(expression *ast.BinaryExpression, left, right float64) (object.Value, error) {
	switch expression.Operator {
	case lexer.Plus:
		return object.Float{Value: left + right}, nil
	case lexer.Minus:
		return object.Float{Value: left - right}, nil
	case lexer.Asterisk:
		return object.Float{Value: left * right}, nil
	case lexer.Slash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: left / right}, nil
	case lexer.DoubleSlash:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: floorDivideFloat(left, right)}, nil
	case lexer.Percent:
		if right == 0 {
			return nil, divisionByZero(expression)
		}
		return object.Float{Value: floorModuloFloat(left, right)}, nil
	default:
		return nil, unsupportedBinaryOperator(expression)
	}
}

func toFloat(value object.Value) (float64, bool) {
	switch value := value.(type) {
	case object.Integer:
		return float64(value.Value), true
	case object.Float:
		return value.Value, true
	default:
		return 0, false
	}
}

func divisionByZero(expression *ast.BinaryExpression) error {
	return Error{
		Kind:     ZeroDivisionError,
		Message:  "division by zero",
		Position: expression.Position(),
	}
}

func unsupportedBinaryOperator(expression *ast.BinaryExpression) error {
	return Error{
		Kind:     RuntimeError,
		Message:  fmt.Sprintf("unsupported binary operator %s", expression.Operator),
		Position: expression.Position(),
	}
}

func floorModulo(dividend, divisor int64) int64 {
	remainder := dividend % divisor
	if remainder != 0 && (remainder < 0) != (divisor < 0) {
		remainder += divisor
	}
	return remainder
}

func floorModuloFloat(dividend, divisor float64) float64 {
	remainder := math.Mod(dividend, divisor)
	if remainder == 0 {
		return math.Copysign(0, divisor)
	}
	if (remainder < 0) != (divisor < 0) {
		remainder += divisor
	}
	return remainder
}

func floorDivide(dividend, divisor int64) int64 {
	quotient := dividend / divisor
	if dividend%divisor != 0 && (dividend < 0) != (divisor < 0) {
		quotient--
	}
	return quotient
}

func floorDivideFloat(dividend, divisor float64) float64 {
	remainder := math.Mod(dividend, divisor)
	quotient := (dividend - remainder) / divisor
	if remainder != 0 && (remainder < 0) != (divisor < 0) {
		quotient--
	}
	if quotient == 0 {
		return math.Copysign(0, dividend/divisor)
	}
	floored := math.Floor(quotient)
	if quotient-floored > 0.5 {
		floored++
	}
	return floored
}
