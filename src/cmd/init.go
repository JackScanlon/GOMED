package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/shared"
	"snomed/src/trud"
)

const (
	usage               string = "%s init arguments:\n"
	defaultBinDirectory string = "./bin"
)

type InitCommand struct {
	fs *flag.FlagSet

	binPath          string
	nhsTrudKey       string
	postgresHost     string
	postgresPort     uint
	postgresUsername string
	postgresPassword string
	postgresDatabase string
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

	fs.StringVar(
		&cc.binPath, "bin", defaultBinDirectory,
		"The temporary output directory for downloaded content",
	)

	config := shared.Config
	fs.StringVar(&cc.nhsTrudKey, "key", config.NhsTrudKey, "NHS Trud API key")
	fs.StringVar(&cc.postgresHost, "host", config.PostgresHost, "Postgres host")
	fs.UintVar(&cc.postgresPort, "port", config.PostgresPort, "Postgres port")
	fs.StringVar(&cc.postgresUsername, "uid", config.PostgresUsername, "Postgres username")
	fs.StringVar(&cc.postgresPassword, "pwd", config.PostgresPassword, "Postgres password")
	fs.StringVar(&cc.postgresDatabase, "db", config.PostgresDatabase, "Postgres database name")

	return cc
}

func (c *InitCommand) Name() string {
	return c.fs.Name()
}

func (c *InitCommand) GetFlagSet() *flag.FlagSet {
	return c.fs
}

func (c *InitCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if err := trud.DownloadPackages(trud.SNOMED_ALL, c.nhsTrudKey, c.binPath); err != nil {
		return err
	}

	return nil
}

func (c *InitCommand) Run() error {

	return nil
}
