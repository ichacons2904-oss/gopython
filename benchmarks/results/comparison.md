# Benchmarks

2026-09-27 · arm64 Darwin · hyperfine, 10 corridas. Media en ms; entre paréntesis, relativo a CPython.

| Intérprete | Versión |
| --- | --- |
| gopython | commit 934cdaf |
| cpython | Python 3.14.6 |
| pypy | PyPy 8.0.0 (Python 3.11.16) |
| rustpython | RustPython 0.5.0 (Python 3.14.0alpha) |
| gpython | v0.2.0 |

| Programa | gopython | cpython | pypy | rustpython | gpython |
| --- | --- | --- | --- | --- | --- |
| `closures.py` | 295 (3.07x) | 96 (1.00x) | 18 (0.19x) | 420 (4.37x) | 297 (3.09x) |
| `fib.py` | 344 (5.13x) | 67 (1.00x) | 27 (0.41x) | 583 (8.68x) | 380 (5.67x) |
| `floats.py` | 312 (1.34x) | 233 (1.00x) | 26 (0.11x) | 847 (3.64x) | 383 (1.65x) |
| `loops.py` | 155 (1.20x) | 130 (1.00x) | 19 (0.14x) | 325 (2.51x) | 199 (1.53x) |
| `startup.py` | 2 (0.12x) | 14 (1.00x) | 15 (1.09x) | 19 (1.35x) | 2 (0.16x) |
| `strings.py` | 206 (1.37x) | 151 (1.00x) | 38 (0.25x) | 478 (3.17x) | 251 (1.66x) |
