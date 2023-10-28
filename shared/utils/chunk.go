package utils

func Chunk[T any](slice []T, size int) [][]T {

	var results [][]T

	var startChunk = 0

	sliceLength := len(slice)

	for {
		if startChunk >= sliceLength {
			break
		}

		endChunk := startChunk + size

		if endChunk >= sliceLength {
			endChunk = sliceLength
		}

		chunkSlice := slice[startChunk:endChunk:size]

		startChunk = endChunk

		results = append(results, chunkSlice)

	}

	return results
}
