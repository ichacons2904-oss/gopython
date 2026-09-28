# Benchmarks

2026-09-27 · arm64 Darwin · hyperfine, 10 corridas. Tiempo medio de gopython; para el resto, cuántas veces ese tiempo (menos es más rápido).

| Intérprete | Versión |
| --- | --- |
| gopython | commit 3caea2c |
| gpython | v0.2.0 |
| rustpython | RustPython 0.5.0 (Python 3.14.0alpha) |
| cpython | Python 3.14.6 |
| pypy | PyPy 8.0.0 (Python 3.11.16) |

| Programa | gopython | gpython | rustpython | cpython | pypy |
| --- | --- | --- | --- | --- | --- |
| `closures.py` | 292 ms | 1.01x | 1.43x | 0.33x | 0.06x |
| `fib.py` | 345 ms | 1.10x | 1.68x | 0.19x | 0.08x |
| `floats.py` | 308 ms | 1.23x | 2.72x | 0.74x | 0.08x |
| `loops.py` | 154 ms | 1.28x | 2.10x | 0.84x | 0.12x |
| `startup.py` | 2 ms | 1.19x | 10.56x | 7.79x | 8.37x |
| `strings.py` | 205 ms | 1.21x | 2.28x | 0.76x | 0.19x |
