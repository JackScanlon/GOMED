package cmd

import (
	"context"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"snomed/src/pg"
	"snomed/src/templates"
)

var (
	desiredCleanupItems string = ""

	availableCleanupItems = map[string]any{
		"releases": false,
		"codelist": false,
		"ontology": false,
		"all":      false,
	}
)

type CleanCommand struct {
	fs *flag.FlagSet

	driver *pg.Driver
}

func NewCleanCommand() *CleanCommand {
	fs := flag.NewFlagSet("clean", flag.ContinueOnError)
	cc := &CleanCommand{
		fs: fs,
	}

	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageFmt, filepath.Base(os.Args[0]), fs.Name())
		fs.PrintDefaults()
	}

	fs.StringVar(
		&desiredCleanupItems,
		"items",
		desiredCleanupItems,
		"Specify items to cleanup, either 'all' or from: releases, codelist, ontology (delimited by comma)",
	)

	config := pg.Config
	fs.StringVar(&config.PostgresHost, "host", config.PostgresHost, "Postgres host")
	fs.UintVar(&config.PostgresPort, "port", config.PostgresPort, "Postgres port")
	fs.StringVar(&config.PostgresUsername, "uid", config.PostgresUsername, "Postgres username")
	fs.StringVar(&config.PostgresPassword, "pwd", config.PostgresPassword, "Postgres password")
	fs.StringVar(&config.PostgresDatabase, "db", config.PostgresDatabase, "Postgres database name")

	return cc
}

func (c *CleanCommand) Name() string {
	return c.fs.Name()
}

func (c *CleanCommand) GetFlagSet() *flag.FlagSet {
	return c.fs
}

func (c *CleanCommand) Init(ctx context.Context, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	driver, err := pg.GetDB(ctx)
	if err != nil {
		return err
	}
	c.driver = driver

	if _, err := templates.InitContainer(ctx); err != nil {
		return err
	}

	return nil
}

func (c *CleanCommand) Run(ctx context.Context) (err error) {
	var data map[string]any = maps.Clone(availableCleanupItems)
	valid := false
	items := strings.Split(desiredCleanupItems, ",")

	for _, item := range items {
		item := strings.TrimSpace(strings.ToLower(item))
		if _, ok := availableCleanupItems[item]; !ok {
			continue
		}
		valid = true
		data[item] = true

		if item == "all" {
			break
		}
	}

	if valid {
		err = templates.
			GetContainer().
			Source(
				"operations:cleanup",
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
