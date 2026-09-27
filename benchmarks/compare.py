import json
import os
import platform
import re
import shutil
import subprocess
import sys
import tempfile
from datetime import date

currentdir = os.path.dirname(os.path.abspath(__file__))
rootdir = os.path.dirname(currentdir)
programsdir = os.path.join(currentdir, "programs")
resultsdir = os.path.join(currentdir, "results")
gopython = os.path.join(rootdir, "bin", "gopython")

CANDIDATES = {
    "gopython": [gopython],
    "cpython": ["python3"],
    "pypy": ["pypy3"],
    "rustpython": ["rustpython"],
    "gpython": ["gpython", os.path.expanduser("~/go/bin/gpython")],
}


def build_gopython():
    subprocess.run(["go", "build", "-o", gopython, "./cmd/gopython"], cwd=rootdir, check=True)


def available_interpreters():
    interpreters = {}
    for name, paths in CANDIDATES.items():
        for path in paths:
            if shutil.which(path):
                interpreters[name] = path
                break
        else:
            print(f"skip {name}: not installed")
    return interpreters


def run(interpreter, program):
    result = subprocess.run([interpreter, program], cwd=programsdir, capture_output=True, text=True)
    return result.returncode, result.stdout


def version(name, interpreter):
    if name == "gopython":
        commit = subprocess.run(["git", "rev-parse", "--short", "HEAD"], cwd=rootdir, capture_output=True, text=True)
        return f"commit {commit.stdout.strip()}"
    if name == "gpython":
        info = subprocess.run(["go", "version", "-m", interpreter], capture_output=True, text=True).stdout
        for line in info.splitlines():
            fields = line.split()
            if fields[:2] == ["mod", "github.com/go-python/gpython"]:
                return fields[2]
        return "unknown"
    result = subprocess.run([interpreter, "--version"], capture_output=True, text=True)
    text = " ".join((result.stdout or result.stderr).split())
    language = text.split(" (")[0]
    implementation = re.search(r"\[(\w+) ([\w.]+)", text)
    if implementation:
        return f"{implementation.group(1)} {implementation.group(2)} ({language})"
    return language


def measure(program, interpreters):
    with tempfile.NamedTemporaryFile(suffix=".json") as export:
        command = ["hyperfine", "--warmup", "2", "--runs", "10", "-N", "--export-json", export.name]
        for name, interpreter in interpreters.items():
            command += ["-n", name, f"{interpreter} {program}"]
        subprocess.run(command, cwd=programsdir, check=True, stdout=subprocess.DEVNULL)
        results = json.load(open(export.name))["results"]
    return {result["command"]: (result["mean"], result["stddev"]) for result in results}


def main():
    selected = sys.argv[1:]
    build_gopython()
    interpreters = available_interpreters()
    programs = sorted(p for p in os.listdir(programsdir) if p.endswith(".py"))
    if selected:
        programs = [p for p in programs if p.removesuffix(".py") in selected]

    table = {}
    for program in programs:
        _, expected = run(interpreters["cpython"], program)
        valid = {}
        for name, interpreter in interpreters.items():
            returncode, output = run(interpreter, program)
            if returncode == 0 and output == expected:
                valid[name] = interpreter
            else:
                print(f"{program}: {name} output differs from cpython, not measured")
        print(f"measuring {program} ...")
        table[program] = measure(program, valid)

    os.makedirs(resultsdir, exist_ok=True)
    with open(os.path.join(resultsdir, "comparison.md"), "w") as report:
        report.write(f"# Comparación entre implementaciones\n\n")
        report.write(f"Fecha: {date.today()}. Máquina: {platform.machine()}, {platform.platform()}.\n\n")
        report.write("hyperfine, 2 corridas de calentamiento y 10 medidas. Tiempo medio ± desvío; entre paréntesis, relativo a CPython.\n\n")
        report.write("| Intérprete | Versión |\n| --- | --- |\n")
        for name, interpreter in interpreters.items():
            report.write(f"| {name} | {version(name, interpreter)} |\n")
        report.write("\n| Programa | " + " | ".join(interpreters) + " |\n")
        report.write("| --- |" + " --- |" * len(interpreters) + "\n")
        for program, results in table.items():
            baseline = results["cpython"][0]
            cells = []
            for name in interpreters:
                if name not in results:
                    cells.append("—")
                    continue
                mean, stddev = results[name]
                cells.append(f"{mean * 1000:.0f} ± {stddev * 1000:.0f} ms ({mean / baseline:.2f}x)")
            report.write(f"| `{program}` | " + " | ".join(cells) + " |\n")
    print(f"written {os.path.join(resultsdir, 'comparison.md')}")


main()
