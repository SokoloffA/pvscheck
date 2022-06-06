package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/docopt/docopt-go"
)

const AppVersion = "0.21"

const description = `
The pvscheck tool checks C & C++ projects using the pvs-studio utility..
`

const usage = `
Usage:
  pvscheck [options]
  pvscheck [options] check
  pvscheck [options] init
  pvscheck [options] report
  pvscheck [options] info

These are commands used in various situations:
  check			Check the project in the current directory
  init			Init default config file
  report		Create a report for a previously analysed project.
  info			Show the information about project and exit
  suppress 	    Suppressing all analyzer warnings

Options:
  -f --force            	Rebuild a project from scratch
  -v --verbose          	Show compiler output
  -h --help             	Show this screen.
  -V --version             	Show version of the program and PVS utilities.
  -c --config=CONFIGFILE   	Use alternate configuration file
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
	// Commands ............
	Check    bool
	Init     bool
	Report   bool
	Info     bool
	Suppress bool

	// Options .............
	Verbose bool
	Version bool
	Force   bool
	Config  string
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

	handleError := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(RetCommandError)
		}

		os.Exit(0)
	}

	if args.Version {
		handleError(runShowVersion())
		return
	}

	if args.Init {
		handleError(runInitConfig(args))
		return
	}

	if args.Info {
		handleError(runInfo(args))
		return
	}

	chk := Checker{args: args}
	handleError(chk.run())

}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func copyFile(src string, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

func runShowVersion() error {
	fmt.Printf("pvscheck   %s\n", AppVersion)
	cmd := exec.Command("pvs-studio", "--version")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func runInitConfig(args Args) error {
	fileName := args.Config
	if len(fileName) == 0 {
		fileName = DefaultConfigFile
	}

	if fileExists(fileName) {
		return fmt.Errorf("file '%s' already exists", fileName)
	}

	err := ioutil.WriteFile(fileName, []byte(defaultConfigYml), 0644)
	if err != nil {
		return fmt.Errorf("can't writre '%s' file: %s", fileName, err)
	}

	return nil
}

func runInfo(args Args) error {
	proj, err := newProject("")
	if err != nil {
		return err
	}

	projectType := ""
	switch proj.ProjectType {
	case UnknownProjectType:
		projectType = "Unknown"
	case CMakeProjectType:
		projectType = "CMake"
	case QMakeProjectType:
		projectType = "QMake"
	}

	fmt.Println("*****************************")
	fmt.Println("Project dir: ", proj.ProjectDir)
	fmt.Println("Project type:", projectType)
	fmt.Println("Output file: ", proj.OutFile)
	fmt.Println("Config file: ", proj.ConfigFile)
	fmt.Println("")
	fmt.Println("Build dir:   ", proj.BuildDir)
	if args.Verbose {
		fmt.Println("...........................")
		fmt.Println("Data dir:    ", proj.DataDir)
		fmt.Println("Tmp coonfig: ", proj.TmpCfgFile)
		fmt.Println("Tmp rules:   ", proj.TmpRulesFile)
		fmt.Println("Log file:    ", proj.LogFile)
		fmt.Println("Tasks file:  ", proj.TasksFile)
	}
	fmt.Println("*****************************")

	return nil
}
