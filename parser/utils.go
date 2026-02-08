package parser

import (
	"reflect"
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

func isNil(x any) bool {
	if x == nil {
		return true
	}
	value := reflect.ValueOf(x)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}
