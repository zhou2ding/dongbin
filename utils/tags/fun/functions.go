package fun

import (
	"github.com/brianvoe/gofakeit/v6"
	"strconv"
	"unicode"
)

func Sum[T ZdbBaseType](p []T) T {
	var sum T
	for _, v := range p {
		sum += v
	}
	return sum
}

func Const[T ZdbBaseType](p []T) T {
	return p[0]
}

// IncrByN 形参的第一个元素为要执行增加的数字，第二个元素是该数字增加多少，第三个元素是已经加了几次
func IncrByN[T ZdbNumber](p []T) T {
	ret := p[0] + p[1]
	if len(p) > 2 {
		ret += p[2] * p[1]
	}
	return ret
}

// IncrByNStr 形参的第一个元素为要执行增加的数字，第二个元素是该数字增加多少，第三个元素是已经加了几次
func IncrByNStr(p []string) string {
	num, _ := strconv.Atoi(p[0])
	step, _ := strconv.Atoi(p[1])
	ret := num + step
	if len(p) > 2 {
		idx, _ := strconv.Atoi(p[2])
		ret += idx * step
	}
	return strconv.Itoa(ret)
}

// Rand 形参的第一个元素为随机范围的下限，第二个元素为随机范围的上限
func Rand[T ZdbNumber](p []T) T {
	i := interface{}(p[0])
	switch i.(type) {
	case int:
		return T(gofakeit.IntRange(int(p[0]), int(p[1])))
	case float32:
		return T(gofakeit.Float32Range(float32(p[0]), float32(p[1])))
	case float64:
		return T(gofakeit.Float64Range(float64(p[0]), float64(p[1])))
	default:
		return 0
	}
}

// RandStr 形参的第一个元素（只能是数字或字母）为随机范围的下限，第二个元素（只能是数字或字母）为随机范围的上限
func RandStr(p []string) string {
	if unicode.IsLetter(rune(p[0][0])) {
		var s string
		for i := 0; i < len(p[0]); i++ {
			s += gofakeit.RandomString(getLetters(string(p[0][i]), string(p[1][i])))
		}
		return s
	} else if unicode.IsDigit(rune(p[0][0])) {
		min, _ := strconv.Atoi(p[0])
		max, _ := strconv.Atoi(p[1])
		return strconv.Itoa(gofakeit.IntRange(min, max))
	}
	return ""
}

func getLetters(min, max string) []string {
	start := rune(min[0])
	end := rune(max[0])
	letters := make([]string, 0)
	for i := int(start); i <= int(end); i++ {
		letters = append(letters, string(rune(i)))
	}
	return letters
}
