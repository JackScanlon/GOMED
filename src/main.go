package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/cmd"
	"snomed/src/shared"
)

func init() {
	if _, err := shared.RegisterEnvironment(); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), errors.New(shared.EnvironmentUsage()))
		os.Exit(1)
	}

	cmd.GenerateCommands()
}

func main() {
	fmt.Printf("%s/%s\n", filepath.Base(os.Args[0]), shared.GetVersion())

	ctx := context.Background()
	if err := cmd.Execute(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		os.Exit(1)
	}
}
