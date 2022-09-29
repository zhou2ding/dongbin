package tag

import "github.com/brianvoe/gofakeit/v6"

func Sum[T ZdbType](nums []T) T {
	var sum T
	for _, n := range nums {
		sum = sum + n
	}
	return sum
}

func Const[T ZdbType](c T) T {
	return c
}

func IncrByN[T ZdbNumber](val, n T) T {
	return val + n
}

func DecrByN[T ZdbNumber](val, n T) T {
	return val - n
}

func Rand[T ZdbNumber](start, end T) T {
	switch interface{}(start).(type) {
	case int:
		return T(gofakeit.IntRange(int(start), int(end)))
	case float64:
		return T(gofakeit.Float64Range(float64(start), float64(end)))
	default:
		return 0
	}
}
