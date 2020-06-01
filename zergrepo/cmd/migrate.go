package cmd

import (
	"errors"

	zergrepo "github.com/ZergsLaw/zerg-repo"
	"github.com/ZergsLaw/zerg-repo/zergrepo/core"
	"github.com/ZergsLaw/zerg-repo/zergrepo/fs"
	"github.com/ZergsLaw/zerg-repo/zergrepo/migrater"
	"github.com/urfave/cli/v2"
)

var Migrate = &cli.Command{
	Name:         "migrate",
	Usage:        "exec migrate",
	Description:  "Migrate the DB To the most recent version available",
	BashComplete: cli.DefaultAppComplete,
	Before:       beforeMigrateAction,
	After:        afterMigrateAction,
	Action:       migrateAction,
	Flags:        dbFlags,
}

func beforeMigrateAction(ctx *cli.Context) error {
	log.Info("starting migration...")
	return nil
}

func afterMigrateAction(ctx *cli.Context) error {
	log.Info("finished")
	return nil
}

const (
	pg = `postgres`
	ms = `mysql`
)

func migrateAction(ctx *cli.Context) error {
	cmd, err := parse(ctx.String(Operation.Name))
	if err != nil {
		return err
	}

	dbDriver := ctx.String(Driver.Name)
	if dbDriver != pg && dbDriver != ms {
		return ErrUnknownDriver
	}

	conn, err := zergrepo.ConnectByCfg(ctx.Context, dbDriver, zergrepo.Config{
		Host:     ctx.String(Host.Name),
		Port:     ctx.Int(Port.Name),
		User:     ctx.String(User.Name),
		Password: ctx.String(Pass.Name),
		DBName:   ctx.String(Name.Name),
		SSLMode:  zergrepo.DBSSLMode,
	})
	if err != nil {
		return err
	}

	metric := zergrepo.MustMetric("zergrepo", "migrater")
	r := zergrepo.New(conn, log, metric, zergrepo.NewMapper())
	defer r.Close()

	m := migrater.New(r)
	filesSystem := fs.New()

	c := core.New(filesSystem, m)

	return c.Migrate(ctx.Context, ctx.String(Dir.Name), core.Config{
		Cmd: cmd,
		To:  ctx.Uint(To.Name),
	})
}

var (
	ErrUnknownOperation = errors.New("unknown Operation")
	ErrUnknownDriver    = errors.New("unknown Driver")
)

func parse(op string) (cmd core.MigrateCmd, err error) {
	switch op {
	case core.Up.String():
		cmd = core.Up
	case core.UpTo.String():
		cmd = core.UpTo
	case core.UpOne.String():
		cmd = core.UpOne
	case core.Down.String():
		cmd = core.Down
	case core.DownTo.String():
		cmd = core.DownTo
	case core.Reset.String():
		cmd = core.Reset
	default:
		return 0, ErrUnknownOperation
	}

	return cmd, nil
}
