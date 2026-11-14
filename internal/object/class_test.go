package object

import "testing"

func TestClassLookupUsesDepthFirstBaseOrder(t *testing.T) {
	root := NewClass("Root", nil, map[string]Value{"value": Integer{Value: 1}})
	left := NewClass("Left", []*Class{root}, nil)
	right := NewClass("Right", []*Class{root}, map[string]Value{
		"value": Integer{Value: 2},
		"other": Integer{Value: 3},
	})
	child := NewClass("Child", []*Class{left, right}, nil)

	if got, ok := child.LookupAttribute("value"); !ok || got != (Integer{Value: 1}) {
		t.Fatalf("value = %#v, %v; want the left branch's inherited value", got, ok)
	}
	if got, ok := child.LookupAttribute("other"); !ok || got != (Integer{Value: 3}) {
		t.Fatalf("other = %#v, %v; want the right branch's value", got, ok)
	}
	child.SetAttribute("value", Integer{Value: 4})
	if got, ok := child.LookupAttribute("value"); !ok || got != (Integer{Value: 4}) {
		t.Fatalf("overridden value = %#v, %v; want 4", got, ok)
	}
}

func TestClassLookupHandlesRepeatedAncestors(t *testing.T) {
	root := NewClass("Root", nil, nil)
	child := NewClass("Child", []*Class{root, root}, nil)
	root.bases = []*Class{child}
	if _, ok := child.LookupAttribute("missing"); ok {
		t.Fatal("LookupAttribute() found a missing attribute")
	}
}

func TestInstancesKeepSeparateAttributes(t *testing.T) {
	class := NewClass("C", nil, nil)
	first := NewInstance(class)
	second := NewInstance(class)
	first.SetAttribute("value", Integer{Value: 1})
	if first.Class() != class || second.Class() != class {
		t.Fatal("instances do not reference their class")
	}
	if _, ok := second.OwnAttribute("value"); ok {
		t.Fatal("second instance received the first instance's attribute")
	}
}
