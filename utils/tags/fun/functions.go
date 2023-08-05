package fun

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
