package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/ricochhet/dbmod/internal/config"
	"github.com/ricochhet/dbmod/internal/handlers"
	"github.com/ricochhet/dbmod/pkg/cmdx"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/ricochhet/dbmod/pkg/mongox"
)

var (
	name      = "dbmod"
	buildDate string
	gitHash   string
	buildOn   string
)

func version() {
	logx.Infof("%s-%s\n", name, gitHash)
	logx.Infof("Build date: %s\n", buildDate)
	logx.Infof("Build on: %s\n", buildOn)
	os.Exit(0)
}

type Context struct {
	*handlers.Context
}

func main() {
	var err error

	logx.LogTime.Store(true)
	logx.MaxProcNameLength.Store(0)
	logx.New(name, 0, logx.ModeAllRelease)

	_ = cmdx.QuickEdit(false)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flags.Args = flag.Args()

	appCtx := Context{
		&handlers.Context{
			Flags: flags,
			WFData: &config.WFData{
				Exports: *readExports(flags.WFExportData),
				Custom:  *readCustom(flags.WFCustomData),
			},
		},
	}

	appCtx.Connector, err = mongox.New(flags.MongoURI)
	exitOnErr(err)

	defer func() { exitOnErr(appCtx.Connector.Disconnect(context.Background())) }()

	if flag.NArg() == 0 {
		exitOnErr(errorx.WithFrameStr("no args specified"))
	}

	if _, err := appCtx.commands(ctx); err != nil {
		logx.Errorf("Error running command: %v\n", err)
	}
}

func (c *Context) commands(ctx context.Context) (bool, error) {
	cmd := strings.ToLower(c.Flags.Args[0])
	rest := []string{"all"}

	if len(c.Flags.Args) > 1 {
		rest = c.Flags.Args[1:]
	}

	switch cmd {
	case "help":
		cmds.Usage()
	case "version":
		version()
	case "i", "inventory":
		switch c.Flags.Mode {
		case "c", "cheat":
			return true, c.InventoryCheats(ctx, rest)
		case "p", "patch":
			return true, c.InventoryPatches(ctx, rest)
		default:
			cmds.Usage()
		}
	case "s", "stats":
		switch c.Flags.Mode {
		case "c", "cheat":
			return true, c.StatsCheat(ctx, rest)
		default:
			cmds.Usage()
		}
	default:
		cmds.Usage()
	}

	return false, nil
}

func exitOnErr(err error) {
	if err != nil {
		logx.Errorf("%s: %v\n", name, err)
		os.Exit(1)
	}
}
