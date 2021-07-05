package strutil

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// FirstRuneLower returns the same string but with the first rune in the string
// transformed to lowercase.
func FirstRuneLower(value string) string {
	if len(value) == 0 {
		return value
	}
	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError {
		return value
	}
	if unicode.IsLower(first) {
		return value
	}
	return fmt.Sprintf("%c%s", unicode.ToLower(first), value[size:])
}

// FirstRuneUpper returns the same string but with the first rune in the string
// transformed to uppercase.
func FirstRuneUpper(value string) string {
	if len(value) == 0 {
		return value
	}
	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError {
		return value
	}
	if unicode.IsUpper(first) {
		return value
	}
	return fmt.Sprintf("%c%s", unicode.ToUpper(first), value[size:])
}
