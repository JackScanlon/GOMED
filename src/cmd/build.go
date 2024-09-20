package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/codes"
	"snomed/src/pg"
	"snomed/src/trud"
)

const (
	usage               string = "%s build arguments:\n"
	defaultBinDirectory string = "./bin"
)

type BuildCommand struct {
	fs *flag.FlagSet

	binPath  string
	category trud.Category
	driver   *pg.Driver
	releases []*trud.Release
}

func NewBuildCommand() *BuildCommand {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	cc := &BuildCommand{
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

func (c *BuildCommand) Name() string {
	return c.fs.Name()
}

func (c *BuildCommand) GetFlagSet() *flag.FlagSet {
	return c.fs
}

func (c *BuildCommand) Init(ctx context.Context, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	driver, err := pg.GetDB(ctx)
	if err != nil {
		return err
	}
	c.driver = driver

	releases, err := trud.DownloadPackages(ctx, trud.SNOMED_ALL, pg.Config.NhsTrudKey, c.binPath)
	if err != nil {
		return err
	}
	c.releases = releases

	return nil
}

func (c *BuildCommand) Run(ctx context.Context) error {
	/*
		TODO:
			- [x] det. whether tables exist; create them if not - could also look at doing delta update?
			- [x] parse tab delimited text files -> upload to db
			- [x] create top-level code map
			- [ ] build simplified code map
			- [ ] build simplified ontology
			- [ ] add logger
	*/

	for _, release := range c.releases {
		if err := codes.BuildRelease(c.driver, release, c.binPath); err != nil {
			return err
		}
	}

	return nil
}
