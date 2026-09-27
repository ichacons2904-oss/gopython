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

func TestMultilineBlocks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "function definition",
			input: "def double(x):\n    return x * 2\n\ndouble(4)\n",
			want:  ">>> ... ... >>> 8\n>>> ",
		},
		{
			name:  "if with else",
			input: "if 1 > 2:\n    print(\"a\")\nelse:\n    print(\"b\")\n\n",
			want:  ">>> ... ... ... ... b\n>>> ",
		},
		{
			name:  "nested blocks",
			input: "for i in range(3):\n    if i != 1:\n        print(i)\n\n",
			want:  ">>> ... ... ... 0\n2\n>>> ",
		},
		{
			name:  "error reports the line inside the block",
			input: "if True:\n    x = 1\n    print(missing)\n\n",
			want:  ">>> ... ... ... " + `NameError at 3:11: name "missing" is not defined` + "\n>>> ",
		},
		{
			name:  "end of input finishes the block",
			input: "if True:\n    print(1)\n",
			want:  ">>> ... ... 1\n>>> ",
		},
		{
			name:  "colon inside a string",
			input: `x = "a:"` + "\nx\n",
			want:  ">>> >>> 'a:'\n>>> ",
		},
		{
			name:  "colon inside a comment",
			input: "x = 1 # note:\nx\n",
			want:  ">>> >>> 1\n>>> ",
		},
		{
			name:  "block header with a comment",
			input: "if True: # check\n    print(1)\n\n",
			want:  ">>> ... ... 1\n>>> ",
		},
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
