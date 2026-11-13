package object

import "strings"

func Repr(value Value) string {
	return repr(value, make(map[Value]bool))
}

func repr(value Value, active map[Value]bool) string {
	switch value := value.(type) {
	case *List:
		return formatList(value, active)
	case *Dict:
		return formatDict(value, active)
	case String:
		return quote(value.Value)
	default:
		return value.Display()
	}
}

func quote(text string) string {
	delimiter := '\''
	if strings.ContainsRune(text, '\'') && !strings.ContainsRune(text, '"') {
		delimiter = '"'
	}

	var builder strings.Builder
	builder.WriteRune(delimiter)
	for _, char := range text {
		switch char {
		case '\\':
			builder.WriteString(`\\`)
		case '\n':
			builder.WriteString(`\n`)
		case '\t':
			builder.WriteString(`\t`)
		case '\r':
			builder.WriteString(`\r`)
		case delimiter:
			builder.WriteRune('\\')
			builder.WriteRune(char)
		default:
			builder.WriteRune(char)
		}
	}
	builder.WriteRune(delimiter)
	return builder.String()
}
