package object

import "testing"

func TestDictKeepsInsertionOrder(t *testing.T) {
	dict := NewDict()
	dict.Set(String{Value: "b"}, Integer{Value: 1})
	dict.Set(String{Value: "a"}, Integer{Value: 2})
	dict.Set(String{Value: "b"}, Integer{Value: 3})
	if got := dict.Display(); got != "{'b': 3, 'a': 2}" {
		t.Fatalf("Display() = %q, want %q", got, "{'b': 3, 'a': 2}")
	}
	if dict.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", dict.Len())
	}
}

func TestDictNumericKeysMatchEquality(t *testing.T) {
	dict := NewDict()
	dict.Set(Integer{Value: 1}, String{Value: "int"})
	dict.Set(Float{Value: 1.0}, String{Value: "float"})
	if dict.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", dict.Len())
	}
	if got := dict.Display(); got != "{1: 'float'}" {
		t.Fatalf("Display() = %q, want %q", got, "{1: 'float'}")
	}
	if _, ok := dict.Get(Boolean{Value: true}); ok {
		t.Fatal("Get(True) found the key 1")
	}
	if _, ok := dict.Get(Float{Value: 1.5}); ok {
		t.Fatal("Get(1.5) found the key 1")
	}
}

func TestDictRejectsUnhashableKeys(t *testing.T) {
	dict := NewDict()
	if dict.Set(NewList(nil), None{}) {
		t.Fatal("Set() accepted a list key")
	}
	if dict.Set(NewDict(), None{}) {
		t.Fatal("Set() accepted a dict key")
	}
	if Hashable(NewList(nil)) || !Hashable(None{}) {
		t.Fatal("Hashable() misclassified a value")
	}
}

func TestDictDeleteReindexesLaterEntries(t *testing.T) {
	dict := NewDict()
	for index, name := range []string{"a", "b", "c"} {
		dict.Set(String{Value: name}, Integer{Value: int64(index)})
	}
	if value, ok := dict.Delete(String{Value: "a"}); !ok || value != (Integer{Value: 0}) {
		t.Fatalf("Delete('a') = %#v, %v", value, ok)
	}
	if _, ok := dict.Delete(String{Value: "a"}); ok {
		t.Fatal("second Delete('a') succeeded")
	}
	if value, ok := dict.Get(String{Value: "c"}); !ok || value != (Integer{Value: 2}) {
		t.Fatalf("Get('c') = %#v, %v", value, ok)
	}
	dict.Set(String{Value: "c"}, Integer{Value: 5})
	if got := dict.Display(); got != "{'b': 1, 'c': 5}" {
		t.Fatalf("Display() = %q, want %q", got, "{'b': 1, 'c': 5}")
	}
}

func TestDictIteratesOverKeySnapshot(t *testing.T) {
	dict := NewDict()
	dict.Set(Integer{Value: 1}, None{})
	iterator := dict.Iterator()
	dict.Set(Integer{Value: 2}, None{})
	if got, ok := iterator.Next(); !ok || got != (Integer{Value: 1}) {
		t.Fatalf("first Next() = %#v, %v", got, ok)
	}
	if _, ok := iterator.Next(); ok {
		t.Fatal("Next() returned a key added after the iterator was created")
	}
}

func TestDictCopyIsIndependent(t *testing.T) {
	dict := NewDict()
	dict.Set(Integer{Value: 1}, None{})
	copied := dict.Copy()
	copied.Set(Integer{Value: 2}, None{})
	if dict.Len() != 1 || copied.Len() != 2 {
		t.Fatalf("lengths = %d, %d; want 1, 2", dict.Len(), copied.Len())
	}
}

func TestReprHandlesNestedAndCyclicContainers(t *testing.T) {
	dict := NewDict()
	list := NewList([]Value{dict})
	dict.Set(String{Value: "self"}, dict)
	dict.Set(String{Value: "list"}, list)
	if got := list.Display(); got != "[{'self': {...}, 'list': [...]}]" {
		t.Fatalf("Display() = %q", got)
	}
}
