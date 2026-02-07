package main

import (
	"strings"
	"unicode"
)

func AnyToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	upperNext := true

	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			upperNext = true
			continue
		}

		if upperNext {
			result.WriteRune(unicode.ToUpper(r))
			upperNext = false
		} else {
			if unicode.IsUpper(r) {
				result.WriteRune(r)
			} else {
				result.WriteRune(r)
			}
		}
	}

	res := result.String()
	if len(res) > 0 {
		return strings.ToUpper(res[:1]) + res[1:]
	}
	return res
}
