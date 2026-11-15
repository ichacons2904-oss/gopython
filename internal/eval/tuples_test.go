package eval

import (
	"gopython/internal/builtin"
	"gopython/internal/object"
	"testing"
)

func TestTupleBehavior(t *testing.T) {
	tests := []struct{ name, input, want string }{
		{"empty", "()\n", "()"},
		{"singleton", "(1,)\n", "(1,)"},
		{"parentheses group", "(1)\n", "1"},
		{"bare and nested", "x = 1, (2, 3),\nx\n", "(1, (2, 3))"},
		{"bare return", "def f():\n    return 1, 2\nf()\n", "(1, 2)"},
		{"call separators", "def f(a, b):\n    return a, b\nf((1, 2), 3)\n", "((1, 2), 3)"},
		{"list separators", "[(1, 2), 3]\n", "[(1, 2), 3]"},
		{"index negative", "(1, 2, 3)[-1]\n", "3"},
		{"empty falsy", "not ()\n", "True"},
		{"nonempty truthy", "not (0,)\n", "False"},
		{"len", "len((1, 2))\n", "2"},
		{"equal nested", "([1], (2,)) == ([1.0], (2,))\n", "True"},
		{"tuple distinct from list", "(1, 2) == [1, 2]\n", "False"},
		{"membership", "2 in (1, 2) and 4 not in (1, 2)\n", "True"},
		{"concat", "(1,) + (2, 3)\n", "(1, 2, 3)"},
		{"repeat", "2 * (1, 2)\n", "(1, 2, 1, 2)"},
		{"negative repeat", "(1,) * -1\n", "()"},
		{"tuple constructor", "tuple([1, 2])\n", "(1, 2)"},
		{"tuple of dict keys", "tuple({1: 2, 3: 4})\n", "(1, 3)"},
		{"list of tuple", "list((1, 2))\n", "[1, 2]"},
		{"iteration", "total = 0\nfor x in 1, 2, 3:\n    total = total + x\ntotal\n", "6"},
		{"mutable contents", "a = [1]\nt = (a,)\na.append(2)\nt\n", "([1, 2],)"},
		{"recursive repr", "a = []\nt = (a,)\na.append(t)\nt\n", "([(...)],)"},
		{"recursive equality", "a = []\nb = []\nx = (a,)\ny = (b,)\na.append(x)\nb.append(y)\nx == y\n", "True"},
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

func TestTupleErrors(t *testing.T) {
	tests := []struct {
		name, input string
		kind        ErrorKind
	}{
		{"index type", "(1,)[\"x\"]\n", TypeError},
		{"index range", "(1,)[1]\n", IndexError},
		{"immutable", "t = (1,)\nt[0] = 2\n", TypeError},
		{"unhashable key", "{(1, 2): 3}\n", TypeError},
		{"unhashable lookup", "{}[(1,)]\n", TypeError},
		{"unhashable membership", "(1,) in {}\n", TypeError},
		{"mixed concatenation", "(1,) + [2]\n", TypeError},
		{"bad repetition", "(1,) * 1.5\n", TypeError},
		{"repetition overflow", "(1, 2) * 9223372036854775807\n", MemoryError},
		{"constructor type", "tuple(1)\n", TypeError},
		{"constructor arity", "tuple(1, 2)\n", TypeError},
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
