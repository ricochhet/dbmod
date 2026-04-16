package main

import (
	"flag"

	"github.com/ricochhet/dbmod/internal/config"
	"github.com/ricochhet/dbmod/pkg/cmdx"
)

var (
	flags = &config.Config{}
	cmds  = cmdx.Info{
		{Usage: "dbmod help", Desc: "Show this help"},
		{Usage: "dbmod version", Desc: "Display dbmod version"},
		{Usage: "dbmod [MODE] inventory challenges", Desc: "Unlock all challenges (cheat)"},
		{Usage: "dbmod [MODE] inventory capturaScenes", Desc: "Unlock all captura scenes (cheat)"},
		{Usage: "dbmod [MODE] inventory flavourItems", Desc: "Unlock all flavor items (cheat)"},
		{Usage: "dbmod [MODE] inventory missions", Desc: "Unlock all missions (cheat)"},
		{
			Usage: "dbmod [MODE] inventory shipDecorations",
			Desc:  "Unlock all ship decorations (cheat, patch)",
		},
		{Usage: "dbmod [MODE] inventory weaponSkins", Desc: "Unlock all weapon skins (cheat)"},
		{Usage: "dbmod [MODE] inventory xpInfo", Desc: "Patch XP info (patch)"},
		{Usage: "dbmod [MODE] stats codexScans", Desc: "Unlock all codex scans (cheat)"},
		{Usage: "dbmod [MODE] stats enemyStats", Desc: "Modify enemy stats (cheat)"},
	}
)

//nolint:gochecknoinits // wontfix
func init() {
	registerFlags(flag.CommandLine, flags)
	flag.Parse()
}

func registerFlags(fs *flag.FlagSet, f *config.Config) {
	fs.BoolVar(&f.DryRun, "dry-run", false, "write changes locally without modifying the database")
	fs.StringVar(&f.MongoURI, "uri", "mongodb://localhost:27017", "MongoDB URI")
	fs.StringVar(&f.Database, "db-name", "openWF", "MongoDB database name")
	fs.StringVar(&f.DBData, "db-data", "dbdata", "backup and dry-run output directory")
	fs.StringVar(
		&f.WFExportData,
		"wf-export-data",
		"assets/exports",
		"Warframe public export directory",
	)
	fs.StringVar(
		&f.WFCustomData,
		"wf-custom-data",
		"assets/custom",
		"Warframe custom data directory",
	)
	fs.IntVar(&f.Index, "i", 0, "index of the document to modify")
	fs.StringVar(&f.Mode, "m", "cheat", "edit mode: c|cheat or p|patch")
	fs.BoolVar(&f.Debug, "debug", false, "enable debug output")
}
