package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/pg"
	"snomed/src/shared"
	"snomed/src/templates"
)

type CopyCommand struct {
	fs *flag.FlagSet

	driver     *pg.Driver
	filePath   string
	delimiter  string
	targetName string
}

func NewCopyCommand() *CopyCommand {
	fs := flag.NewFlagSet("copy", flag.ContinueOnError)
	cc := &CopyCommand{
		fs: fs,
	}

	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageFmt, filepath.Base(os.Args[0]), fs.Name())
		fs.PrintDefaults()
	}

	fs.StringVar(
		&cc.filePath, "file", "",
		"Path to content file",
	)

	fs.StringVar(
		&cc.delimiter, "delimiter", "E'\\t'",
		"Content file field delimiter",
	)

	fs.StringVar(
		&cc.targetName, "target", "",
		"Table target name",
	)

	config := shared.Config
	fs.StringVar(&config.NhsTrudKey, "key", config.NhsTrudKey, "NHS Trud API key")
	fs.StringVar(&config.PostgresHost, "host", config.PostgresHost, "Postgres host")
	fs.UintVar(&config.PostgresPort, "port", config.PostgresPort, "Postgres port")
	fs.StringVar(&config.PostgresUsername, "uid", config.PostgresUsername, "Postgres username")
	fs.StringVar(&config.PostgresPassword, "pwd", config.PostgresPassword, "Postgres password")
	fs.StringVar(&config.PostgresDatabase, "db", config.PostgresDatabase, "Postgres database name")

	return cc
}

func (c *CopyCommand) Name() string {
	return c.fs.Name()
}

func (c *CopyCommand) GetFlagSet() *flag.FlagSet {
	return c.fs
}

func (c *CopyCommand) Init(ctx context.Context, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if c.filePath == "" {
		return errors.New("expected filePath as non-empty string describing a path to a file")
	} else if _, err := os.Stat(c.filePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file<%s> does not exist", c.filePath)
	} else {
		absPath, err := filepath.Abs(c.filePath)
		if err != nil {
			return err
		}

		c.filePath = absPath
	}

	if c.delimiter == "" {
		return errors.New("expected delimiter as non-empty string describing the field delimiter")
	}

	if c.targetName == "" {
		return errors.New("expected targetName as non-empty string describing the target table name")
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

func (c *CopyCommand) Run(ctx context.Context) (err error) {
	data := map[string]any{
		"filePath":   c.filePath,
		"delimiter":  c.delimiter,
		"targetName": c.targetName,
	}

	fmt.Printf("%s | %s | %s\n", c.filePath, c.delimiter, c.targetName)

	err = templates.
		GetContainer().
		Source(
			"copy:file",
			templates.WithData(data),
			templates.WithEcho(),
		).
		Exec()

	if err != nil {
		return err
	}

	return nil
}
