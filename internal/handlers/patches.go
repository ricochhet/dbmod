package handlers

import (
	"github.com/ricochhet/dbmod/internal/api/patches"
	"github.com/ricochhet/dbmod/pkg/errorx"
)

func (c *Context) inventoryPatches() Queries {
	return Queries{
		{Name: "xpInfo", Query: c.xpInfoPatch},
		{Name: "shipDecorations", Query: c.shipDecorationsPatch},
	}
}

func (c *Context) xpInfoPatch(db *Database) (*Database, error) {
	weapons := c.Exports.Weapons
	if len(weapons) == 0 {
		return nil, errorx.WithFramef("Weapons data is %d bytes", len(weapons))
	}

	warframes := c.Exports.Warframes
	if len(warframes) == 0 {
		return nil, errorx.WithFramef("Warframes data is %d bytes", len(warframes))
	}

	sentinels := c.Exports.Sentinels
	if len(sentinels) == 0 {
		return nil, errorx.WithFramef("Sentinels data is %d bytes", len(sentinels))
	}

	inv, err := patches.ApplyXPInfo(weapons, warframes, sentinels, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) shipDecorationsPatch(db *Database) (*Database, error) {
	resources := c.Exports.Resources
	if len(resources) == 0 {
		return nil, errorx.WithFramef("Resources data is %d bytes", len(resources))
	}

	inv, err := patches.ApplyShipDecorations(resources, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}
