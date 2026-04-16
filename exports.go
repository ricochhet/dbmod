package main

import (
	"path/filepath"

	"github.com/ricochhet/dbmod/internal/config"
	"github.com/ricochhet/dbmod/pkg/fsx"
	"github.com/ricochhet/dbmod/pkg/logx"
)

func readExports(path string) *config.Exports {
	return &config.Exports{
		Achievements: maybeRead(filepath.Join(path, "ExportAchievements.json")),
		Codex:        maybeRead(filepath.Join(path, "ExportCodex.json")),
		Customs:      maybeRead(filepath.Join(path, "ExportCustoms.json")),
		Enemies:      maybeRead(filepath.Join(path, "ExportEnemies.json")),
		Flavor:       maybeRead(filepath.Join(path, "ExportFlavour.json")),
		Regions:      maybeRead(filepath.Join(path, "ExportRegions.json")),
		Resources:    maybeRead(filepath.Join(path, "ExportResources.json")),
		Virtuals:     maybeRead(filepath.Join(path, "ExportVirtuals.json")),
		Weapons:      maybeRead(filepath.Join(path, "ExportWeapons.json")),
		Warframes:    maybeRead(filepath.Join(path, "ExportWarframes.json")),
		Sentinels:    maybeRead(filepath.Join(path, "ExportSentinels.json")),
	}
}

func readCustom(path string) *config.Custom {
	return &config.Custom{
		WarframesU41: maybeRead(filepath.Join(path, "ExportWarframes_41.1.0.json")),
		WarframesU42: maybeRead(filepath.Join(path, "ExportWarframes_42.0.6.json")),
		AllScans:     maybeRead(filepath.Join(path, "allScans.json")),
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
