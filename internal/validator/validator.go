package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

// Create a new empty validator
func New() *Validator {
	return &Validator{
		Errors: map[string]string{},
	}
}

// Valid reports whether the validator has errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the validator if the given key is not found.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the validator only if the given OK param is not 'OK'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if the given value is in in the list of strings.
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Matches returns tue if the given value string matches the regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all string values in the slice are unique.
func Unique[T int64 | string](values []T) bool {
	uniqueValues := map[T]bool{}
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}

// NotBlank returns true if the given value is not empty after triming spaces.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// NotEmptyList returns true if the given list has one or more items.
func NotEmptyList[T any](value []T) bool {
	return len(value) > 0
}

// MaxChars returns true when the number of runes in value is less or equal to n.
// Erroneous and short encodings are treated as single runes of width 1 byte.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
