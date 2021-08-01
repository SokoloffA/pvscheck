package main

import "github.com/SokoloffA/pvscheck/engine"

type ReportArgs struct {
	CommandArgs
}

func (args ReportArgs) Execute(_ []string) (err error) {

	e, err := engine.New(args.Pos.Directory)
	if err != nil {
		return err
	}

	rep, err := e.BuildReport()
	if err != nil {
		return err
	}

	rep.Print()

	return nil
}
