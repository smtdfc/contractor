package helpers

import "unicode"

func ToCamelCase(s string) string {
	if len(s) == 0 {
		return ""
	}

	runes := []rune(s)

	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func ToPascalCase(s string) string {
	if len(s) == 0 {
		return ""
	}

	runes := []rune(s)

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}
