# GoPython

Un intérprete de un subconjunto de Python, escrito en Go. Lee el código, lo tokeniza, arma un AST y lo recorre para ejecutarlo (tree-walking).

## Qué soporta

Enteros, floats, strings, booleanos y `None`; variables; aritmética con las reglas de Python (`/` da float, `//` y `%` redondean hacia abajo); comparaciones encadenadas; `and`, `or`, `not`; `if`/`elif`/`else`, `while`, `for` sobre `range`, `break`, `continue`; funciones con recursión y closures, `nonlocal`, `global`, `pass`; `print`, `range` y `len`. Los errores indican línea y columna.

No hay listas, diccionarios, clases ni excepciones por el momento. Apuntamos tenerlos para la entrega final. 

## Cómo probarlo

```bash
go test ./...
go run ./cmd/gopython examples/informe.py   # ejecutar un archivo
go run ./cmd/gopython                       # REPL
python3 real-tests/script.py                # tests de la cátedra
python3 benchmarks/compare.py               # benchmarks 
```


## Layout

```text
cmd/gopython/      CLI: ejecuta un archivo o abre el REPL
internal/lexer/    tokens e indentación significativa
internal/ast/      nodos del AST
internal/parser/   parser descendente recursivo
internal/object/   valores, funciones y entornos
internal/builtin/  print, range y len
internal/eval/     evaluador tree-walking
internal/repl/     loop interactivo
real-tests/        tests de la cátedra
benchmarks/        comparación con otras implementaciones
examples/          programas de ejemplo
```

## Tests de la cátedra

Tradujimos los cinco programas de `plox/real-tests` a Python y agregamos dos para floats y closures. 

Como Python no tiene scope de bloque, los bloques de Lox pasaron a ser funciones, y `a < b == c` necesita paréntesis porque en Python las comparaciones se encadenan.

## Correr un programa de ejemplo

`examples/informe.py` muestra, para cada número del 1 al 20, su numeral romano, si es primo y cuántos pasos de Collatz necesita. Usa recursión, closures, floats y strings, y da la misma salida que CPython.

## Benchmarks

En la siguiente tabla mostramos el tiempo de nuestra implementacion de gopython, comparandolos con los de otras implementaciones. Fueron corridos en un Apple M4, generado con `benchmarks/compare.py`:

| Programa | gopython | gpython | RustPython | CPython | PyPy |
| --- | --- | --- | --- | --- | --- |
| arranque | 2 ms | 1.19x | 10.56x | 7.79x | 8.37x |
| bucles | 154 ms | 1.28x | 2.10x | 0.84x | 0.12x |
| strings | 205 ms | 1.21x | 2.28x | 0.76x | 0.19x |
| floats | 308 ms | 1.23x | 2.72x | 0.74x | 0.08x |
| closures | 292 ms | 1.01x | 1.43x | 0.33x | 0.06x |
| `fib` | 345 ms | 1.10x | 1.68x | 0.19x | 0.08x |

Nuestras observaciones e interpretación de los resultados: 

- **Arranque:** gopython y gpython son binarios de Go que leen el archivo y empiezan. CPython, PyPy y RustPython primero cargan su runtime y varios módulos.
- **gpython y RustPython:** compilan a bytecode y nosotros no, pero implementan Python completo: enteros de precisión arbitraria y operadores que buscan el método según el tipo. Nosotros usamos un `int64` y un `switch`. Somos más rápidos en parte porque hacemos menos.
- **CPython:** nos gana en todo menos en el arranque. La mayor diferencia está en las llamadas a funciones (`fib`). CPython sabe al compilar dónde está cada variable local y la lee por índice; nosotros creamos un mapa por llamada y buscamos cada variable por nombre. Una pasada de resolución antes de ejecutar, como la de *Crafting Interpreters*, atacaría justo eso.
- **PyPy.** Tiene un JIT que traduce los bucles calientes a código de máquina, por eso está en otra liga.

## Diferencias con Lox

- **Lenguaje:** bloques por indentación en vez de llaves; scope de función en vez de bloque; `int` y `float` separados; `0` y `""` son falsos; el `for` recorre un iterable; `print` es una función.
- **Implementación:** en vez del patrón Visitor usamos un `switch` sobre el tipo de cada nodo. `return`, `break` y `continue` no son excepciones, porque Go no las tiene: cada sentencia devuelve cómo terminó. El lexer genera tokens de indentación con una pila de niveles. No tenemos resolver, así que `nonlocal` y `global` se resuelven al ejecutar.

## Diferencias con CPython

Los enteros son de 64 bits y, si una cuenta se pasa, dan la vuelta en vez de crecer como en Python. `True + 1` da error (en Python `bool` es un `int`), los errores de los built-ins salen como `RuntimeError`, una expresión entre paréntesis no puede seguir en otra línea, y los mensajes de error tienen formato propio. 

## REPL

`gopython` abre un intérprete interactivo, donde las variables y
funciones se conservan entre entradas, las expresiones muestran su valor y los errores no cierran la sesión. Con Ctrl-D se sale.

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

Para ejecutar otro archivo del host:

```bash
docker run --rm -v "$PWD:/workspace" gopython /workspace/path/to/program.py
```

Para abrir el REPL en el contenedor:

```bash
docker run --rm -it --entrypoint /gopython gopython
```

