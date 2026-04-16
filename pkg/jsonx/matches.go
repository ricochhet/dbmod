package jsonx

import (
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

func walk(node gjson.Result, path, key string, search SearchCriteria, out *[]Match) {
	switch {
	case node.IsObject():
		if keyMatches(key, search.Key, search.KeyMode) {
			match := (len(search.Fields) > 0 && objectMatchesFields(node, search.Fields, search.MatchAll)) ||
				(len(search.Fields) == 0 && search.Value == "" && search.Key != "")

			if match {
				*out = append(*out, Match{Path: path, Key: key, Value: node})
			}
		}

		node.ForEach(func(k, v gjson.Result) bool {
			walk(v, joinPath(path, k.String()), k.String(), search, out)
			return true
		})

	case node.IsArray():
		node.ForEach(func(k, v gjson.Result) bool {
			walk(v, joinPath(path, k.String()), k.String(), search, out)
			return true
		})

	default:
		if len(search.Fields) == 0 &&
			keyMatches(key, search.Key, search.KeyMode) &&
			valueMatches(node.String(), search.Value, search.ValueMode) {
			*out = append(*out, Match{Path: path, Key: key, Value: node})
		}
	}
}
