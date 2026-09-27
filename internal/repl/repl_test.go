package repl

import (
	"bytes"
	"strings"
	"testing"
)

func run(input string) string {
	var output bytes.Buffer
	Start(strings.NewReader(input), &output)
	return output.String()
}

func TestStateIsKeptBetweenLines(t *testing.T) {
	got := run("x = 5\nprint(x + 1)\n")

	want := ">>> >>> 6\n>>> "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestErrorsDoNotStopTheSession(t *testing.T) {
	got := run("print(missing)\nprint(1)\n")

	want := `>>> NameError at 1:7: name "missing" is not defined` + "\n>>> 1\n>>> "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEmptyLineIsIgnored(t *testing.T) {
	got := run("\nprint(1)\n")

	want := ">>> >>> 1\n>>> "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
