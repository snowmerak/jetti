package strcase

import "strings"

func PascalToSnake(value string) string {
	buffer := strings.Builder{}
	for i, v := range value {
		if i == 0 {
			buffer.WriteRune(v + 32)
			continue
		}
		if v >= 'A' && v <= 'Z' {
			buffer.WriteRune('_')
			buffer.WriteRune(v + 32)
			continue
		}
		buffer.WriteRune(v)
	}

	return buffer.String()
}

func SnakeToPascal(value string) string {
	builder := strings.Builder{}
	for _, v := range strings.Split(value, "_") {
		builder.WriteString(strings.ToUpper(v[:1]) + v[1:])
	}

	return builder.String()
}

func CamelToSnake(value string) string {
	buffer := strings.Builder{}
	for i, v := range value {
		if i == 0 {
			buffer.WriteRune(v)
			continue
		}
		if v >= 'A' && v <= 'Z' {
			buffer.WriteRune('_')
			buffer.WriteRune(v + 32)
			continue
		}
		buffer.WriteRune(v)
	}

	return buffer.String()
}

func SnakeToCamel(value string) string {
	builder := strings.Builder{}
	for i, v := range strings.Split(value, "_") {
		if i == 0 {
			builder.WriteString(v)
			continue
		}
		builder.WriteString(strings.ToUpper(v[:1]) + v[1:])
	}

	return builder.String()
}
