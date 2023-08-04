package fun

type ZdbBaseType interface {
	ZdbNumber | ~string
}

type ZdbArray interface {
	~[]string | ~[]int | ~[]uint64 | ~[]int64 | ~[]float32 | ~[]float64
}

type ZdbNumber interface {
	~int | ~uint64 | ~int64 | ~float32 | ~float64
}
