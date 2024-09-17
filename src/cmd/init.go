package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/pg"
	"snomed/src/trud"
)

const (
	usage               string = "%s init arguments:\n"
	defaultBinDirectory string = "./bin"
)

type InitCommand struct {
	fs *flag.FlagSet

	binPath  string
	category trud.Category
}

func NewInitCommand() *InitCommand {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	cc := &InitCommand{
		fs: fs,
	}

	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, filepath.Base(os.Args[0]))
		fs.PrintDefaults()
	}

	config := pg.Config
	fs.StringVar(
		&cc.binPath, "bin", defaultBinDirectory,
		"The temporary output directory for downloaded content",
	)
	fs.StringVar(&config.NhsTrudKey, "key", config.NhsTrudKey, "NHS Trud API key")
	fs.StringVar(&config.PostgresHost, "host", config.PostgresHost, "Postgres host")
	fs.UintVar(&config.PostgresPort, "port", config.PostgresPort, "Postgres port")
	fs.StringVar(&config.PostgresUsername, "uid", config.PostgresUsername, "Postgres username")
	fs.StringVar(&config.PostgresPassword, "pwd", config.PostgresPassword, "Postgres password")
	fs.StringVar(&config.PostgresDatabase, "db", config.PostgresDatabase, "Postgres database name")

	var cat string
	fs.StringVar(&cat, "cat", "SNOMED_ALL", "Desired SNOMED release categories")
	if succ, res := trud.ParseCategory(cat); succ {
		cc.category = res
	} else {
		cc.category = trud.SNOMED_NONE
	}

	return cc
}

func (c *InitCommand) Name() string {
	return c.fs.Name()
}

func (c *InitCommand) GetFlagSet() *flag.FlagSet {
	return c.fs
}

func (c *InitCommand) Init(ctx context.Context, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	_, err := pg.GetDb(ctx)
	if err != nil {
		return err
	}

	if err := trud.DownloadPackages(ctx, trud.SNOMED_ALL, pg.Config.NhsTrudKey, c.binPath); err != nil {
		return err
	}

	return nil
}

func (c *InitCommand) Run(ctx context.Context) error {

	return nil
}
