package object

type Iterable interface {
	Value
	Iterator() Iterator
}

type Iterator interface {
	Next() (Value, bool)
}
