# Benchmarks

2026-11-16 · arm64 Darwin · hyperfine, 10 corridas. Tiempo medio de gopython; para el resto, cuántas veces ese tiempo (menos es más rápido).

| Intérprete | Versión |
| --- | --- |
| gopython | commit 2c198bf |
| gpython | v0.2.0 |
| rustpython | RustPython 0.5.0 (Python 3.14.0alpha) |
| cpython | Python 3.15.0 |
| pypy | PyPy 8.0.0 (Python 3.11.16) |

| Programa | gopython | gpython | rustpython | cpython | pypy |
| --- | --- | --- | --- | --- | --- |
| `closures.py` | 315 ms | 0.97x | 1.37x | 0.37x | 0.06x |
| `fib.py` | 373 ms | 1.05x | 1.60x | 0.17x | 0.08x |
| `floats.py` | 334 ms | 1.19x | 2.54x | 0.81x | 0.08x |
| `lists_dicts.py` | 222 ms | — | 1.35x | 0.47x | 0.13x |
| `loops.py` | 159 ms | 1.27x | 2.05x | 0.93x | 0.12x |
| `startup.py` | 2 ms | 1.17x | 9.63x | 6.21x | 7.66x |
| `strings.py` | 218 ms | 1.17x | 2.19x | 0.88x | 0.18x |
