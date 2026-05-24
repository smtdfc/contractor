package core

import (
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Validator registry for Generic validators (registered at generated code init time).
var validatorRegistry = map[string]func(interface{}) error{}

// RegisterValidator registers a validator function for a model hash. The function should
// accept any value (map or struct) and return an error (nil on success or ValidationError/other error).
func RegisterValidator(modelHash string, fn func(interface{}) error) {
	validatorRegistry[modelHash] = fn
}

// GetRegisteredValidator returns the registered validator function for a model hash.
func GetRegisteredValidator(modelHash string) (func(interface{}) error, bool) {
	fn, ok := validatorRegistry[modelHash]
	return fn, ok
}

func NotNull(value interface{}, msg string) (bool, string) {
	if value == nil {
		return false, msg
	}
	if s, ok := value.(string); ok {
		return s != "", msg
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface:
		return !rv.IsNil(), msg
	}
	return true, msg
}

func CheckIs(value interface{}, target interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}

	return reflect.DeepEqual(value, target), msg
}

func Min(value interface{}, min float64, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	switch v := value.(type) {
	case int:
		return float64(v) >= min, msg
	case int32:
		return float64(v) >= min, msg
	case int64:
		return float64(v) >= min, msg
	case float32:
		return float64(v) >= min, msg
	case float64:
		return v >= min, msg
	case string:
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n >= min, msg
		}
	}
	return false, msg
}

func Max(value interface{}, max float64, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	switch v := value.(type) {
	case int:
		return float64(v) <= max, msg
	case int32:
		return float64(v) <= max, msg
	case int64:
		return float64(v) <= max, msg
	case float32:
		return float64(v) <= max, msg
	case float64:
		return v <= max, msg
	case string:
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n <= max, msg
		}
	}
	return false, msg
}

func Range(value interface{}, min, max float64, msg string) (bool, string) {
	if valid, _ := Min(value, min, msg); !valid {
		return false, msg
	}
	return Max(value, max, msg)
}

func Length(value interface{}, length int, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len() == length, msg
	}
	return false, msg
}

func MinLength(value interface{}, min int, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len() >= min, msg
	}
	return false, msg
}

func MaxLength(value interface{}, max int, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len() <= max, msg
	}
	return false, msg
}

func Matches(value interface{}, pattern string, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, msg
	}
	return re.MatchString(s), msg
}

func Contains(value interface{}, substr string, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	return strings.Contains(s, substr), msg
}

func StartsWith(value interface{}, prefix string, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	return strings.HasPrefix(s, prefix), msg
}

func EndsWith(value interface{}, suffix string, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	return strings.HasSuffix(s, suffix), msg
}

func In(value interface{}, list interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false, msg
	}
	for i := 0; i < rv.Len(); i++ {
		if reflect.DeepEqual(rv.Index(i).Interface(), value) {
			return true, msg
		}
	}
	return false, msg
}

func IsEmail(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return re.MatchString(s), msg
}

func IsNumber(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	switch value.(type) {
	case int, int32, int64, float32, float64:
		return true, msg
	case string:
		_, err := strconv.ParseFloat(value.(string), 64)
		return err == nil, msg
	}
	return false, msg
}

func IsURL(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	_, err := url.ParseRequestURI(s)
	return err == nil, msg
}

func IsUUID(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return re.MatchString(strings.ToLower(s)), msg
}

func IsDate(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	switch v := value.(type) {
	case time.Time:
		return !v.IsZero(), msg
	case string:
		_, err := time.Parse(time.RFC3339, v)
		if err == nil {
			return true, msg
		}
		// try date only
		_, err = time.Parse("2006-01-02", v)
		return err == nil, msg
	}
	return false, msg
}

func IsAlpha(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	re := regexp.MustCompile(`^[A-Za-z]+$`)
	return re.MatchString(s), msg
}

func IsAlnum(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	s, ok := value.(string)
	if !ok {
		return false, msg
	}
	re := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	return re.MatchString(s), msg
}

func IsBool(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	switch value.(type) {
	case bool:
		return true, msg
	case string:
		_, err := strconv.ParseBool(value.(string))
		return err == nil, msg
	}
	return false, msg
}

func IsModel(value interface{}, msg string) (bool, string) {
	if valid, _ := NotNull(value, msg); !valid {
		return false, msg
	}
	rv := reflect.ValueOf(value)
	return rv.Kind() == reflect.Struct || (rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct), msg
}

// IsZeroValue reports whether the provided value is the zero value for its type.
func IsZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return true
		}
		return IsZeroValue(rv.Elem().Interface())
	}
	zero := reflect.Zero(rv.Type()).Interface()
	return reflect.DeepEqual(rv.Interface(), zero)
}

// GetModelHash attempts to extract the model hash from a value. It looks for map keys "__k" or "_k",
// or struct fields with those names (exported or unexported). Returns empty string when not found.
func GetModelHash(value interface{}) string {
	if value == nil {
		return ""
	}

	// maps
	if m, ok := value.(map[string]interface{}); ok {
		if v, exists := m["__k"]; exists {
			if s, ok := v.(string); ok {
				return s
			}
		}
		if v, exists := m["_k"]; exists {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}

	// reflect into structs
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ""
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		// try exported field __k or _k
		if f := rv.FieldByName("__k"); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
		if f := rv.FieldByName("_k"); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}

		mt2 := rv.MethodByName("GetModelHash")
		if mt2.IsValid() {
			res := mt2.Call(nil)
			if len(res) == 1 && res[0].Kind() == reflect.String {
				return res[0].String()
			}
		}
	}

	return ""
}
