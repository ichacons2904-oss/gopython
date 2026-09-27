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

func TestExpressionValuesAreEchoed(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "integer", input: "2 + 3\n", want: ">>> 5\n>>> "},
		{name: "string uses repr", input: `"hola"` + "\n", want: ">>> 'hola'\n>>> "},
		{name: "variable", input: "x = 7\nx\n", want: ">>> >>> 7\n>>> "},
		{name: "none is not echoed", input: "None\n", want: ">>> >>> "},
		{name: "assignment is not echoed", input: "x = 1\n", want: ">>> >>> "},
		{name: "print output is not repeated", input: "print(1)\n", want: ">>> 1\n>>> "},
		{name: "function", input: "print\n", want: ">>> <built-in function print>\n>>> "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := run(test.input); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestEmptyLineIsIgnored(t *testing.T) {
	got := run("\nprint(1)\n")

	want := ">>> >>> 1\n>>> "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
