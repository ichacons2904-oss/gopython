package object

import "testing"

func TestEnvironmentStoresValues(t *testing.T) {
	environment := NewEnvironment(nil)
	want := Integer{Value: 42}
	environment.Set("answer", want)

	got, ok := environment.Get("answer")
	if !ok {
		t.Fatal("Get() did not find stored value")
	}

	if got != want {
		t.Fatalf("Get() = %#v, want %#v", got, want)
	}
}

func TestEnvironmentReadsFromParent(t *testing.T) {
	global := NewEnvironment(nil)
	local := NewEnvironment(global)
	global.Set("answer", Integer{Value: 42})

	got, ok := local.Get("answer")
	if !ok {
		t.Fatal("Get() did not find value in parent environment")
	}

	want := Integer{Value: 42}
	if got != want {
		t.Fatalf("Get() = %#v, want %#v", got, want)
	}
}

func TestEnvironmentShadowsParentValue(t *testing.T) {
	global := NewEnvironment(nil)
	local := NewEnvironment(global)
	global.Set("answer", Integer{Value: 42})
	local.Set("answer", Integer{Value: 7})

	got, ok := local.Get("answer")
	if !ok {
		t.Fatal("Get() did not find shadowing value")
	}

	want := Integer{Value: 7}
	if got != want {
		t.Fatalf("Get() = %#v, want %#v", got, want)
	}
}

func TestEnvironmentReportsMissingValue(t *testing.T) {
	environment := NewEnvironment(nil)

	if _, ok := environment.Get("missing"); ok {
		t.Fatal("Get() found a value that was never stored")
	}
}
