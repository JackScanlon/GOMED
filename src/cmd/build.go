package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/codes"
	"snomed/src/pg"
	"snomed/src/templates"
	"snomed/src/trud"
)

const (
	usage               string = "%s build arguments:\n"
	defaultBinDirectory string = "./bin"
)

var (
	desiredCategory string = "SNOMED_ALL"
)

type BuildCommand struct {
	fs *flag.FlagSet

	binPath  string
	managed  bool
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
	fs.BoolVar(&cc.managed, "managed", false, "Specifies whether this application is managing the SNOMED tables")
	fs.StringVar(&config.NhsTrudKey, "key", config.NhsTrudKey, "NHS Trud API key")
	fs.StringVar(&config.PostgresHost, "host", config.PostgresHost, "Postgres host")
	fs.UintVar(&config.PostgresPort, "port", config.PostgresPort, "Postgres port")
	fs.StringVar(&config.PostgresUsername, "uid", config.PostgresUsername, "Postgres username")
	fs.StringVar(&config.PostgresPassword, "pwd", config.PostgresPassword, "Postgres password")
	fs.StringVar(&config.PostgresDatabase, "db", config.PostgresDatabase, "Postgres database name")
	fs.StringVar(&desiredCategory, "cat", desiredCategory, "Desired SNOMED release categories")

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

	if succ, res := trud.ParseCategory(desiredCategory); succ {
		c.category = res
	} else {
		c.category = trud.SNOMED_NONE
	}

	driver, err := pg.GetDB(ctx)
	if err != nil {
		return err
	}
	c.driver = driver

	if _, err := templates.InitContainer(ctx); err != nil {
		return err
	}

	releases, err := trud.DownloadPackages(ctx, c.category, pg.Config.NhsTrudKey, c.binPath)
	if err != nil {
		return err
	}
	c.releases = releases

	return nil
}

func (c *BuildCommand) Run(ctx context.Context) (err error) {
	/*
		TODO:
			- [x] det. whether tables exist; create them if not - could also look at doing delta update?
			- [x] parse tab delimited text files -> upload to db
			- [x] create top-level code map
			- [x] build simplified code map
			- [x] build simplified ontology
			- [ ] intermediate table cleanup
			- [ ] add logger
	*/

	var rebuilt bool = false
	for _, release := range c.releases {
		rebuilt, err = codes.BuildRelease(c.driver, release, c.binPath)
		if err != nil {
			return err
		}
	}

	if rebuilt {
		data := map[string]any{
			"managed": c.managed,
		}

		err = templates.
			GetContainer().
			Source(
				"concept:descriptionIdentifier",
				templates.WithData(data),
				templates.WithEcho(),
			).
			Exec()

		if err != nil {
			return err
		}

		err = templates.
			GetContainer().
			Source(
				"concept:simplifyCodelist",
				templates.WithData(data),
				templates.WithEcho(),
			).
			Exec()

		if err != nil {
			return err
		}

		err = templates.
			GetContainer().
			Source(
				"ontology:network",
				templates.WithData(data),
				templates.WithEcho(),
			).
			Exec()

		if err != nil {
			return err
		}
	}

	return nil
}
