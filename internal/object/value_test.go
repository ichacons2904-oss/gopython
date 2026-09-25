package object

import "testing"

func TestPrimitiveValues(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		wantType Type
		wantText string
	}{
		{name: "integer", value: Integer{Value: 42}, wantType: IntegerType, wantText: "42"},
		{name: "negative integer", value: Integer{Value: -7}, wantType: IntegerType, wantText: "-7"},
		{name: "string", value: String{Value: "hello"}, wantType: StringType, wantText: "hello"},
		{name: "true", value: Boolean{Value: true}, wantType: BooleanType, wantText: "True"},
		{name: "false", value: Boolean{Value: false}, wantType: BooleanType, wantText: "False"},
		{name: "none", value: None{}, wantType: NoneType, wantText: "None"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.value.Type(); got != test.wantType {
				t.Fatalf("Type() = %q, want %q", got, test.wantType)
			}
			if got := test.value.Inspect(); got != test.wantText {
				t.Fatalf("Inspect() = %q, want %q", got, test.wantText)
			}
		})
	}
}
