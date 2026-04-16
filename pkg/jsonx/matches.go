package jsonx

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type MatchMode int

const (
	MatchExact MatchMode = iota
	MatchPartial
)

type Match struct {
	Path  string
	Key   string
	Value gjson.Result
}

func keyMatches(actual, pattern string, mode MatchMode) bool {
	if pattern == "" {
		return true
	}

	return strMatch(actual, pattern, mode)
}

func valueMatches(actual, pattern string, mode MatchMode) bool {
	if pattern == "" {
		return true
	}

	return strMatch(actual, pattern, mode)
}

func strMatch(actual, pattern string, mode MatchMode) bool {
	if mode == MatchPartial {
		return strings.Contains(strings.ToLower(actual), strings.ToLower(pattern))
	}

	return actual == pattern
}

func objectMatchesFields(node gjson.Result, fields []FieldCriteria, all bool) bool {
	for _, fc := range fields {
		matched := false

		node.ForEach(func(k, v gjson.Result) bool {
			if keyMatches(k.String(), fc.Key, fc.KeyMode) &&
				valueMatches(v.String(), fc.Value, fc.ValueMode) {
				matched = true
				return false
			}

			return true
		})

		if all && !matched {
			return false
		}

		if !all && matched {
			return true
		}
	}

	return all
}

func searchForMatches(
	node gjson.Result,
	path, key string,
	search SearchCriteria,
	matches *[]Match,
) {
	switch {
	case node.IsObject():
		if shouldAppendObject(node, key, search) {
			fmt.Println("Path: " + path)
			*matches = append(*matches, Match{Path: path, Key: key, Value: node})
		}

		node.ForEach(func(k, v gjson.Result) bool {
			searchForMatches(v, Join(path, k.String()), k.String(), search, matches)
			return true
		})

	case node.IsArray():
		node.ForEach(func(k, v gjson.Result) bool {
			searchForMatches(v, Join(path, k.String()), k.String(), search, matches)
			return true
		})

	default:
		if len(search.Fields) == 0 &&
			keyMatches(key, search.Key, search.KeyMode) &&
			valueMatches(node.String(), search.Value, search.ValueMode) {
			*matches = append(*matches, Match{Path: path, Key: key, Value: node})
		}
	}
}

func shouldAppendObject(node gjson.Result, key string, search SearchCriteria) bool {
	if !keyMatches(key, search.Key, search.KeyMode) {
		return false
	}

	if len(search.Fields) > 0 {
		return objectMatchesFields(node, search.Fields, search.MatchAll)
	}

	return search.Value == "" && search.Key != ""
}
