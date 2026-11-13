package object

import "math"

type DictEntry struct {
	Key   Value
	Value Value
}

type Dict struct {
	entries []DictEntry
	index   map[hashKey]int
}

// hashKey mirrors valuesEqual for hashable values: 1 and 1.0 are the same
// key, while True and 1 are not because bool is not an int in this subset.
type hashKey struct {
	kind    Type
	integer int64
	float   float64
	text    string
}

func NewDict() *Dict {
	return &Dict{index: make(map[hashKey]int)}
}

func (*Dict) Type() Type { return DictType }

func (dict *Dict) Display() string {
	return Repr(dict)
}

func Hashable(value Value) bool {
	_, ok := hashKeyOf(value)
	return ok
}

func hashKeyOf(value Value) (hashKey, bool) {
	switch value := value.(type) {
	case Integer:
		return hashKey{kind: IntegerType, integer: value.Value}, true
	case Float:
		if value.Value == math.Trunc(value.Value) && value.Value >= math.MinInt64 && value.Value < math.MaxInt64 {
			return hashKey{kind: IntegerType, integer: int64(value.Value)}, true
		}
		return hashKey{kind: FloatType, float: value.Value}, true
	case String:
		return hashKey{kind: StringType, text: value.Value}, true
	case Boolean:
		if value.Value {
			return hashKey{kind: BooleanType, integer: 1}, true
		}
		return hashKey{kind: BooleanType}, true
	case None:
		return hashKey{kind: NoneType}, true
	default:
		return hashKey{}, false
	}
}

func (dict *Dict) Len() int {
	return len(dict.entries)
}

func (dict *Dict) Get(key Value) (Value, bool) {
	hash, ok := hashKeyOf(key)
	if !ok {
		return nil, false
	}
	position, ok := dict.index[hash]
	if !ok {
		return nil, false
	}
	return dict.entries[position].Value, true
}

func (dict *Dict) Set(key, value Value) bool {
	hash, ok := hashKeyOf(key)
	if !ok {
		return false
	}
	if position, ok := dict.index[hash]; ok {
		dict.entries[position].Value = value
		return true
	}
	dict.index[hash] = len(dict.entries)
	dict.entries = append(dict.entries, DictEntry{Key: key, Value: value})
	return true
}

func (dict *Dict) Delete(key Value) (Value, bool) {
	hash, ok := hashKeyOf(key)
	if !ok {
		return nil, false
	}
	position, ok := dict.index[hash]
	if !ok {
		return nil, false
	}
	value := dict.entries[position].Value
	delete(dict.index, hash)
	dict.entries = append(dict.entries[:position], dict.entries[position+1:]...)
	for index := position; index < len(dict.entries); index++ {
		shifted, _ := hashKeyOf(dict.entries[index].Key)
		dict.index[shifted] = index
	}
	return value, true
}

func (dict *Dict) Entries() []DictEntry {
	return append([]DictEntry{}, dict.entries...)
}

func (dict *Dict) Keys() []Value {
	keys := make([]Value, len(dict.entries))
	for index, entry := range dict.entries {
		keys[index] = entry.Key
	}
	return keys
}

func (dict *Dict) Values() []Value {
	values := make([]Value, len(dict.entries))
	for index, entry := range dict.entries {
		values[index] = entry.Value
	}
	return values
}

func (dict *Dict) Clear() {
	dict.entries = nil
	dict.index = make(map[hashKey]int)
}

func (dict *Dict) Copy() *Dict {
	copied := NewDict()
	for _, entry := range dict.entries {
		copied.Set(entry.Key, entry.Value)
	}
	return copied
}

// The iterator walks a snapshot of the keys, so mutating the dict while
// looping does not raise RuntimeError as it does in CPython.
func (dict *Dict) Iterator() Iterator {
	return &listIterator{list: NewList(dict.Keys())}
}

func formatDict(dict *Dict, active map[Value]bool) string {
	if active[dict] {
		return "{...}"
	}
	active[dict] = true
	defer delete(active, dict)

	parts := make([]string, len(dict.entries))
	for index, entry := range dict.entries {
		parts[index] = repr(entry.Key, active) + ": " + repr(entry.Value, active)
	}
	return "{" + join(parts, ", ") + "}"
}
