package main

import (
	"path/filepath"

	"github.com/ricochhet/dbmod/internal/config"
	"github.com/ricochhet/dbmod/pkg/fsx"
	"github.com/ricochhet/dbmod/pkg/logx"
)

func readExports(path string) *config.Exports {
	join := func(name string) string { return filepath.Join(path, name) }

	return &config.Exports{
		Achievements: maybeRead(join("ExportAchievements.json")),
		Codex:        maybeRead(join("ExportCodex.json")),
		Customs:      maybeRead(join("ExportCustoms.json")),
		Enemies:      maybeRead(join("ExportEnemies.json")),
		Flavor:       maybeRead(join("ExportFlavour.json")),
		Regions:      maybeRead(join("ExportRegions.json")),
		Resources:    maybeRead(join("ExportResources.json")),
		Virtuals:     maybeRead(join("ExportVirtuals.json")),
		Weapons:      maybeRead(join("ExportWeapons.json")),
		Warframes:    maybeRead(join("ExportWarframes.json")),
		WarframesU41: maybeRead(join("ExportWarframes_41.1.0.json")),
		WarframesU42: maybeRead(join("ExportWarframes_42.0.6.json")),
		Sentinels:    maybeRead(join("ExportSentinels.json")),
		AllScans:     maybeRead(join("allScans.json")),
	}
}

func maybeRead(path string) []byte {
	data, err := fsx.Read(path)
	if err != nil {
		logx.Errorf("failed to read %s: %v\n", path, err)
		return nil
	}

	return data
}
