package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gopython/internal/ast"
	"gopython/internal/builtin"
	"gopython/internal/eval"
	"gopython/internal/lexer"
	"gopython/internal/object"
	"gopython/internal/parser"
)

const (
	prompt             = ">>> "
	continuationPrompt = "... "
)

func Start(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	environment := object.NewEnvironment(nil)
	builtin.Register(environment, output)

	for {
		fmt.Fprint(output, prompt)
		source, ok := readEntry(scanner, output)
		if !ok {
			return
		}

		value, err := execute(source, environment)
		if err != nil {
			fmt.Fprintln(output, err)
			continue
		}
		if value != nil {
			fmt.Fprintln(output, object.Repr(value))
		}
	}
}

func readEntry(scanner *bufio.Scanner, output io.Writer) (string, bool) {
	if !scanner.Scan() {
		return "", false
	}
	first := scanner.Text()
	if !opensBlock(first) {
		return first, true
	}

	lines := []string{first}
	for {
		fmt.Fprint(output, continuationPrompt)
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n") + "\n", true
}

func opensBlock(line string) bool {
	tokens, err := lexer.New(line).Lex()
	if err != nil {
		return false
	}

	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i].Type {
		case lexer.Newline, lexer.Dedent, lexer.EOF:
			continue
		}
		return tokens[i].Type == lexer.Colon
	}
	return false
}

func execute(source string, environment *object.Environment) (object.Value, error) {
	tokens, err := lexer.New(source).Lex()
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
