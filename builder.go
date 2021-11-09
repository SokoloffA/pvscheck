package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Builder interface {
	configure(checker Checker, verbose bool) error
	build(checker Checker, verbose bool) error
}

func dumpCommandArgs(exe string, args ...string) string {
	res := []string{}

	res = append(res, exe)
	res = append(res, args...)

	for i, s := range res {
		if strings.ContainsAny(s, " \t\n") {
			res[i] = `"` + s + `"`
		} else if s == "" {
			res[i] = `""`
		}
	}

	return strings.Join(res, " ")
}

func runWithProgress(caption string, verbose bool, exe string, args ...string) error {
	if verbose {
		fmt.Println(dumpCommandArgs(exe, args...))

		proc := exec.Command(exe, args...)
		proc.Stderr = os.Stderr
		proc.Stdout = os.Stdout
		return proc.Run()
	}

	proc := exec.Command(exe, args...)
	proc.Stderr = os.Stderr

	stdout, _ := proc.StdoutPipe()
	if err := proc.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	showMarks := true
	marks := "|/-\\"
	mark := 0

	getPercent := func(line string) int {
		if line[0] != '[' {
			return -1
		}

		e := strings.Index(line, "%]")
		if e < 0 {
			return -1
		}

		res, err := strconv.Atoi(strings.TrimSpace(line[1:e]))
		if err != nil {
			return -1
		}

		return res
	}

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if line == "" {
			continue
		}

		percent := getPercent(line)
		if percent > -1 {
			showMarks = false
			fmt.Printf(" %s   %d%%\r", caption, percent)
		} else if showMarks {
			mark += 1
			fmt.Printf(" %s   %c\r", caption, marks[mark%len(marks)])
		}

	}

	if err := proc.Wait(); err != nil {
		for _, line := range lines {
			fmt.Println(line)
		}
		return err
	}

	return nil
}

type CMakeBuilder struct{}

func (builder CMakeBuilder) configure(checker Checker, verbose bool) error {
	return runWithProgress("Cmake", verbose,
		"cmake",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=On",
		"-B"+checker.proj.BuildDir,
		checker.proj.ProjectDir,
	)

}

func (builder CMakeBuilder) build(checker Checker, verbose bool) error {
	return runWithProgress("Build", verbose,
		"make",
		"-C"+checker.proj.BuildDir,
		"-j", fmt.Sprintf("%v", checker.cfg.BuildThreads),
	)
}
