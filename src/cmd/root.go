package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var cmds []Command

type Command interface {
	Init([]string) error
	Run() error
	Name() string

	GetFlagSet() *flag.FlagSet
}

func GenerateCommands() {
	cmds = []Command{
		NewInitCommand(),
	}
}

func Execute(args []string) error {
	if len(args) < 1 {
		return nil
	}

	command := os.Args[1]
	for _, cmd := range cmds {
		if cmd.Name() == command {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	if command == "-h" || command == "-help" {
		msg := fmt.Sprintf("Usage: %s [command] [options]\nAvailable commands:\n", filepath.Base(os.Args[0]))
		for _, cmd := range cmds {
			msg += fmt.Sprintf("  %s | arguments:\n", cmd.Name())

			cmd.GetFlagSet().VisitAll(func(f *flag.Flag) {
				msg += fmt.Sprintf("    -%s (default %q): %s\n", f.Name, f.DefValue, f.Usage)
			})
		}

		fmt.Print(msg)
		return nil
	}

	return fmt.Errorf("unknown command: %q", command)
}
