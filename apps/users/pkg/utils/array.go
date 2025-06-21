package utils

import "slices"

func ArraySome[T any](slice []T, predicate func(T) bool) bool {
	return slices.ContainsFunc(slice, predicate)
}

func ArrayMap[A any, B any](slice []A, fn func(A) B) []B {
	result := make([]B, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}
