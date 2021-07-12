package strutil

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// FirstRuneLower returns the same string but with the first rune in the string
// transformed to lowercase.
func FirstRuneLower(value string) string {
	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError || unicode.IsLower(first) {
		return value
	}
	return concatRuneString(unicode.ToLower(first), value[size:])
}

// FirstRuneUpper returns the same string but with the first rune in the string
// transformed to uppercase.
func FirstRuneUpper(value string) string {
	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError || unicode.IsUpper(first) {
		return value
	}
	return concatRuneString(unicode.ToUpper(first), value[size:])
}

func concatRuneString(r rune, s string) string {
	return fmt.Sprintf("%c%s", r, s)
}
