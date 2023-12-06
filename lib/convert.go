package lib

import (
	"strings"
)

func RuneToDigit(r rune) int {
	switch r {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	default:
		return 0
	}
}

func WordToNum(s string) (int, bool) {
	switch {
	case strings.HasSuffix(s, "one"):
		return 1, true
	case strings.HasSuffix(s, "two"):
		return 2, true
	case strings.HasSuffix(s, "three"):
		return 3, true
	case strings.HasSuffix(s, "four"):
		return 4, true
	case strings.HasSuffix(s, "five"):
		return 5, true
	case strings.HasSuffix(s, "six"):
		return 6, true
	case strings.HasSuffix(s, "seven"):
		return 7, true
	case strings.HasSuffix(s, "eight"):
		return 8, true
	case strings.HasSuffix(s, "nine"):
		return 9, true
	}
	return 0, false
}
