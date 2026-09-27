# Comparación entre implementaciones

Fecha: 2026-09-27. Máquina: arm64, macOS-27.0-arm64-arm-64bit-Mach-O.

hyperfine, 2 corridas de calentamiento y 10 medidas. Tiempo medio ± desvío; entre paréntesis, relativo a CPython.

| Intérprete | Versión |
| --- | --- |
| gopython | commit 4ccde84 |
| cpython | Python 3.14.6 |
| pypy | PyPy 8.0.0 (Python 3.11.16) |
| rustpython | RustPython 0.5.0 (Python 3.14.0alpha) |
| gpython | v0.2.0 |

| Programa | gopython | cpython | pypy | rustpython | gpython |
| --- | --- | --- | --- | --- | --- |
| `closures.py` | 298 ± 2 ms (3.03x) | 98 ± 4 ms (1.00x) | 19 ± 0 ms (0.19x) | 419 ± 5 ms (4.27x) | 301 ± 11 ms (3.07x) |
| `fib.py` | 348 ± 4 ms (5.07x) | 69 ± 1 ms (1.00x) | 28 ± 0 ms (0.41x) | 586 ± 4 ms (8.52x) | 382 ± 4 ms (5.56x) |
| `floats.py` | 337 ± 19 ms (1.43x) | 236 ± 10 ms (1.00x) | 27 ± 0 ms (0.11x) | 848 ± 19 ms (3.60x) | 386 ± 5 ms (1.64x) |
| `loops.py` | 162 ± 2 ms (1.25x) | 130 ± 2 ms (1.00x) | 19 ± 0 ms (0.14x) | 324 ± 3 ms (2.49x) | 200 ± 7 ms (1.54x) |
| `startup.py` | 2 ± 0 ms (0.14x) | 15 ± 0 ms (1.00x) | 15 ± 0 ms (0.98x) | 19 ± 0 ms (1.21x) | 3 ± 0 ms (0.16x) |
| `strings.py` | 216 ± 1 ms (1.39x) | 155 ± 5 ms (1.00x) | 38 ± 0 ms (0.25x) | 471 ± 23 ms (3.03x) | 251 ± 4 ms (1.62x) |
