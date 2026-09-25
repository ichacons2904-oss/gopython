package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func Evaluate(node ast.Node, environment *object.Environment) (object.Value, error) {
	if environment == nil {
		environment = object.NewEnvironment(nil)
	}

	result, err := evaluate(node, environment)
	if err != nil {
		return nil, err
	}
	if result.Flow != NoFlow {
		return nil, Error{
			Kind:    RuntimeError,
			Message: fmt.Sprintf("unexpected %s outside its valid context", result.Flow),
		}
	}
	return result.Value, nil
}

func evaluate(node ast.Node, environment *object.Environment) (Result, error) {
	switch node := node.(type) {
	case ast.Program:
		return evaluateProgram(node, environment)
	case *ast.Program:
		return evaluateProgram(*node, environment)
	case ast.ExpressionStatement:
		value, err := evaluateExpression(node.Expression, environment)
		return Result{Value: value}, err
	case ast.Assignment:
		return evaluateAssignment(node, environment)
	case ast.IfStatement:
		return evaluateIfStatement(node, environment)
	case *ast.IfStatement:
		return evaluateIfStatement(*node, environment)
	case ast.WhileStatement:
		return evaluateWhileStatement(node, environment)
	case *ast.WhileStatement:
		return evaluateWhileStatement(*node, environment)
	case ast.ForStatement:
		return evaluateForStatement(node, environment)
	case *ast.ForStatement:
		return evaluateForStatement(*node, environment)
	case ast.BreakStatement:
		return Result{Flow: BreakFlow}, nil
	case *ast.BreakStatement:
		return Result{Flow: BreakFlow}, nil
	case ast.ContinueStatement:
		return Result{Flow: ContinueFlow}, nil
	case *ast.ContinueStatement:
		return Result{Flow: ContinueFlow}, nil
	default:
		return Result{}, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported statement %T", node),
			Position: node.Position(),
		}
	}
}

func evaluateProgram(program ast.Program, environment *object.Environment) (Result, error) {
	return evaluateStatements(program.Statements, environment)
}

func evaluateStatements(statements []ast.Statement, environment *object.Environment) (Result, error) {
	result := Result{Value: object.None{}}
	for _, statement := range statements {
		value, err := evaluate(statement, environment)
		if err != nil {
			return Result{}, err
		}
		result = value
		if result.Flow != NoFlow {
			return result, nil
		}
	}
	return result, nil
}

func evaluateAssignment(statement ast.Assignment, environment *object.Environment) (Result, error) {
	value, err := evaluateExpression(statement.Value, environment)
	if err != nil {
		return Result{}, err
	}

	environment.Set(statement.Name.Name, value)
	return Result{Value: object.None{}}, nil
}
