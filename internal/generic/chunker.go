//    HipparchiaGoBuilder
//    Copyright: E Gunderson 2025-26
//    License: GNU GENERAL PUBLIC LICENSE 3
//        (see LICENSE in the top level directory of the distribution)

package generic

// ChunkSlice - Function to split a slice into chunks of size chunksize
func ChunkSlice[T any](slice []T, chunksize int) [][]T {
	var result [][]T

	// Iterate over the slice in steps of chunksize
	for i := 0; i < len(slice); i += chunksize {
		end := i + chunksize
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}

	return result
}

// SplitIntoNSlices - Function to split a slice into N approximately equal slices
func SplitIntoNSlices[T any](slice []T, n int) [][]T {
	var result [][]T

	// Calculate the base chunk size
	chunksize := len(slice) / n
	// Calculate the number of slices that will get one extra element due to division remainder
	overflow := len(slice) % n

	// Initialize indices
	start := 0

	for i := 0; i < n; i++ {
		// Calculate end index for the chunk
		end := start + chunksize
		if i < overflow {
			// Distribute the remaining elements (overflow) to the first few chunks
			end++
		}
		// Append the chunk to the result
		result = append(result, slice[start:end])

		// Update the start index for the next chunk
		start = end
	}

	return result
}
