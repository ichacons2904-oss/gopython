package object

import "strings"

func Repr(value Value) string {
	if value, ok := value.(String); ok {
		return quote(value.Value)
	}
	return value.Display()
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
