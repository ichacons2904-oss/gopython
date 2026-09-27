import os
import subprocess
import sys

INTERPRETER = sys.argv[1].split() if len(sys.argv) > 1 else ["go", "run", "./cmd/gopython"]

currentdir = os.path.dirname(os.path.abspath(__file__))
rootdir = os.path.dirname(currentdir)
thisfile = os.path.basename(__file__)


def fail():
    print(" -------- ")
    print("|  ERROR  |")
    print(" -------- ")
    sys.exit(1)


for test_file in sorted(os.listdir(currentdir)):
    if not test_file.endswith(".py") or test_file == thisfile:
        continue

    print(f"$ {' '.join(INTERPRETER)} real-tests/{test_file}")

    result = subprocess.run(
        [*INTERPRETER, os.path.join(currentdir, test_file)],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        stdin=subprocess.DEVNULL,
        cwd=rootdir,
    )

    out = result.stdout.decode().strip()
    err = result.stderr.decode().strip()
    print(out)
    if err:
        print(err)
    print()

    if result.returncode != 0 or "ERROR".lower() in out.lower():
        fail()

print(" -------- ")
print("| Todo OK |")
print(" -------- ")
