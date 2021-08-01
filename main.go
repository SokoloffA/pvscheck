package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jessevdk/go-flags"
)

const AppVersion = "0.17"

const (
	RetArgsParseError = 1
	RetCommandError   = 2
	RetExecNotFound   = 3
)

var argv Args

type Args struct {
	Check      CheckArgs    `command:"check" description:"check the project, this command is executed if no other is specified"`
	ReportArgs ReportArgs   `command:"report" description:"create a report for a previously analysed project."`
	Suppress   SuppressArgs `command:"suppress" description:"suppressing all analyzer warnings"`
	Info       InfoArgs     `command:"info" description:"show the information"`
	Version    VersionArgs  `command:"version" alias:"ver" description:"show the version of the program and PVS utilities"`

	Verbose bool `short:"v" long:"verbose" description:"show compiler output"`
}

type CommandArgs struct {
	Pos struct {
		Directory string `positional-arg-name:"PROJECT DIRECTORY"`
	} `positional-args:"yes"`
}

func processErrors() {
	r := recover()
	if r == nil {
		return
	}

	if execErr, ok := r.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
		fmt.Printf("%s command not found, it looks like the program is not installed\n", execErr.Name)
		os.Exit(RetExecNotFound)
	}

	fmt.Println(r)
	os.Exit(RetCommandError)
}

func run(command flags.Commander, args []string) error {
	defer processErrors()

	if err := command.Execute(args); err != nil {
		panic(err)
	}
	return nil
}

func main() {

	parser := flags.NewParser(&argv, flags.PassDoubleDash|flags.HelpFlag)
	parser.CommandHandler = run

	_, err := parser.Parse()

	// We use "check" as default command .............
	/*
		if err != nil {
			fmt.Println("......................")
			fmt.Println(err)
			fmt.Println("......................")
			//if flagsErr, ok := err.(*flags.Error); ok {
			//if flagsErr.Type == flags.ErrCommandRequired || flagsErr.Type == flags.ErrUnknownFlag {
			args := os.Args
			args[0] = "check"
			_, err = parser.ParseArgs(args)
			//}
			//}
		}
	*/

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(RetArgsParseError)
	}
}
