package cmd

import (
	"github.com/ZergsLaw/zerg-repo/cli/core"
	"github.com/ZergsLaw/zerg-repo/cli/fs"
	"github.com/urfave/cli/v2"
)

var NewMigrate = &cli.Command{
	Name:         "create",
	Usage:        "create migrate file",
	Description:  "Creates a new migration file with test data.",
	BashComplete: cli.DefaultAppComplete,
	Action:       newMigrateAction,
	Flags:        []cli.Flag{migrateName, dir},
}

func newMigrateAction(ctx *cli.Context) error {
	filesSystem := fs.New()

	c := core.New(filesSystem, nil)

	return c.NewMigrate(ctx.Context, ctx.String(dir.Name), ctx.String(migrateName.Name))
}