package arrays

import "slices" // Make sure you import the slices package

// Some checks if at least one element in a slice satisfies the given predicate function.
// It leverages slices.ContainsFunc from the standard library for efficiency and readability.
func Some[T any](slice []T, predicate func(T) bool) bool {
	return slices.ContainsFunc(slice, predicate)
}
