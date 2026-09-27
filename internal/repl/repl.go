package repl

import (
	"bufio"
	"fmt"
	"io"

	"gopython/internal/ast"
	"gopython/internal/builtin"
	"gopython/internal/eval"
	"gopython/internal/lexer"
	"gopython/internal/object"
	"gopython/internal/parser"
)

const prompt = ">>> "

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	environment := object.NewEnvironment(nil)
	builtin.Register(environment, output)

	for {
		fmt.Fprint(output, prompt)
		if !scanner.Scan() {
			return
		}

		value, err := execute(scanner.Text(), environment)
		if err != nil {
			fmt.Fprintln(output, err)
			continue
		}
		if value != nil {
			fmt.Fprintln(output, object.Repr(value))
		}
	}
}

func execute(line string, environment *object.Environment) (object.Value, error) {
	tokens, err := lexer.New(line).Lex()
	if err != nil {
		return nil, err
	}

	program, err := parser.New(tokens).Parse()
	if err != nil {
		return nil, err
	}

	value, err := eval.EvaluateProgram(program, environment)
	if err != nil {
		return nil, err
	}

	if !endsWithExpression(program) {
		return nil, nil
	}
	if _, isNone := value.(object.None); isNone {
		return nil, nil
	}
	return value, nil
}

func endsWithExpression(program *ast.Program) bool {
	if len(program.Statements) == 0 {
		return false
	}
	_, ok := program.Statements[len(program.Statements)-1].(*ast.ExpressionStatement)
	return ok
}
