package repl

import (
	"bufio"
	"fmt"
	"io"

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

		if err := execute(scanner.Text(), environment); err != nil {
			fmt.Fprintln(output, err)
		}
	}
}

func execute(line string, environment *object.Environment) error {
	tokens, err := lexer.New(line).Lex()
	if err != nil {
		return err
	}

	program, err := parser.New(tokens).Parse()
	if err != nil {
		return err
	}

	_, err = eval.EvaluateProgram(program, environment)
	return err
}
