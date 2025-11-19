package utils

type Numeric interface {
	~int32 | ~int64 | ~float32 | ~float64
}

func PtrNumeric[T Numeric, U Numeric](v T) *U {
	u := U(v)
	return &u
}
