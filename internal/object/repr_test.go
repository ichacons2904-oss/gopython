package object

import "testing"

func TestRepr(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{name: "integer", value: Integer{Value: 42}, want: "42"},
		{name: "boolean", value: Boolean{Value: true}, want: "True"},
		{name: "none", value: None{}, want: "None"},
		{name: "string", value: String{Value: "hola"}, want: "'hola'"},
		{name: "empty string", value: String{Value: ""}, want: "''"},
		{name: "single quote", value: String{Value: "it's"}, want: `"it's"`},
		{name: "double quote", value: String{Value: `say "hi"`}, want: `'say "hi"'`},
		{name: "both quotes", value: String{Value: `it's "hi"`}, want: `'it\'s "hi"'`},
		{name: "escapes", value: String{Value: "a\\b\nc\td\r"}, want: `'a\\b\nc\td\r'`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Repr(test.value); got != test.want {
				t.Fatalf("Repr() = %s, want %s", got, test.want)
			}
		})
	}
}
