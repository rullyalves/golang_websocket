package utils

func MapKeys[K comparable](data map[K]any) []K {

	var keys []K
	for key := range data {
		keys = append(keys, key)
	}

	return keys
}
