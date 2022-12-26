package utils

// return -1 if not found
func SliceFindIndex(condition uint32, data []uint32) int {

	for idx, d := range data {
		if d == condition {
			return idx
		}
	}

	return -1
}
