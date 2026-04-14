package mongox

import (
	"context"

	"github.com/ricochhet/dbmod/pkg/errorx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Client struct {
	*mongo.Client
}

func New(uri string) (*Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Client{client}, nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	if err := c.Client.Disconnect(ctx); err != nil {
		return errorx.WithFrame(err)
	}

	return nil
}

func (c *Client) Get(
	ctx context.Context,
	database, collection string,
) ([]byte, *mongo.Collection, error) {
	col := c.Database(database).Collection(collection)

	data, err := c.fetchAll(ctx, col)
	if err != nil {
		return nil, nil, errorx.WithFrame(err)
	}

	return data, col, nil
}

func (c *Client) Set(ctx context.Context, collection *mongo.Collection, data []byte) error {
	if err := c.replaceAll(ctx, collection, data); err != nil {
		return errorx.WithFrame(err)
	}

	return nil
}

func (c *Client) fetchAll(ctx context.Context, collection *mongo.Collection) ([]byte, error) {
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errorx.New("collection.Find", err)
	}
	defer cursor.Close(ctx)

	var results []bson.D
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errorx.New("cursor.All", err)
	}

	docs := make([]any, len(results))
	for i := range results {
		docs[i] = results[i]
	}

	return MarshalExtJSONArray(docs)
}

func (c *Client) replaceAll(ctx context.Context, collection *mongo.Collection, data []byte) error {
	if err := collection.Drop(ctx); err != nil {
		return errorx.New("collection.Drop", err)
	}

	var rawDocs []bson.Raw
	if err := bson.UnmarshalExtJSON(data, false, &rawDocs); err == nil {
		docs := make([]any, len(rawDocs))
		for i, raw := range rawDocs {
			docs[i] = raw
		}

		_, err := collection.InsertMany(ctx, docs)

		return errorx.New("collection.InsertMany", err)
	}

	var single bson.Raw
	if err := bson.UnmarshalExtJSON(data, false, &single); err != nil {
		return errorx.New("bson.UnmarshalExtJSON", err)
	}

	_, err := collection.InsertOne(ctx, single)

	return errorx.New("collection.InsertOne", err)
}
