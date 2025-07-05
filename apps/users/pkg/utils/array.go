package utils

func ArrayMap[A any, B any](slice []A, fn func(A) B) []B {
	result := make([]B, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}
