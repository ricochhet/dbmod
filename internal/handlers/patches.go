package handlers

import (
	"github.com/ricochhet/dbmod/internal/api/patches"
	"github.com/ricochhet/dbmod/pkg/errorx"
)

func (c *Context) inventoryPatches() Queries {
	return Queries{
		{Name: "xpInfo", Query: c.xpInfoPatch},
		{Name: "shipDecorations", Query: c.shipDecorationsPatch},
		{Name: "abilityPaths", Query: c.abilityPathsPatch},
	}
}

func (c *Context) xpInfoPatch(db *Database) (*Database, error) {
	weapons := c.WFData.Exports.Weapons
	if len(weapons) == 0 {
		return nil, errorx.WithFramef("Weapons data is %d bytes", len(weapons))
	}

	warframes := c.WFData.Exports.Warframes
	if len(warframes) == 0 {
		return nil, errorx.WithFramef("Warframes data is %d bytes", len(warframes))
	}

	sentinels := c.WFData.Exports.Sentinels
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
	resources := c.WFData.Exports.Resources
	if len(resources) == 0 {
		return nil, errorx.WithFramef("Resources data is %d bytes", len(resources))
	}

	inv, err := patches.ApplyShipDecorations(resources, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) abilityPathsPatch(db *Database) (*Database, error) {
	warframesU41 := c.WFData.Custom.WarframesU41
	if len(warframesU41) == 0 {
		return nil, errorx.WithFramef("WarframesU41 data is %d bytes", len(warframesU41))
	}

	warframesU42 := c.WFData.Custom.WarframesU42
	if len(warframesU42) == 0 {
		return nil, errorx.WithFramef("WarframesU42 data is %d bytes", len(warframesU42))
	}

	inv, err := patches.ApplyU42AbilityPaths(
		warframesU41,
		warframesU42,
		db.Inventory,
		c.Flags.Index,
	)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}
