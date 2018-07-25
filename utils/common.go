package utils

import (
	"unicode"
)

func UpperInitial(str string) string {

	runes := []rune(str)
	return string(unicode.ToUpper(runes[0])) + str[1:]
}
