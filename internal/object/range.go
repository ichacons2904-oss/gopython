package object

import "fmt"

type Iterable interface {
	Value
	Iterator() Iterator
}

type Iterator interface {
	Next() (Value, bool)
}

type Range struct {
	Start int64
	Stop  int64
	Step  int64
}

func NewRange(start, stop, step int64) (Range, error) {
	if step == 0 {
		return Range{}, fmt.Errorf("range() arg 3 must not be zero")
	}
	return Range{Start: start, Stop: stop, Step: step}, nil
}

func (Range) Type() Type {
	return RangeType
}

func (value Range) Inspect() string {
	return fmt.Sprintf("range(%d, %d, %d)", value.Start, value.Stop, value.Step)
}

func (value Range) Iterator() Iterator {
	return &rangeIterator{
		current: value.Start,
		stop:    value.Stop,
		step:    value.Step,
	}
}

type rangeIterator struct {
	current int64
	stop    int64
	step    int64
}

func (iterator *rangeIterator) Next() (Value, bool) {
	if iterator.step > 0 && iterator.current >= iterator.stop {
		return nil, false
	}
	if iterator.step < 0 && iterator.current <= iterator.stop {
		return nil, false
	}

	value := Integer{Value: iterator.current}
	iterator.current += iterator.step
	return value, true
}
