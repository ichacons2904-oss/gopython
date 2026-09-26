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
go run ./cmd/gopython
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
