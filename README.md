# GoPython

Un intérprete de un subconjunto de Python, escrito en Go. Lee el código, lo tokeniza, arma un AST y lo recorre para ejecutarlo (tree-walking).

## Funcionalidades implementadas

gopython interpreta un subconjunto de Python. Los ejemplos siguientes muestran
las funcionalidades disponibles y las principales partes de la implementación
que intervienen en cada una.

### 1. Valores y operadores

Se soportan enteros de 64 bits, floats, strings, booleanos y `None`. Las
operaciones aritméticas incluyen `+`, `-`, `*`, `/`, `//` y `%`; la división
`/` devuelve un float, mientras que `//` y `%` siguen las reglas de redondeo de
Python. También hay comparaciones encadenadas y operadores lógicos con
cortocircuito.

```python
print(7 / 2)          # 3.5
print(-7 // 2)        # -4
print(1 < 2 < 3)      # True
print(False and 1 / 0)  # False
```

**Implementación:** el lexer reconoce literales y operadores; el parser
construye expresiones con precedencia; `internal/object` representa los valores
y el evaluator ejecuta las operaciones y comparaciones.

### 2. Control de flujo

Hay condicionales `if` / `elif` / `else`, bucles `while` y `for`, y sentencias
`break`, `continue` y `pass`. `for` recorre rangos, listas, tuplas y claves de
diccionarios. Los bloques se delimitan por indentación.

```python
for number in range(5):
    if number == 2:
        continue
    if number == 4:
        break
    print(number)      # 0, 1, 3
```

**Implementación:** el lexer genera tokens de indentación; el parser construye
los bloques; el evaluator ejecuta las sentencias y propaga las señales de
`break` y `continue`.

### 3. Funciones y scopes

Las funciones se definen con `def`, admiten recursión y argumentos posicionales,
y son valores que pueden pasarse o retornarse. Una función sin `return` explícito
retorna `None`. Las funciones anidadas pueden capturar variables de su entorno;
`nonlocal` y `global` permiten asignar en scopes exteriores.

```python
def make_counter():
    count = 0
    def increment():
        nonlocal count
        count = count + 1
        return count
    return increment

counter = make_counter()
print(counter())       # 1
print(counter())       # 2
```

