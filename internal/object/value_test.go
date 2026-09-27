package object

import (
	"math"
	"testing"
)

func TestPrimitiveValues(t *testing.T) {
	tenth, fifth := 0.1, 0.2
	tests := []struct {
		name     string
		value    Value
		wantType Type
		wantText string
	}{
		{name: "integer", value: Integer{Value: 42}, wantType: IntegerType, wantText: "42"},
		{name: "negative integer", value: Integer{Value: -7}, wantType: IntegerType, wantText: "-7"},
		{name: "whole float", value: Float{Value: 3}, wantType: FloatType, wantText: "3.0"},
		{name: "fractional float", value: Float{Value: -3.5}, wantType: FloatType, wantText: "-3.5"},
		{name: "shortest float", value: Float{Value: tenth + fifth}, wantType: FloatType, wantText: "0.30000000000000004"},
		{name: "large float", value: Float{Value: 1e6}, wantType: FloatType, wantText: "1000000.0"},
		{name: "huge float", value: Float{Value: 1e16}, wantType: FloatType, wantText: "1e+16"},
		{name: "tiny float", value: Float{Value: 1.5e-5}, wantType: FloatType, wantText: "1.5e-05"},
		{name: "negative zero", value: Float{Value: math.Copysign(0, -1)}, wantType: FloatType, wantText: "-0.0"},
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
			if got := test.value.Display(); got != test.wantText {
				t.Fatalf("Display() = %q, want %q", got, test.wantText)
			}
		})
	}
}
