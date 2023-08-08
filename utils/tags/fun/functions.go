package fun

import (
	"github.com/brianvoe/gofakeit/v6"
	"strconv"
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
