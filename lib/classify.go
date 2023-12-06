package lib

import "unicode"

func IsNum(r rune) bool {
	return unicode.IsDigit(r)
}
