package main

import (
	"fmt"
	"os"
	"os/exec"
)

type VersionArgs struct {
	Force bool `short:"f" long:"force" description:"rebuild a project from scratch."`
}

func (args *VersionArgs) Execute(_ []string) error {
	fmt.Printf("pvscheck   %s\n", AppVersion)
	cmd := exec.Command("pvs-studio", "--version")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
