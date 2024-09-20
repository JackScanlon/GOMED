package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var cmds []Command

type Command interface {
	Init(context.Context, []string) error
	Run(context.Context) error
	Name() string

	GetFlagSet() *flag.FlagSet
}

func GenerateCommands() {
	cmds = []Command{
		NewBuildCommand(),
	}
}

func Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return nil
	}

	command := os.Args[1]
	for _, cmd := range cmds {
		if cmd.Name() != command {
			continue
		}

		if err := cmd.Init(ctx, os.Args[2:]); err != nil {
			return err
		}

		return cmd.Run(ctx)
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
