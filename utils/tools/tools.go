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
	var (
		chars  string
		numStr string
	)
	start := 0
	for i, c := range s {
		if unicode.IsDigit(c) {
			start = i
			break
		}
	}

	end := start
	for i, c := range s[start:] {
		if unicode.IsDigit(c) {
			numStr += string(c)
		} else {
			break
		}
		end = start + i
	}

	if start == 0 && end == 0 {
		return s, 0
	}

	if end < len(s)-1 {
		chars = s[:start] + s[end+1:]
	} else {
		chars = s[:start]
	}

	num, _ := strconv.ParseInt(numStr, 10, 64)
	return chars, int(num)
}
