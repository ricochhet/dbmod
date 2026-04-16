package strx

import (
	"sort"
	"strings"
)

func Replace(value string, rules []struct{ From, To string }) string {
	for _, r := range rules {
		if strings.Contains(value, r.From) {
			return strings.ReplaceAll(value, r.From, r.To)
		}
	}

	return value
}

func ReplaceMap(value string, rules map[string]string) string {
	for from, to := range rules {
		if strings.Contains(value, from) {
			return strings.ReplaceAll(value, from, to)
		}
	}

	return value
}

func ReplaceMapDeterministic(value string, rules map[string]string) string {
	keys := make([]string, 0, len(rules))
	for k := range rules {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, from := range keys {
		if strings.Contains(value, from) {
			return strings.ReplaceAll(value, from, rules[from])
		}
	}

	return value
}

// EndsWithAny checks if a string ends with any suffixes in the slice.
func EndsWithAny(str string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(str, s) {
			return true
		}
	}

	return false
}