**Implementación:** `Function` guarda los parámetros, el cuerpo y el entorno
donde se definió; cada llamada crea un entorno nuevo. La búsqueda de nombres es
dinámica y tiene diferencias conocidas con Python; ver [Resolución de scopes](#resolución-de-scopes).

### 4. Listas

Las listas son mutables, admiten elementos de distintos tipos y pueden anidarse.
Se pueden indexar con enteros, incluidos índices negativos, y asignar a índices
existentes. También admiten igualdad, concatenación, repetición, `in`, `not in`,
iteración y `len`.

```python
items = [1, "dos"]
items.append(3)
items[-1] = 4
print(items)           # [1, 'dos', 4]
print(4 in items)      # True
print(items + [5])     # [1, 'dos', 4, 5]
```

Los métodos implementados son `append`, `extend`, `insert`, `remove`, `pop`,
`clear` y `copy`. `list()` crea una lista vacía o copia una lista, tupla, rango o las
claves de un diccionario. Los slices `items[inicio:fin:paso]` devuelven una
lista nueva y admiten índices negativos y paso negativo. No se soportan
comprehensions; los strings todavía no son iterables ni admiten slices.

**Implementación:** el AST representa literales y acceso por índice;
`internal/object` almacena los elementos y el evaluator maneja las operaciones
y resuelve los métodos mediante un registro por tipo.

### 5. Tuplas

Las tuplas son secuencias inmutables: sus elementos no se pueden reemplazar,
aunque pueden contener valores mutables. Admiten literales con y sin paréntesis,
indexación entera (incluidos índices negativos), iteración, `len`, igualdad,
pertenencia, concatenación y repetición.

```python
pair = 1, 2
single = (1,)
print(pair[-1])          # 2
print(pair + single)     # (1, 2, 1)
print(tuple([3, 4]))     # (3, 4)
```

`tuple()` acepta los iterables soportados por gopython. Los slices de una tupla
devuelven otra tupla. La asignación múltiple desempaqueta cualquier iterable en
nombres, índices o atributos (`a, b = b, a`), con `ValueError` si no coincide la
cantidad; no se admite en el `for` (`for k, v in ...`) ni con `*resto`. Las
tuplas aún no se admiten como claves de diccionario, incluso si sus elementos
son hashables. No hay comparaciones de orden entre tuplas.

**Implementación:** el parser distingue las comas de tuplas de las que separan
argumentos o elementos; el runtime conserva sus elementos y el evaluator
implementa sus operaciones.

### 6. Diccionarios

Los diccionarios almacenan pares clave-valor mutables. Admiten acceso y asignación
por clave, igualdad, pertenencia, `len` e iteración por claves en orden de
inserción. Las claves pueden ser enteros, floats, strings, booleanos o `None`.

```python
person = {"name": "Ana", "age": 30}
person["age"] = 31
print(person["name"])  # Ana
print("age" in person) # True
print(person.get("city", "unknown"))  # unknown
```

Los métodos implementados son `get`, `keys`, `values`, `pop`, `update`, `clear`
y `copy`. `dict()` crea un diccionario vacío o una copia superficial de otro;
`keys()` y `values()` devuelven listas. En este subset, `1` y `1.0` comparten
clave, pero los booleanos son distintos de los enteros.

**Implementación:** el parser construye los literales y accesos; el runtime
almacena las entradas y el evaluator ejecuta las operaciones y métodos.

### 7. Clases y objetos

Una clase puede definir atributos y métodos. Al llamarla se crea una instancia
y se ejecuta `__init__` si existe. Los métodos accedidos mediante la instancia
quedan ligados a ella; los atributos propios tienen prioridad sobre los de la
clase.

```python
class Counter:
    def __init__(self, start):
        self.value = start
    def increment(self):
        self.value = self.value + 1
        return self.value

counter = Counter(3)
print(counter.increment())  # 4
```

Se admite herencia múltiple entre clases de usuario. La búsqueda recorre las
bases en profundidad de izquierda a derecha, por lo que puede diferir del MRO
de Python en un diamante. Heredar de `list` o `dict` produce `TypeError`;
`super()` no está implementado.

**Implementación:** el parser produce nodos para clases y atributos;
`internal/object` representa clases, instancias y métodos ligados; el evaluator
construye instancias y busca atributos en las bases.

### 8. Funciones incorporadas y errores

Están disponibles `print`, `range`, `len`, `list`, `tuple` y `dict`. Los errores indican
línea y columna; incluyen, entre otros, `NameError`, `TypeError`, `IndexError`,
`KeyError` y `AttributeError`. El lenguaje aún no permite capturar excepciones.

**Implementación:** `internal/builtin` registra las funciones incorporadas y el
evaluator propaga los errores con sus posiciones de origen.

## Qué corre y qué no

Una guía rápida para saber si un programa va a correr. No es exhaustiva; el
detalle está en la
[especificación](../lengcomp/wiki/03-especificacion-subconjunto.md#estado-de-implementación).
Lo que no está soportado da un error de sintaxis, `NameError`, `TypeError` o
`AttributeError`; nunca se ignora en silencio.

| Tema | Corre | No corre |
| --- | --- | --- |
| Valores | `int` (64 bits), `float`, `str`, `bool`, `None` | enteros grandes (`2 ** 70`), `1e3`, `0xff`, strings con triple comilla, sets |
| Operadores | `+ - * / // %`, comparaciones encadenadas, `and`/`or`/`not`, `in`/`not in` | `**`, `+=` y demás asignaciones aumentadas, `is`, operadores de bits, `:=`, `a if c else b` |
| Control de flujo | `if`/`elif`/`else`, `while`, `for`, `break`, `continue`, `pass`, bloques de una línea | `else` en bucles, `try`/`except`/`raise`, `with`, `assert`, `del`, `;` |
| Funciones | parámetros posicionales, recursión, closures, `nonlocal`, `global`, funciones como valores | argumentos por defecto o con nombre (`print(x, sep="")` incluido), `*args`, `lambda`, generadores (`yield`), decoradores, docstrings |
| Strings | literales con escapes, `+`, comparación, `len` | indexar, slices, recorrer, `in`, `*`, métodos (`upper`, `split`, ...), f-strings |
| Listas y tuplas | literales, índices negativos, slices, `+`, `*`, `append`, `extend`, `insert`, `remove`, `pop`, `clear`, `copy` | `sort`, `index`, `count`, comprensiones, literales en varias líneas |
| Diccionarios | literales, `d[k]`, `get`, `keys`, `values`, `pop`, `update`, `clear`, `copy` | `items()`, tuplas como claves, comprensiones |
| Asignación | nombres, índices, atributos, `a, b = b, a` | `a = b = 1`, `for k, v in ...`, `a, *resto = ...` |
| Clases | atributos, métodos, `__init__`, herencia múltiple | `super()`, `isinstance`, `type`, métodos especiales como `__str__` (se pueden definir, pero no se usan) |
| Built-ins | `print`, `range`, `len`, `list`, `tuple`, `dict` | `str`, `int`, `float`, `bool`, `min`, `max`, `sum`, `sorted`, `input`, `open` |
| Módulos | un solo archivo | `import`, `__name__` |

Una recursión infinita no da `RecursionError`: termina el proceso con un stack
overflow de Go.

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
internal/builtin/  print, range, len, list, tuple y dict
internal/eval/     evaluador tree-walking
internal/repl/     loop interactivo
real-tests/        tests de la cátedra
benchmarks/        comparación con otras implementaciones
examples/          programas de ejemplo
```

## Tests de la cátedra

Tradujimos los cinco programas de `plox/real-tests` a Python y agregamos casos para floats, closures, listas y diccionarios.

Como Python no tiene scope de bloque, los bloques de Lox pasaron a ser funciones, y `a < b == c` necesita paréntesis porque en Python las comparaciones se encadenan.

## Correr un programa de ejemplo

`examples/informe.py` muestra, para cada número del 1 al 20, su numeral romano, si es primo y cuántos pasos de Collatz necesita. Usa recursión, closures, floats y strings, y da la misma salida que CPython.

## Benchmarks

En la siguiente tabla mostramos el tiempo de nuestra implementacion de gopython, comparandolos con los de otras implementaciones. Fueron corridos en un Apple M4, generado con `benchmarks/compare.py`:

| Programa | gopython | gpython | RustPython | CPython | PyPy |
| --- | --- | --- | --- | --- | --- |
| arranque | 2 ms | 1.17x | 9.63x | 6.21x | 7.66x |
| bucles | 159 ms | 1.27x | 2.05x | 0.93x | 0.12x |
| strings | 218 ms | 1.17x | 2.19x | 0.88x | 0.18x |
| floats | 334 ms | 1.19x | 2.54x | 0.81x | 0.08x |
| listas y diccionarios | 222 ms | — | 1.35x | 0.47x | 0.13x |
| closures | 315 ms | 0.97x | 1.37x | 0.37x | 0.06x |
| `fib` | 373 ms | 1.05x | 1.60x | 0.17x | 0.08x |

Tiempo medio de gopython; para el resto, cuántas veces ese tiempo (menos de 1 es más rápido que gopython). 2026-11-16, CPython 3.15.0, PyPy 8.0.0, RustPython 0.5.0 y gpython 0.2.0; detalle en [`benchmarks/results/comparison.md`](benchmarks/results/comparison.md). gpython no corre `lists_dicts.py` porque sus diccionarios solo aceptan strings como claves.

### Dónde se va el tiempo

Para respaldar las conclusiones con datos, cada programa también se puede correr como benchmark de Go y perfilar:

```bash
go test ./cmd/gopython -run '^$' -bench Programs -benchmem
go test ./cmd/gopython -run '^$' -bench 'Programs/fib$' -cpuprofile cpu.out
go tool pprof -top cpu.out
```

Las asignaciones de memoria de `-benchmem` son exactas:

| Programa | Qué repite | Asignaciones | Por repetición |
| --- | --- | --- | --- |
| `fib` | 2.692.537 llamadas | 10,8 millones (1 GB) | 4 asignaciones, 384 B |
| closures | 2.000.000 llamadas | 12,0 millones (704 MB) | 6 asignaciones, 352 B |
| floats | 3.000.000 iteraciones | 18,0 millones | 6 asignaciones |
| bucles | 2.250.000 iteraciones | 6,4 millones | 2,8 asignaciones |

En `fib`, la mitad de las asignaciones son de `NewEnvironment` (el struct del entorno y su mapa), un cuarto es el slice de argumentos y el otro cuarto aparece al escribir el parámetro en el mapa. El perfil de CPU da lo mismo: más de la mitad de las muestras del intérprete pasan por `Environment.Get`/`Set` y por mapas de Go, y una de cada cinco por `NewEnvironment`. En closures, tres de cada cuatro muestras pasan por los métodos de `Environment`, incluida la redirección de `nonlocal`, que se vuelve a armar en cada llamada. En macOS el perfil de CPU de Go es impreciso (la señal de muestreo llega a cualquier hilo), así que tomamos sus porcentajes como orientativos y los cruzamos con las asignaciones.

### Conclusiones

- **Arranque:** gopython y gpython son binarios de Go que leen el archivo y empiezan. CPython, PyPy y RustPython primero cargan su runtime y varios módulos.
- **Llamadas a funciones, nuestro punto más débil:** en `fib` CPython es casi 6 veces más rápido. Cada llamada nuestra crea un entorno con un mapa, copia los argumentos y busca cada nombre recorriendo mapas (`fib` es global, así que primero falla en el entorno local). CPython decide al compilar qué nombres son locales, les asigna un índice en un arreglo del *frame* y lee `n` sin hashear nada. Una pasada de resolución antes de ejecutar, como el `Resolver` de plox pero con índices en lugar de distancias, atacaría justo esto.
- **Números:** cada `int` o `float` intermedio se guarda en una interfaz de Go, y eso lo copia al heap (6 asignaciones por iteración en floats). CPython también crea un objeto por número, pero reutiliza los enteros chicos y tiene listas de objetos libres para floats, y por eso solo nos saca entre 7% y 19% en bucles, strings y floats.
- **Listas y diccionarios:** CPython es el doble de rápido. Nuestro `d[k]` normaliza la clave y la busca en un mapa de Go, y `d.get(k)` crea un método ligado y después busca `get` por nombre en una tabla; CPython tiene instrucciones de bytecode especializadas para subíndices y llamadas a métodos, que evitan crear el objeto intermedio.
- **gpython y RustPython:** compilan a bytecode, pero implementan Python completo: enteros de precisión arbitraria y operadores que buscan el método según el tipo. Nosotros usamos un `int64` y un `switch`. Les ganamos en casi todo en parte porque hacemos menos; en `fib` y closures estamos parejos con gpython.
- **PyPy:** su JIT traduce los bucles calientes a código de máquina y especializa por tipo; queda entre 5 y 17 veces más rápido que gopython.

## Diferencias con Lox

Nos concentramos en las diferencias semánticas, es decir, en qué hace distinto un programa al ejecutarse, y en cómo cambiaron la implementación. Los ejemplos de Python los corrimos en CPython 3.15 y en gopython. Como referencia de Lox usamos [plox](https://github.com/FdelMazo/plox), el intérprete de la cátedra: la rama `main`, que ya tiene resolver, y la rama `tree-walk` para el intérprete sin resolver. plox no tiene clases; para clases tomamos el Lox del libro.

### Scope y resolución de nombres

- **El scope es la función, no el bloque.** `if`, `while` y `for` no crean scope: una variable asignada adentro sigue viva al salir, incluida la del `for`. *Implementación:* hay un entorno nuevo por llamada a función y no por bloque. plox crea uno por cada bloque y otro por cada llamada.
- **Asignar es declarar.** En Lox, `x = 2` modifica una `x` ya declarada, aunque esté en un entorno de afuera, y es un error si no existe (`Env.assign` en plox). En Python siempre crea o actualiza `x` en el scope actual, así que para modificar una variable de afuera hay que declararlo con `nonlocal` o `global`. Sin `nonlocal`, el contador clásico de closures en CPython da `UnboundLocalError`. *Implementación:* `Environment.Set` escribe siempre en el entorno actual; `nonlocal x` y `global x` registran una redirección de `x` al entorno de la función contenedora o al global.
- **El "closure bug".** Traducimos `plox/examples/closure-bug.lox` (el mismo caso está en `real-tests/2-functions.lox`) usando una función en lugar del bloque:

  ```python
  a = "global"
  def block():
      def show_a():
          print(a)
      show_a()
      a = "block"
      show_a()
  block()
  ```

  | Intérprete | Salida |
  | --- | --- |
  | plox `tree-walk`, sin resolver | `global`, `block` (el bug) |
  | plox `main`, con resolver | `global`, `global` |
  | CPython | `NameError` en la primera llamada |
  | gopython | `global`, `block` |

  En Python no es un bug porque un scope no puede ganar una variable a mitad de camino, como pasa con `var a` en un bloque de Lox. Antes de ejecutar, el compilador decide qué nombres son locales de cada función: como `block` asigna `a`, `a` es local de `block` desde la primera línea, y en la primera llamada todavía no tiene valor. Python resuelve de forma estática *a qué variable* se refiere un nombre, como el `Resolver` de plox, y lee *su valor* al ejecutar (*late binding*). Por eso, si se reasigna una global que ya existe, la función ve el valor nuevo: ese caso (`LATE BINDING` en `real-tests/2-functions.py`) da lo mismo en plox, CPython y gopython, porque plox no resuelve las globales: las busca al ejecutar.

  gopython no tiene resolver: busca cada nombre por la cadena de entornos al leerlo. Por eso, en el ejemplo de arriba se comporta como la rama `tree-walk` de plox, y no da `UnboundLocalError`. Es una diferencia que aceptamos como parte del dialecto (ver [Resolución de scopes](#resolución-de-scopes)).

  Aclaración sobre `real-tests`: el test `CLOSURE BUG` de plox no tiene una traducción fiel a Python, porque depende de declarar `var a` en un bloque. En `real-tests/2-functions.py` lo reemplazamos por `LATE BINDING`, que sí pasa. Que ese test pase no quiere decir que gopython resuelva el closure bug: con el ejemplo original, se comporta como plox sin resolver.
- **Closures dentro de un bucle.** Capturan la variable, no su valor: funciones creadas en un `for` ven todas el último valor. En Lox pasa lo mismo: el parser de plox traduce el `for` a `{ inicializador; while (condición) { cuerpo; incremento; } }`, así que la variable se declara una sola vez.
- **El cuerpo de una clase no es un scope para sus métodos.** Un método no ve los nombres del cuerpo de su clase; hay que escribir `self.x`. *Implementación:* el cuerpo de la clase se ejecuta en un entorno marcado como *de clase*, y las funciones definidas ahí guardan como closure el entorno de afuera de la clase.

### Qué es expresión y qué es statement

| Construcción | Lox | Python | gopython |
| --- | --- | --- | --- |
| Asignación | expresión (`print a = 2;` imprime `2`) | statement | statement |
| `a = b = 1` | expresión asociativa a derecha | statement especial | no implementado |
| Asignar dentro de una expresión | la asignación misma | `:=` | no implementado |
| `print` | statement | función, un valor más | función |
| Función anónima | no hay | `lambda` | no implementado |
| Condicional como expresión | no hay | `a if c else b` | no implementado |
| `break` y `continue` | no hay | sí | sí |

*Implementación:* las asignaciones son nodos de statement. El parser lee una expresión y, si después viene `=`, verifica que sea un destino válido (nombre, índice, atributo o tupla). Es el mismo truco del *l-value* de `assignment()` en plox, pero solo al principio de un statement; por eso `print(x = 1)` es un error de sintaxis. `print` es un built-in registrado en el entorno global: se puede guardar en una variable o pasar como argumento. plox implementa `return` lanzando una excepción de Python (`ReturnValue`). Go no tiene excepciones, así que en gopython `return`, `break` y `continue` no lo son: cada sentencia devuelve cómo terminó, y los bucles y las llamadas lo inspeccionan.

### Valores, iteración y objetos

- **Qué es falso:** en Lox solo `nil` y `false`. En Python también `0`, `0.0`, `""`, `[]`, `{}` y `()`.
- **Números:** Lox tiene solo doubles; plox además agrega `%`, que calcula con el de Python. Python separa `int` y `float`: `/` siempre da float, y `//` y `%` redondean hacia abajo. *Implementación:* dos tipos con promoción a float, y funciones propias para `//` y `%`, porque las de Go truncan hacia cero.
- **Igualdad y orden:** en Lox, valores de distinto tipo nunca son iguales. En Python, `1 == 1.0`, las colecciones se comparan por contenido, `"a" < "b"` es válido y `a < b < c` se encadena evaluando `b` una sola vez. *Implementación:* la comparación encadenada es un nodo con una lista de operadores, y la igualdad estructural detecta ciclos.
- **Iteración:** el `for` de Lox se traduce a un `while`. El de Python le pide un iterador a cualquier iterable. *Implementación:* interfaces `Iterable` e `Iterator`, que implementan rangos, listas, tuplas y diccionarios.
- **Colecciones:** Lox no tiene. En Python, listas y diccionarios son mutables y se comparten por referencia (`b = a` no copia), las tuplas no se pueden modificar, y las claves de un diccionario tienen que ser hasheables.
- **Clases:** el Lox del libro usa `this` implícito, `init` y herencia simple con `super`. Python usa `self` explícito, `__init__`, atributos de clase compartidos, un cuerpo de clase que puede tener cualquier código, y herencia múltiple. Un método es una función común que se liga a la instancia al leerla (`obj.metodo`). *Implementación:* en el libro, el método se liga con un entorno que contiene `this`; gopython devuelve un objeto que guarda la función y la instancia, y al llamarlo agrega la instancia como primer argumento.
- **Errores:** Python tiene excepciones que se pueden capturar. gopython implementa los tipos de error, pero no `try`/`except`, así que un error termina el programa, como en Lox.

### Implementación

Como plox, gopython es un intérprete *tree-walking* con un parser descendente recursivo, y no compila a bytecode. plox elige qué hacer con cada nodo según su tipo con `singledispatchmethod`; gopython hace lo mismo con un `switch`. El lexer genera tokens `INDENT` y `DEDENT` con una pila de niveles de indentación. A diferencia de plox, no tiene resolver: los nombres se buscan al ejecutar y `nonlocal`/`global` se resuelven con redirecciones.

## Diferencias con CPython

Los enteros son de 64 bits y, si una cuenta se pasa, dan la vuelta en vez de crecer como en Python. `True + 1` da error (en Python `bool` es un `int`), algunos errores de built-ins salen como `RuntimeError`, `keys()` y `values()` devuelven listas nuevas en vez de vistas, no hay `items()`, modificar un diccionario mientras se lo recorre no da error, una expresión entre paréntesis no puede seguir en otra línea, y los mensajes de error tienen formato propio.

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
docker build -t gopython .
docker run --rm -v "$PWD:/workspace" gopython /workspace/path/to/program.py
```

Para abrir el REPL en el contenedor:

```bash
docker run --rm -it --entrypoint /gopython gopython
```

## Resolución de scopes

La búsqueda de nombres en gopython es dinámica: cada nombre se busca por la
cadena de entornos en el momento de leerlo. Decidimos conservar las diferencias
con Python que esto produce como parte del dialecto, y no corregirlas. Si más
adelante agregamos una pasada de resolución antes de la evaluación, sería solo
para acelerar el acceso a las variables y tendría que conservar este
comportamiento.

### Referencia: LEGB en Python

En Python, la resolución léxica de nombres suele describirse como **LEGB**:

- **Local**: nombres del scope actual, como los parámetros y variables de una
  función.
- **Enclosing**: scopes de funciones que contienen a la función actual, desde
  el más cercano hacia afuera; aquí se encuentran las variables capturadas.
- **Global**: el namespace del módulo donde se definió la función.
- **Built-in**: nombres incorporados de Python, como `len` y `print`.

LEGB es una guía para entender la búsqueda, no una regla de fallback idéntica
para todos los nombres. Python analiza el cuerpo completo de una función para
clasificarlos: si asigna a un nombre, normalmente es local en toda la función,
incluso antes de esa asignación. Una lectura prematura produce
`UnboundLocalError` en vez de buscar en scopes exteriores. Las declaraciones
`global` y `nonlocal` cambian esa clasificación para todo el bloque;
`nonlocal` requiere un binding en una función contenedora.

Los cuerpos de clase tienen reglas particulares: ejecutan su código en un
namespace propio, pero ese namespace no es un scope léxico contenedor para los
métodos. Por eso una búsqueda simple LEGB no alcanza para describir todos los
casos de clases anidadas.

En gopython, una lectura puede encontrar un nombre exterior antes de que una
asignación posterior cree un nombre local. Esta diferencia se acepta como
comportamiento del dialecto; no se propone cambiarla solo por diferir de Python.
Las funciones incorporadas se registran en el entorno global actual, no en un
scope built-in separado como en Python.

Casos que difieren de CPython:

- **Lectura de un local antes de asignarlo:**

  ```python
  x = 10
  def show():
      print(x)
      x = 20
  show()
  ```

  gopython imprime `10`; CPython detecta que `x` es local a toda la función y
  produce `UnboundLocalError` porque todavía no tiene valor.

- **Asignación con el mismo nombre dentro de una clase anidada:**

  ```python
  x = "global"
  def outer():
      x = "enclosing"
      class C:
          x = x
      return C.x
  print(outer())
  ```

  gopython imprime `enclosing`; CPython imprime `global`. El cuerpo de una clase
  tiene reglas de búsqueda propias y no es un scope léxico para sus métodos.

- **Validación anticipada de `nonlocal`:**

  ```python
  print("started")
  def f():
      nonlocal missing
  f()
  ```

  gopython imprime `started` y luego informa `SyntaxError` al ejecutar la
  declaración. CPython rechaza el programa antes de imprimir, porque no existe
  un binding válido en una función contenedora.
