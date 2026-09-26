package main

import (
	"fmt"
	"io"
	"os"

	"gopython/internal/builtin"
	"gopython/internal/eval"
	"gopython/internal/lexer"
	"gopython/internal/object"
	"gopython/internal/parser"
)

func main() {
	filename := "examples/functions.py"
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: gopython [file.py]")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	if err := runFile(filename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runFile(filename string, output io.Writer) error {
	source, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read %q: %w", filename, err)
	}
	return runSource(string(source), output)
}

func runSource(source string, output io.Writer) error {
	tokens, err := lexer.New(source).Lex()
	if err != nil {
		return err
	}

	program, err := parser.New(tokens).Parse()
	if err != nil {
		return err
	}

	environment := object.NewEnvironment(nil)
	builtin.Register(environment, output)
	_, err = eval.EvaluateProgram(program, environment)
	return err
}
