package utils

import "time"

type targetType interface {
	~int64 | ~float64 | ~string | time.Time
}

func NilToValueType[T targetType](nv *T) T {
	if nv != nil {
		return *nv
	}

	var zero T
	return zero
}
