package handlers

import (
	"github.com/ricochhet/dbmod/internal/api/cheats"
	"github.com/ricochhet/dbmod/pkg/errorx"
)

func (c *Context) inventoryCheats() Queries {
	return Queries{
		{Name: "accolades", Query: c.accoladesCheat},
		{Name: "challenges", Query: c.challengesCheat},
		{Name: "capturaScenes", Query: c.capturaScenesCheat},
		{Name: "flavourItems", Query: c.flavourItemsCheat},
		{Name: "missions", Query: c.missionsCheat},
		{Name: "shipDecorations", Query: c.shipDecorationsCheat},
		{Name: "weaponSkins", Query: c.weaponSkinsCheat},
	}
}

func (c *Context) statsCheats() Queries {
	return Queries{
		{Name: "codexScans", Query: c.codexScansCheat},
		{Name: "enemyStats", Query: c.enemyStatsCheat},
	}
}

func (c *Context) accoladesCheat(db *Database) (*Database, error) {
	accolades := cheats.Accolades{
		Staff:     true,
		Founder:   4,
		Guide:     2,
		Moderator: true,
		Partner:   true,
		Heirloom:  true,
		Counselor: true,
	}

	inv, err := accolades.Apply(db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) challengesCheat(db *Database) (*Database, error) {
	achievements := c.WFData.Exports.Achievements
	if len(achievements) == 0 {
		return nil, errorx.WithFramef("Achievements data is %d bytes", len(achievements))
	}

	inv, err := cheats.ApplyChallenges(achievements, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) capturaScenesCheat(db *Database) (*Database, error) {
	resources := c.WFData.Exports.Resources
	if len(resources) == 0 {
		return nil, errorx.WithFramef("Resources data is %d bytes", len(resources))
	}

	virtuals := c.WFData.Exports.Virtuals
	if len(virtuals) == 0 {
		return nil, errorx.WithFramef("Virtuals data is %d bytes", len(virtuals))
	}

	inv, err := cheats.ApplyCapturaScenes(resources, virtuals, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) flavourItemsCheat(db *Database) (*Database, error) {
	flavor := c.WFData.Exports.Flavor
	if len(flavor) == 0 {
		return nil, errorx.WithFramef("Flavor data is %d bytes", len(flavor))
	}

	inv, err := cheats.ApplyFlavourItems(flavor, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) missionsCheat(db *Database) (*Database, error) {
	regions := c.WFData.Exports.Regions
	if len(regions) == 0 {
		return nil, errorx.WithFramef("Regions data is %d bytes", len(regions))
	}

	inv, err := cheats.ApplyMissions(regions, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) shipDecorationsCheat(db *Database) (*Database, error) {
	shipDecorations := cheats.ShipDecorations{MaxCount: 999}

	resources := c.WFData.Exports.Resources
	if len(resources) == 0 {
		return nil, errorx.WithFramef("Resources data is %d bytes", len(resources))
	}

	inv, err := shipDecorations.Apply(resources, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) weaponSkinsCheat(db *Database) (*Database, error) {
	customs := c.WFData.Exports.Customs
	if len(customs) == 0 {
		return nil, errorx.WithFramef("Customs data is %d bytes", len(customs))
	}

	inv, err := cheats.ApplyWeaponSkins(customs, db.Inventory, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: inv, Stats: db.Stats}, nil
}

func (c *Context) enemyStatsCheat(db *Database) (*Database, error) {
	enemyStats := cheats.EnemyStats{Kills: 25, Assists: 5, Headshots: 10}

	enemies := c.WFData.Exports.Enemies
	if len(enemies) == 0 {
		return nil, errorx.WithFramef("Enemies data is %d bytes", len(enemies))
	}

	stats, err := enemyStats.Apply(enemies, db.Stats, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: db.Inventory, Stats: stats}, nil
}

func (c *Context) codexScansCheat(db *Database) (*Database, error) {
	codexScans := cheats.CodexScans{MaxScans: 99}

	allScans := c.WFData.Custom.AllScans
	if len(allScans) == 0 {
		return nil, errorx.WithFramef("AllScans data is %d bytes", len(allScans))
	}

	codex := c.WFData.Exports.Codex
	if len(codex) == 0 {
		return nil, errorx.WithFramef("Codex data is %d bytes", len(codex))
	}

	enemies := c.WFData.Exports.Enemies
	if len(enemies) == 0 {
		return nil, errorx.WithFramef("Enemies data is %d bytes", len(enemies))
	}

	stats, err := codexScans.Apply(allScans, codex, enemies, db.Stats, c.Flags.Index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	return &Database{Inventory: db.Inventory, Stats: stats}, nil
}
