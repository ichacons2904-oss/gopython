package object

import "testing"

func TestTupleOwnsStorage(t *testing.T) {
	elements := []Value{Integer{Value: 1}}
	tuple := NewTuple(elements)
	elements[0] = Integer{Value: 2}
	got, ok := tuple.At(0)
	if !ok || got != (Integer{Value: 1}) {
		t.Fatalf("element = %#v, %v; want 1", got, ok)
	}
	if Hashable(tuple) {
		t.Fatal("tuple should not be a valid dictionary key yet")
	}
}
