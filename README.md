# GoPython

GoPython es un intérprete minimal de Python escrito en Go

## layout

```text
cmd/gopython/      command-line entry point
internal/lexer/   tokens and significant indentation ✅
internal/ast/     AST node definitions ✅
internal/parser/  parser and AST construction ✅
internal/object/  runtime values and environments
internal/builtin/ built-in functions such as print and range
internal/eval/    evaluator
internal/repl/    interactive loop
examples/         python examples
```

Pipeline

```text
source → lexer → parser/AST → evaluator/runtime → output
```

## Commands

```bash
go test ./...
go run ./cmd/gopython examples/functions.py   
go run ./cmd/gopython                         # ejecuta el REPL
```

## REPL

`gopython` abre un intérprete interactivo, donde las variables y
funciones se conservan entre entradas, las expresiones muestran su valor y los errores no cierran la sesión.

Una línea que termina en `:` abre un bloque: el REPL muestra `...` y sigue
leyendo hasta una línea vacía. Ctrl-D sale.

```text
>>> def double(x):
...     return x * 2
...
>>> double(21)
42
>>> "hola"
'hola'
>>> missing
NameError at 1:1: name "missing" is not defined
```

## Docker

```bash
docker build -t gopython .
docker run --rm gopython
```

To run another source file from the host:

```bash
docker run --rm -v "$PWD:/workspace" gopython /workspace/path/to/program.py
```

To start the REPL inside the container:

```bash
docker run --rm -it --entrypoint /gopython gopython
```
