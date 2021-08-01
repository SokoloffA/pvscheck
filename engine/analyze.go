package engine

import (
	"fmt"
	"os"
)

func (e Engine) Analyze(verbose bool) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(e.BuildDir); err != nil {
		return err
	}

	return runWithProgress("Analyzing", verbose,
		"pvs-studio-analyzer",
		"analyze",
		"-j", fmt.Sprintf("%d", e.config.PvsThreads),
		"--cfg", e.TmpCfgFile,
		//#"--incremental",
		"-o", e.LogFile,
	)
}

func (e Engine) Convert(verbose bool) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(e.BuildDir); err != nil {
		return err
	}

	return runWithProgress("Report", verbose,
		"plog-converter",
		"-a", e.config.PvsLevels(),
		"-s", e.TmpCfgFile,
		"--renderTypes=tasklist",
		"-o", e.TasksFile,
		e.LogFile,
	)
}

func (e Engine) FilterSuppressed(verbose bool) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(e.BuildDir); err != nil {
		return err
	}

	return runWithProgress("Filter suppressed messages", verbose,
		"pvs-studio-analyzer",
		"filter-suppressed",
		e.LogFile,
	)
}

func (e Engine) Suppress(verbose bool) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(e.BuildDir); err != nil {
		return err
	}

	return runWithProgress("Filter suppressed messages", verbose,
		"pvs-studio-analyzer",
		"filter-suppressed",
		e.LogFile,
	)
}
