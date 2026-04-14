package handlers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ricochhet/dbmod/internal/config"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/fsx"
	"github.com/ricochhet/dbmod/pkg/mongox"
	"github.com/ricochhet/dbmod/pkg/timex"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	inventoriesCollection = "inventories"
	statsCollection       = "stats"
)

type Database struct {
	Inventory []byte
	Stats     []byte
}

type Context struct {
	Flags     *config.Config
	Connector *mongox.Client
	Exports   *config.Exports
}

func (c *Context) InventoryCheats(ctx context.Context, names []string) error {
	query := c.inventoryCheats()

	data, collection, err := c.readCollection(ctx, inventoriesCollection)
	if err != nil {
		return errorx.New("readCollection", err)
	}

	if err := c.writeBackup(inventoriesCollection, data); err != nil {
		return errorx.New("writeBackup", err)
	}

	result := query.Run(&Database{Inventory: data}, query.skip(names))

	return c.commit(ctx, collection, result.Inventory, "inventories_commit")
}

func (c *Context) StatsCheat(ctx context.Context, names []string) error {
	query := c.statsCheats()

	data, collection, err := c.readCollection(ctx, statsCollection)
	if err != nil {
		return errorx.New("readCollection", err)
	}

	if err := c.writeBackup(statsCollection, data); err != nil {
		return errorx.New("writeBackup", err)
	}

	result := query.Run(&Database{Stats: data}, query.skip(names))

	return c.commit(ctx, collection, result.Stats, "stats_commit")
}

func (c *Context) InventoryPatches(ctx context.Context, names []string) error {
	query := c.inventoryPatches()

	data, collection, err := c.readCollection(ctx, inventoriesCollection)
	if err != nil {
		return errorx.New("readCollection", err)
	}

	if err := c.writeBackup(inventoriesCollection, data); err != nil {
		return errorx.New("writeBackup", err)
	}

	result := query.Run(&Database{Inventory: data}, query.skip(names))

	return c.commit(ctx, collection, result.Inventory, "inventories_commit")
}

func (c *Context) commit(
	ctx context.Context,
	collection *mongo.Collection,
	data []byte,
	backup string,
) error {
	if !c.Flags.DryRun {
		if err := c.writeCollection(ctx, collection, data); err != nil {
			return errorx.WithFrame(err)
		}
	}

	if err := c.writeBackup(backup, data); err != nil {
		return errorx.WithFrame(err)
	}

	return nil
}

func (c *Context) readCollection(
	ctx context.Context,
	name string,
) ([]byte, *mongo.Collection, error) {
	data, collection, err := c.Connector.Get(ctx, c.Flags.Database, name)
	if err != nil {
		return nil, nil, errorx.WithFrame(err)
	}

	return data, collection, nil
}

func (c *Context) writeCollection(
	ctx context.Context,
	collection *mongo.Collection,
	data []byte,
) error {
	if err := c.Connector.Set(ctx, collection, data); err != nil {
		return errorx.WithFrame(err)
	}

	return nil
}

func (c *Context) writeBackup(name string, data []byte) error {
	if len(data) == 0 {
		return errorx.WithFrameStr("cannot write empty backup")
	}

	dest := filepath.Join(
		c.Flags.DBData,
		fmt.Sprintf("%s-%s", name, timex.TimeStamp()),
	)

	if err := fsx.Write(dest, data); err != nil {
		return errorx.WithFrame(err)
	}

	return nil
}
