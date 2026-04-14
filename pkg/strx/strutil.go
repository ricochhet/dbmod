package strx

import "strings"

// EndsWithAny checks if a string ends with any suffixes in the slice.
func EndsWithAny(str string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(str, s) {
			return true
		}
	}

	return false
}
