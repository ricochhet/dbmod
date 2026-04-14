package mongox

import (
	"bytes"

	"github.com/ricochhet/dbmod/pkg/errorx"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MarshalExtJSONArray converts a JSON document into a JSON array.
func MarshalExtJSONArray(docs []any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')

	for i, doc := range docs {
		if i > 0 {
			buf.WriteByte(',')
		}

		b, err := bson.MarshalExtJSON(doc, false, false)
		if err != nil {
			return nil, errorx.WithFrame(err)
		}

		buf.Write(b)
	}

	buf.WriteByte(']')

	return buf.Bytes(), nil
}
