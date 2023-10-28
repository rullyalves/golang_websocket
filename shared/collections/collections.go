package collections

func Index[T comparable](vs []T, t T) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include[T comparable](vs []T, t T) bool {
	return Index(vs, t) >= 0
}

func Any[T any](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All[T any](vs []T, f func(T) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter[T any](vs []T, f func(T) bool) []T {
	vsf := make([]T, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map[T any, E any](vs []T, f func(T) E) []E {
	vsm := make([]E, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func FromList[T any](dest []any) []T {
	results := make([]T, len(dest))

	for i, v := range dest {
		results[i] = v.(T)
	}

	return results
}
