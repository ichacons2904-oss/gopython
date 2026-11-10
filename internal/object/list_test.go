package object

import "testing"

func TestListIndexMutationIterationAndDisplay(t *testing.T) {
	list := NewList([]Value{Integer{Value: 1}, String{Value: "two"}})
	if got := list.Display(); got != "[1, 'two']" {
		t.Fatalf("Display() = %q, want %q", got, "[1, 'two']")
	}
	if got, ok := list.At(-1); !ok || got != (String{Value: "two"}) {
		t.Fatalf("At(-1) = %#v, %v; want string element, true", got, ok)
	}
	if !list.Set(-1, Integer{Value: 3}) {
		t.Fatal("Set(-1) returned false")
	}
	if got := list.Display(); got != "[1, 3]" {
		t.Fatalf("Display() after mutation = %q, want %q", got, "[1, 3]")
	}
	iterator := list.Iterator()
	if got, ok := iterator.Next(); !ok || got != (Integer{Value: 1}) {
		t.Fatalf("first Next() = %#v, %v", got, ok)
	}
	if got, ok := iterator.Next(); !ok || got != (Integer{Value: 3}) {
		t.Fatalf("second Next() = %#v, %v", got, ok)
	}
	if _, ok := iterator.Next(); ok {
		t.Fatal("Next() after last element returned ok")
	}
}

func TestListIndexRejectsOutOfRange(t *testing.T) {
	list := NewList([]Value{Integer{Value: 1}})
	if _, ok := list.At(-2); ok {
		t.Fatal("At(-2) succeeded for a one-element list")
	}
	if list.Set(1, Integer{Value: 2}) {
		t.Fatal("Set(1) succeeded for a one-element list")
	}
}

func TestListDisplayHandlesCycles(t *testing.T) {
	list := NewList([]Value{None{}})
	list.Elements[0] = list
	if got := list.Display(); got != "[[...]]" {
		t.Fatalf("Display() = %q, want %q", got, "[[...]]")
	}
}
