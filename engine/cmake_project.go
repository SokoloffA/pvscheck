package engine

import (
	"fmt"
)

func init() {
	loadProjectFuncs = append(loadProjectFuncs, newCMakeProject)
}

type CMakeProject struct {
	dir string
}

func (p CMakeProject) Directory() string {
	return p.dir
}

func (p CMakeProject) Type() string {
	return "cmake"
}

func newCMakeProject(dir string) Project {
	if !fileExists(dir + "/CMakeLists.txt") {
		return nil
	}

	return &CMakeProject{
		dir: dir,
	}
}

func (p CMakeProject) configure(e Engine, verbose bool) error {
	return runWithProgress("Cmake", verbose,
		"cmake",
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=On",
		"-B"+e.BuildDir,
		e.ProjectDir,
	)

}

func (p CMakeProject) build(e Engine, verbose bool) error {
	return runWithProgress("Build", verbose,
		"make",
		"-C"+e.BuildDir,
		"-j", fmt.Sprintf("%v", e.config.BuildThreads),
	)
}
