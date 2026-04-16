package jsonx

import (
	"fmt"

	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type FieldCriteria struct {
	Key       string
	Value     string
	KeyMode   MatchMode
	ValueMode MatchMode
}

type SearchCriteria struct {
	Key     string
	KeyMode MatchMode

	Value     string
	ValueMode MatchMode

	Fields []FieldCriteria

	MatchAll bool
}

func (m Match) String() string {
	return fmt.Sprintf("path=%q  key=%q  value=%s", m.Path, m.Key, m.Value.Raw)
}

func Search(json string, search SearchCriteria) []Match {
	root := gjson.Parse(json)

	var out []Match
	walk(root, "", "", search, &out)

	return out
}

func Edit(json, path string, value any) (string, error) {
	v, err := sjson.Set(json, path, value)
	if err != nil {
		return json, errorx.WithFrame(err)
	}

	return v, nil
}

func EditRaw(json, path, raw string) (string, error) {
	v, err := sjson.SetRaw(json, path, raw)
	if err != nil {
		return json, errorx.WithFrame(err)
	}

	return v, nil
}

func Delete(json, path string) (string, error) {
	v, err := sjson.Delete(json, path)
	if err != nil {
		return json, errorx.WithFrame(err)
	}

	return v, nil
}

func SearchAndEdit(
	json string,
	search SearchCriteria,
	editor func(m Match) any,
) (string, error) {
	matches := Search(json, search)
	if len(matches) == 0 {
		return json, errorx.WithFrameStr("no matches for criteria")
	}

	m := matches[0]

	return Edit(json, m.Path, editor(m))
}
