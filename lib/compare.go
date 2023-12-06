package lib

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Abs[T constraints.Signed](a T) T {
	if a < T(0) {
		return T(-1) * a
	}
	return a
}
