package main

import (
	"fmt"

	"github.com/SokoloffA/pvscheck/engine"
)

type CheckArgs struct {
	CommandArgs

	Force bool `short:"f" long:"force" description:"rebuild a project from scratch"`
}

func (args CheckArgs) Execute(_ []string) error {

	e, err := engine.New(args.Pos.Directory)
	if err != nil {
		return err
	}

	if argv.Verbose {
		fmt.Println("Build directory:", e.BuildDir)
	}

	if err := e.Prepare(args.Force); err != nil {
		return err
	}

	if err := e.Configure(argv.Verbose); err != nil {
		return err
	}

	if err := e.Build(argv.Verbose); err != nil {
		return err
	}

	if err := e.Analyze(argv.Verbose); err != nil {
		return err
	}

	if err := e.Convert(argv.Verbose); err != nil {
		return err
	}

	if err := e.FilterSuppressed(argv.Verbose); err != nil {
		return err
	}

	rep, err := e.BuildReport()
	if err != nil {
		return err
	}

	rep.Print()
	return nil
}
