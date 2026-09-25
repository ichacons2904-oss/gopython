package object

import "testing"

func TestRangeIterator(t *testing.T) {
	rangeValue, err := NewRange(1, 6, 2)
	if err != nil {
		t.Fatal(err)
	}

	iterator := rangeValue.Iterator()
	for _, want := range []int64{1, 3, 5} {
		value, ok := iterator.Next()
		if !ok {
			t.Fatal("iterator ended before expected")
		}
		if value != (Integer{Value: want}) {
			t.Fatalf("Next() = %#v, want Integer{%d}", value, want)
		}
	}

	if _, ok := iterator.Next(); ok {
		t.Fatal("iterator returned a value after reaching stop")
	}
}

func TestRangeRejectsZeroStep(t *testing.T) {
	if _, err := NewRange(0, 5, 0); err == nil {
		t.Fatal("NewRange() accepted a zero step")
	}
}
