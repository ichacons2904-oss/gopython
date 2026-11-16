package object

type Tuple struct{ elements []Value }

func NewTuple(elements []Value) *Tuple { return &Tuple{elements: append([]Value(nil), elements...)} }
func (*Tuple) Type() Type              { return TupleType }
func (tuple *Tuple) Display() string   { return Repr(tuple) }
func (tuple *Tuple) Len() int          { return len(tuple.elements) }
func (tuple *Tuple) At(index int64) (Value, bool) {
	i, ok := normalizeListIndex(index, len(tuple.elements))
	if !ok {
		return nil, false
	}
	return tuple.elements[i], true
}
func (tuple *Tuple) Iterator() Iterator { return &tupleIterator{tuple: tuple} }

type tupleIterator struct {
	tuple *Tuple
	index int
}

func (it *tupleIterator) Next() (Value, bool) {
	if it.index >= len(it.tuple.elements) {
		return nil, false
	}
	value := it.tuple.elements[it.index]
	it.index++
	return value, true
}
func formatTuple(tuple *Tuple, active map[Value]bool) string {
	if active[tuple] {
		return "(...)"
	}
	active[tuple] = true
	defer delete(active, tuple)
	parts := make([]string, len(tuple.elements))
	for i, value := range tuple.elements {
		parts[i] = repr(value, active)
	}
	if len(parts) == 1 {
		return "(" + parts[0] + ",)"
	}
	return "(" + join(parts, ", ") + ")"
}
