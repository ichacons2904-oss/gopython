package eval

import (
	"gopython/internal/builtin"
	"gopython/internal/object"
	"testing"
)

func TestUnpackingAndSlices(t *testing.T) {
	tests := []struct{ name, input, want string }{
		{"swap", "a = 1\nb = 2\na, b = b, a\n(a, b)\n", "(2, 1)"},
		{"singleton", "a, = [5]\na\n", "5"},
		{"empty", "() = ()\n42\n", "42"},
		{"range", "a, b, c = range(3)\n(a, b, c)\n", "(0, 1, 2)"},
		{"dict keys", "a, b = {1: 2, 3: 4}\n(a, b)\n", "(1, 3)"},
		{"target order", "i = 0\na = [0, 0]\ni, a[i] = 1, 7\na\n", "[0, 7]"},
		{"attribute target", "class C:\n    pass\nc = C()\na, c.x = 4, 9\n(a, c.x)\n", "(4, 9)"},
		{"global", "a = 0\ndef f():\n    global a\n    a, b = 1, 2\nf()\na\n", "1"},
		{"nonlocal", "def outer():\n    a = 0\n    def f():\n        nonlocal a\n        a, b = 3, 4\n    f()\n    return a\nouter()\n", "3"},
		{"slice basics", "a = [0, 1, 2, 3, 4]\n(a[1:4], a[:2], a[3:], a[::2], a[::-1], a[4:0:-2])\n", "([1, 2, 3], [0, 1], [3, 4], [0, 2, 4], [4, 3, 2, 1, 0], [4, 2])"},
		{"negative bounds", "a = [0, 1, 2, 3]\n(a[-3:-1], a[:-1:-1], a[::-1], a[-100:100])\n", "([1, 2], [], [3, 2, 1, 0], [0, 1, 2, 3])"},
		{"extreme bounds", "a = [0, 1, 2]\n(a[-9223372036854775807 - 1:9223372036854775807], a[::9223372036854775807], a[::-9223372036854775807 - 1])\n", "([0, 1, 2], [0], [2])"},
		{"empty and none", "([], ())[0][::-1], (1, 2)[None:None:None], ()[::-1]\n", "([], (1, 2), ())"},
		{"copies shallow", "inner = [1]\na = [inner, 2]\nb = a[:]\nb[1] = 3\ninner.append(4)\n(a, b, (inner,)[:])\n", "([[1, 4], 2], [[1, 4], 3], ([1, 4],))"},
		{"evaluation order", "seen = []\ndef mark(n):\n    seen.append(n)\n    return n\ndef seq():\n    seen.append(0)\n    return [1, 2, 3]\nx = seq()[mark(1):mark(2):mark(3)]\nseen\n", "[0, 1, 2, 3]"},
		{"rhs once", "seen = []\ndef f():\n    seen.append(1)\n    return [2, 3]\na, b = f()\nseen\n", "[1]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnvironment(nil)
			builtin.Register(env, nil)
			result, err := EvaluateProgram(parseProgram(t, tt.input), env)
			if err != nil {
				t.Fatal(err)
			}
			if got := object.Repr(result); got != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestUnpackingAndSliceErrors(t *testing.T) {
	tests := []struct {
		name, input string
		kind        ErrorKind
	}{
		{"too few", "a, b = [1]\n", ValueError},
		{"too many", "a, b = range(1000000000)\n", ValueError},
		{"empty too many", "() = [1]\n", ValueError},
		{"non iterable", "a, b = 1\n", TypeError},
		{"zero step", "[1][::0]\n", ValueError},
		{"bad start", "[1][1.0:]\n", TypeError},
		{"bad stop", "[1][:True]\n", TypeError},
		{"bad step", "[1][::'x']\n", TypeError},
		{"unsupported receiver", "'abc'[:]\n", TypeError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnvironment(nil)
			builtin.Register(env, nil)
			_, err := EvaluateProgram(parseProgram(t, tt.input), env)
			e, ok := err.(Error)
			if !ok || e.Kind != tt.kind {
				t.Fatalf("got %#v, want %s", err, tt.kind)
			}
		})
	}
}

func TestUnpackLengthErrorDoesNotAssignTargets(t *testing.T) {
	env := object.NewEnvironment(nil)
	builtin.Register(env, nil)
	_, err := EvaluateProgram(parseProgram(t, "a = 10\nb = 20\na, b = [1]\n"), env)
	if err == nil {
		t.Fatal("expected error")
	}
	for name, want := range map[string]int64{"a": 10, "b": 20} {
		v, ok := env.Get(name)
		if !ok || v != (object.Integer{Value: want}) {
			t.Fatalf("%s = %v, want %d", name, v, want)
		}
	}
}
