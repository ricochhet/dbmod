package jsonx

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/fsx"
	"github.com/tidwall/gjson"
	"github.com/tidwall/jsonc"
	"github.com/tidwall/sjson"
)

// Marshal v of type T.
func Marshal[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal data into type T, and store it in store(v).
func Unmarshal[T any](data []byte, store func(T)) error {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return errorx.WithFrame(err)
	}

	store(v)

	return nil
}

// ReadAndUnmarshal parses a JSON file from the specified path.
func ReadAndUnmarshal[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.New("os.ReadFile", err)
	}

	var t T
	if err := json.Unmarshal(jsonc.ToJSON(data), &t); err != nil {
		return nil, errorx.New("json.Unmarshal", err)
	}

	return &t, nil
}

// MarshalAndWrite marshales the data to the specified output file.
func MarshalAndWrite[T any](path string, data T) ([]byte, error) {
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}

	return b, fsx.Write(path, b)
}

func ArrayElementFieldValues(json []byte, path string, index int) ([]gjson.Result, error) {
	t, err := ArrayElementField(json, path, index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return t.Array(), nil
}

func ArrayElementField(json []byte, path string, index int) (gjson.Result, error) {
	array := gjson.ParseBytes(json).Array()
	if index < 0 || index >= len(array) {
		return gjson.Result{}, errorx.WithFramef("invalid index: %d", index)
	}

	return array[index].Get(path), nil
}

func ArrayElement(json []byte, index int) (gjson.Result, error) {
	array := gjson.ParseBytes(json).Array()
	if index < 0 || index >= len(array) {
		return gjson.Result{}, errorx.WithFramef("invalid index: %d", index)
	}

	return array[index], nil
}

func SetArrayElementRaw(data, elem []byte, index int) ([]byte, error) {
	return sjson.SetRawBytes(data, strconv.Itoa(index), elem)
}

func SetArrayElementFieldArray(
	json []byte,
	path string,
	value []string,
	index int,
) ([]byte, error) {
	return sjson.SetRawBytes(
		json,
		fmt.Sprintf("%d.%s", index, path),
		[]byte("["+strings.Join(value, ",")+"]"),
	)
}

func SetArrayElementFieldRaw(json []byte, path, value string, index int) ([]byte, error) {
	return sjson.SetRawBytes(json, fmt.Sprintf("%d.%s", index, path), []byte(value))
}

func SetArrayElementField[T any](json, path string, value T, index int) (string, error) {
	return sjson.Set(json, fmt.Sprintf("%d.%s", index, path), value)
}

func Join(path ...string) string {
	if len(path) == 0 {
		return ""
	}

	var s strings.Builder
	s.WriteString(path[0])

	for _, p := range path[1:] {
		if p == "" {
			continue
		}

		s.WriteString("." + p)
	}

	return s.String()
}
