package object

type List struct {
	Elements []Value
}

func NewList(elements []Value) *List {
	return &List{Elements: elements}
}

func (*List) Type() Type { return ListType }

func (value *List) Display() string {
	return formatList(value, make(map[*List]bool))
}

func (value *List) Iterator() Iterator {
	return &listIterator{list: value}
}

func (value *List) At(index int64) (Value, bool) {
	normalizedIndex, ok := normalizeListIndex(index, len(value.Elements))
	if !ok {
		return nil, false
	}
	return value.Elements[normalizedIndex], true
}

func (value *List) Set(index int64, element Value) bool {
	normalizedIndex, ok := normalizeListIndex(index, len(value.Elements))
	if !ok {
		return false
	}
	value.Elements[normalizedIndex] = element
	return true
}

func (value *List) Pop(index int64) (Value, bool) {
	normalizedIndex, ok := normalizeListIndex(index, len(value.Elements))
	if !ok {
		return nil, false
	}
	element := value.Elements[normalizedIndex]
	copy(value.Elements[normalizedIndex:], value.Elements[normalizedIndex+1:])
	value.Elements[len(value.Elements)-1] = nil
	value.Elements = value.Elements[:len(value.Elements)-1]
	return element, true
}

func normalizeListIndex(index int64, length int) (int, bool) {
	if index < 0 {
		index += int64(length)
	}
	if index < 0 || index >= int64(length) {
		return 0, false
	}
	return int(index), true
}

func formatList(value *List, active map[*List]bool) string {
	if active[value] {
		return "[...]"
	}
	active[value] = true
	defer delete(active, value)

	parts := make([]string, len(value.Elements))
	for index, element := range value.Elements {
		if nested, ok := element.(*List); ok {
			parts[index] = formatList(nested, active)
		} else {
			parts[index] = Repr(element)
		}
	}
	return "[" + join(parts, ", ") + "]"
}

func join(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, part := range parts[1:] {
		result += separator + part
	}
	return result
}

type listIterator struct {
	list  *List
	index int
}

func (iterator *listIterator) Next() (Value, bool) {
	if iterator.index >= len(iterator.list.Elements) {
		return nil, false
	}
	value := iterator.list.Elements[iterator.index]
	iterator.index++
	return value, true
}
