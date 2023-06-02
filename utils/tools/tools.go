package tools

import "unicode"

func GetDigitOfStr(s string) string {
	var b, e int
	for i := 0; i < len(s); i++ {
		if unicode.IsDigit(rune(s[i])) {
			b = i
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if unicode.IsDigit(rune(s[i])) {
			e = i
			break
		}
	}
	return s[b : e+1]
}
