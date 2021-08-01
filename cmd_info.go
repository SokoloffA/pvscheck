package main

import (
	"fmt"

	"github.com/SokoloffA/pvscheck/engine"
)

type InfoArgs struct {
	CommandArgs
}

func (args *InfoArgs) Execute(_ []string) error {

	eng, err := engine.New(args.Pos.Directory)
	if err != nil {
		return err
	}

	fmt.Println("  Project directory:", eng.ProjectDir)
	fmt.Println("  Project type:     ", eng.ProjectType())
	fmt.Println("")
	fmt.Println("  Build directory:  ", eng.BuildDir)
	fmt.Println("  Tasks file:       ", eng.TasksFile)
	return nil
}
