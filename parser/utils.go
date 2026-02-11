package parser

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func AnyToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	var builder strings.Builder
	upperNext := false

	runes := []rune(s)

	for i, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			upperNext = true
			continue
		}

		if upperNext {
			builder.WriteRune(unicode.ToUpper(r))
			upperNext = false
		} else {
			if i == 0 {
				builder.WriteRune(unicode.ToLower(r))
			} else {
				builder.WriteRune(r)
			}
		}
	}

	return builder.String()
}

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

func PrintError(err BaseError, code string) {
	if err == nil {
		return
	}

	loc := err.GetLocation()
	fmt.Printf("\033[1;31m[%s]\033[0m: %s\n", err.GetName(), err.GetMessage())

	if loc != nil {
		fmt.Printf("  \033[1;33m->\033[0m %s:%d:%d\n", loc.File, loc.Start.Line, loc.Start.Column)
		fmt.Println("   |")

		lines := strings.Split(code, "\n")
		lineIdx := loc.Start.Line - 1

		if lineIdx >= 0 && lineIdx < len(lines) {
			rawLine := lines[lineIdx]
			displayLine := strings.ReplaceAll(rawLine, "\t", "    ")

			fmt.Printf("%2d |  %s\n", loc.Start.Line, displayLine)
			padding := ""
			//tabCount := 0
			for i := 0; i < loc.Start.Column-1 && i < len(rawLine); i++ {
				if rawLine[i] == '\t' {
					padding += "    "
				} else {
					padding += " "
				}
			}

			length := 1
			if loc.End.Line == loc.Start.Line {
				length = loc.End.Column - loc.Start.Column
			}
			if length <= 0 {
				length = 1
			}

			underline := strings.Repeat("^", length)
			fmt.Printf("   |  %s\033[1;31m%s\033[0m\n", padding, underline)
		}
		fmt.Println("   |")
	}
	fmt.Println()
}
