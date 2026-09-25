package eval

import (
	"fmt"

	"gopython/internal/ast"
	"gopython/internal/object"
)

func EvaluateProgram(program *ast.Program, environment *object.Environment) (object.Value, error) {
	if environment == nil {
		environment = object.NewEnvironment(nil)
	}

	outcome, err := evaluateStatements(program.Statements, environment)
	if err != nil {
		return nil, err
	}
	if outcome.kind != normalCompletion {
		return nil, Error{
			Kind:    RuntimeError,
			Message: fmt.Sprintf("unexpected %s outside its valid context", outcome.kind),
		}
	}
	return outcome.value, nil
}

func evaluateStatement(statement ast.Statement, environment *object.Environment) (completion, error) {
	switch statement := statement.(type) {
	case *ast.ExpressionStatement:
		value, err := evaluateExpression(statement.Expression, environment)
		return completion{value: value}, err
	case *ast.Assignment:
		return evaluateAssignment(statement, environment)
	case *ast.IfStatement:
		return evaluateIfStatement(statement, environment)
	case *ast.WhileStatement:
		return evaluateWhileStatement(statement, environment)
	case *ast.ForStatement:
		return evaluateForStatement(statement, environment)
	case *ast.BreakStatement:
		return completion{kind: breakCompletion}, nil
	case *ast.ContinueStatement:
		return completion{kind: continueCompletion}, nil
	default:
		return completion{}, Error{
			Kind:     RuntimeError,
			Message:  fmt.Sprintf("unsupported statement %T", statement),
			Position: statement.Position(),
		}
	}
}

func evaluateStatements(statements []ast.Statement, environment *object.Environment) (completion, error) {
	outcome := completion{value: object.None{}}
	for _, statement := range statements {
		value, err := evaluateStatement(statement, environment)
		if err != nil {
			return completion{}, err
		}
		outcome = value
		if outcome.kind != normalCompletion {
			return outcome, nil
		}
	}
	return outcome, nil
}

func evaluateAssignment(statement *ast.Assignment, environment *object.Environment) (completion, error) {
	value, err := evaluateExpression(statement.Value, environment)
	if err != nil {
		return completion{}, err
	}

	environment.Set(statement.Name.Name, value)
	return completion{value: object.None{}}, nil
}
