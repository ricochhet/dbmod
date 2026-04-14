package jsonx

import (
	"encoding/json"
	"fmt"
	"os"
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

// ResultAsArray gets a json result as a slice of gjson.Result.
func ResultAsArray(data []byte, name string, index int) ([]gjson.Result, error) {
	t, err := Result(data, name, index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return t.Array(), nil
}

// Result gets a result slice at the specified index.
func Result(data []byte, path string, index int) (gjson.Result, error) {
	array := gjson.ParseBytes(data).Array()
	if index < 0 || index >= len(array) {
		return gjson.Result{}, errorx.WithFramef("invalid index: %d", index)
	}

	target := array[index]

	return target.Get(path), nil
}

// SetSliceInRawBytes sets a slice to the input json bytes at the specified index.
func SetSliceInRawBytes(input []byte, path string, elems []string, index int) ([]byte, error) {
	return sjson.SetRawBytes(
		input,
		fmt.Sprintf("%d.%s", index, path),
		[]byte("["+strings.Join(elems, ",")+"]"),
	)
}

// SetFieldInRawBytes sets a field to the input json bytes at the specified index.
func SetFieldInRawBytes(input []byte, path, elem string, index int) ([]byte, error) {
	return sjson.SetRawBytes(input, fmt.Sprintf("%d.%s", index, path), []byte(elem))
}

// SetFieldInRawBytes sets a field to the input json string at the specified index.
func SetFieldInBytes[T any](input, path string, elem T, index int) (string, error) {
	return sjson.Set(input, fmt.Sprintf("%d.%s", index, path), elem)
}
