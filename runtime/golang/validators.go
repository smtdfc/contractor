package golang

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func IsEmail(s string) bool {
	return emailRegex.MatchString(s)
}

func IsNotEmpty(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func IsString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func IsUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

func IsUrl(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsDateString(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

func MinLength(s string, min int) bool {
	return utf8.RuneCountInString(s) >= min
}

func MaxLength(s string, max int) bool {
	return utf8.RuneCountInString(s) <= max
}

func Length(s string, min, max int) bool {
	l := utf8.RuneCountInString(s)
	return l >= min && l <= max
}

func Min[T Number](v, min T) bool {
	return v >= min
}

func Max[T Number](v, max T) bool {
	return v <= max
}

func IsInt(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func IsFloat(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}

func IsBoolean(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

func IsArray(v interface{}) bool {
	if v == nil {
		return false
	}
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
}

func ArrayMinSize[T any](l []T, min int) bool {
	return len(l) >= min
}

func ArrayMaxSize[T any](l []T, max int) bool {
	return len(l) <= max
}

func ArrayLength[T any](l []T, expected int) bool {
	return len(l) == expected
}

func IsPhoneNumber(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 9 || len(s) > 15 {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && r != '+' && r != ' ' && r != '-' {
			return false
		}
	}
	return true
}
