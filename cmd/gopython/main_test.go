package main

import (
	"bytes"
	"testing"
)

func TestRunSource(t *testing.T) {
	var output bytes.Buffer

	err := runSource(`def square(value):
    return value * value

print(square(4))
`, &output)
	if err != nil {
		t.Fatalf("runSource returned an error: %v", err)
	}
	if output.String() != "16\n" {
		t.Fatalf("unexpected output %q", output.String())
	}
}
