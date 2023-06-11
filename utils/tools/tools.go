package tools

import (
	"strconv"
	"unicode"
)

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

func SplitCharAndNum(s string) (string, int) {
	start := 0
	for _, c := range s {
		if unicode.IsDigit(c) {
			break
		}
		start++
	}

	end := start
	for _, c := range s[start:] {
		if !unicode.IsDigit(c) {
			break
		}
		end++
	}
	num, _ := strconv.ParseInt(s[start:end], 10, 64)
	return s[:start] + s[end:], int(num)
}
