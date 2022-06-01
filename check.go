package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Checker struct {
	args Args
	proj Project
	cfg  Config
}

func (c *Checker) run() error {
	proj, err := newProject("")
	if err != nil {
		return err
	}

	c.proj = proj

	err = c.cfg.load(c.proj.ConfigFile)
	if err != nil {
		return err
	}

	if c.args.Verbose {
		fmt.Println("Build directory:", c.proj.BuildDir)
	}

	if err := c.prepare(c.args.Force); err != nil {
		return err
	}

	if err := os.Chdir(proj.BuildDir); err != nil {
		return err
	}

	var builder Builder
	if proj.ProjectType == CMakeProjectType {
		builder = CMakeBuilder{}
	}

	if err := builder.configure(*c, c.args.Verbose); err != nil {
		return err
	}

	if err := builder.build(*c, c.args.Verbose); err != nil {
		return err
	}

	if err := c.analyze(); err != nil {
		return err
	}

	if err := c.convert(); err != nil {
		return err
	}

	rep, err := saveReport(c.proj.TasksFile, c.proj.OutFile)
	if err != nil {
		return err
	}

	rep.print()

	return nil
}

func (checker Checker) prepare(clear bool) error {

	if clear {
		if err := os.RemoveAll(checker.proj.BuildDir); err != nil {
			return err
		}
	}

	os.MkdirAll(checker.proj.BuildDir, os.ModePerm)
	os.MkdirAll(checker.proj.DataDir, os.ModePerm)

	if err := checker.createTempConfig(); err != nil {
		return err
	}

	if err := checker.createTempRulesFile(); err != nil {
		return err
	}

	return nil

}

func (checker Checker) createTempConfig() error {
	f, err := os.Create(checker.proj.TmpCfgFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// In MacOS os.chdir(/var/folders/XXX) actually changes the directory to /private/var/folders/XXX
	absBuildDir, err := filepath.Abs(checker.proj.BuildDir)
	if err != nil {
		return err
	}

	lines := []string{
		fmt.Sprintf("analysis-mode=%d", checker.cfg.pvsMode()),
		fmt.Sprintf("sourcetree-root=%s", checker.proj.ProjectDir),

		fmt.Sprintf("exclude-path=%s", checker.proj.BuildDir),
		fmt.Sprintf("exclude-path=%s", absBuildDir),
		fmt.Sprintf("exclude-path=%s", "/opt"),
		fmt.Sprintf("exclude-path=%s", "/usr"),

		fmt.Sprintf("rules-config=%s", checker.proj.TmpRulesFile),
	}

	// Write lines ..............................
	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (checker Checker) createTempRulesFile() error {
	f, err := os.Create(checker.proj.TmpRulesFile)
	if err != nil {
		return err
	}
	defer f.Close()

	/*****************************
	Examples:
	//-V::112
	//-V:qCDebug:1044
	******************************/
	for _, c := range checker.cfg.Checks {
		s := fmt.Sprintf("//-%s", c)
		if _, err := f.WriteString(s + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (checker Checker) analyze() error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(checker.proj.BuildDir); err != nil {
		return err
	}

	return runWithProgress("Analyzing", checker.args.Verbose,
		"pvs-studio-analyzer",
		"analyze",
		"-j", fmt.Sprintf("%d", checker.cfg.PvsThreads),
		"--cfg", checker.proj.TmpCfgFile,
		"--incremental",
		//"--disableLicenseExpirationCheck",
		"-R", checker.proj.TmpRulesFile,
		"-o", checker.proj.LogFile,
	)
}

func (c Checker) convert() error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(c.proj.BuildDir); err != nil {
		return err
	}

	return runWithProgress("Report", c.args.Verbose,
		"plog-converter",
		"-a", c.cfg.pvsLevels(), // Specifies analyzer(s) and level(s) to be used for filtering
		"-s", c.proj.TmpCfgFile, // Path to PVS-Studio settings file.
		"-r", ".", // A path to the project directory.
		"-t", "tasklist", //  Render types for output.
		"-o", c.proj.TasksFile, // Output file.
		c.proj.LogFile,
	)
}
