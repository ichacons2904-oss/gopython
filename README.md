# GoPython

GoPython is an educational interpreter for a deliberately small subset of
Python, implemented in Go.

The project context, language contract, architecture, and delivery plan live
in the sibling [`LengComp`](../LengComp) repository. Start with its
[`wiki/index.md`](../LengComp/wiki/index.md), then read the pages relevant to
the change you are making.

## Repository layout

```text
cmd/gopython/      command-line entry point
internal/lexer/   tokens and significant indentation
internal/ast/     AST node definitions
internal/parser/  parser and AST construction
internal/object/  runtime values and environments
internal/eval/    evaluator
internal/repl/    optional interactive loop
examples/         readable example programs
testdata/         programs and expected behavior for comparisons
```

The interpreter pipeline is being built incrementally:

```text
source → lexer → parser/AST → evaluator/runtime → output
```

## Current status

The initial lexer slice is implemented. The parser, runtime, and evaluator are
still scaffolding.

## Commands

```bash
go test ./...
go run ./cmd/gopython
```
