package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docopt/docopt-go"
)

const AppVersion = "0.18"

const description = `
The pvscheck tool checks C & C++ projects using the pvs-studio utility..
`

const usage = `
Usage:
  pvscheck [options]

Options:
  -f --force            Rebuild a project from scratch
  --report              Create a report for a previously analysed project.
  --suppress 	        Suppressing all analyzer warnings
  --info                Show the information about project and exit
  -v --verbose          Show compiler output
  -h --help             Show this screen.
  --version             Show version of the program and PVS utilities.
  --config=CONFIGFILE   Use alternate configuration file 
`

const (
	RetArgsParseError = 1
	RetCommandError   = 2
	RetExecNotFound   = 3
)

// func processErrors() {
// 	r := recover()
// 	if r == nil {
// 		return
// 	}

// 	if execErr, ok := r.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
// 		fmt.Printf("%s command not found, it looks like the program is not installed\n", execErr.Name)
// 		os.Exit(RetExecNotFound)
// 	}

// 	fmt.Println(r)
// 	os.Exit(RetCommandError)
// }

// func run(command flags.Commander, args []string) error {
// 	defer processErrors()

// 	if err := command.Execute(args); err != nil {
// 		panic(err)
// 	}
// 	return nil
// }

type Args struct {
	Report   bool
	Suppress bool
	Info     bool
	Verbose  bool
	Version  bool
	Force    bool
	Config   string
}

func main() {
	args := Args{}

	{
		parser := docopt.Parser{}
		opts, _ := parser.ParseArgs(description+usage, os.Args[1:], "") // AppVersion)

		if err := opts.Bind(&args); err != nil {
			fmt.Println(err)
			parser.HelpHandler(err, usage)
			os.Exit(RetArgsParseError)
		}
	}

	if args.Version {
		showVersion()
		return
	}

	var err error
	{
		chk := Checker{args: args}
		err = chk.run()
	}

	//engine, err := newEngine(args.Input_bundle[0], args.Input_bundle[1], args.Output_bundle)
	//if err != nil {
	//fmt.Fprintln(os.Stderr, err)
	//os.Exit(RetArgsParseError)
	//}

	//engine.verbose = args.Verbose

	//err = engine.run()
	//if err != nil {
	//fmt.Fprintln(os.Stderr, "Error: ", err)
	//os.Exit(RetArgsParseError)
	//}

	// parser := flags.NewParser(&argv, flags.PassDoubleDash|flags.HelpFlag)
	// parser.CommandHandler = run

	// _, err := parser.Parse()

	// // We use "check" as default command .............
	// /*
	// 	if err != nil {
	// 		fmt.Println("......................")
	// 		fmt.Println(err)
	// 		fmt.Println("......................")
	// 		//if flagsErr, ok := err.(*flags.Error); ok {
	// 		//if flagsErr.Type == flags.ErrCommandRequired || flagsErr.Type == flags.ErrUnknownFlag {
	// 		args := os.Args
	// 		args[0] = "check"
	// 		_, err = parser.ParseArgs(args)
	// 		//}
	// 		//}
	// 	}
	// */

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(RetArgsParseError)
	}
}

func showVersion() error {
	fmt.Printf("pvscheck   %s\n", AppVersion)
	cmd := exec.Command("pvs-studio", "--version")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
