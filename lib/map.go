package lib

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func Keys[K constraints.Ordered, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}
