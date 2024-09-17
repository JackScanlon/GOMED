package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"snomed/src/cmd"
	"snomed/src/shared"
)

func init() {
	if err := shared.RegisterEnvironment(); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), errors.New(shared.EnvironmentUsage()))
		os.Exit(1)
	}

	cmd.GenerateCommands()
}

func main() {
	fmt.Printf("%s/%s\n", filepath.Base(os.Args[0]), shared.GetVersion())

	if err := cmd.Execute(os.Args[1:]); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		os.Exit(1)
	}
}
