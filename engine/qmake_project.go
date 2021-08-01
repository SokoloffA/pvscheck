package engine

import (
	"fmt"
	"path/filepath"
)

func init() {
	loadProjectFuncs = append(loadProjectFuncs, newQMakeProject)
}

type QMakeProject struct {
	dir string
}

func (p QMakeProject) Directory() string {
	return p.dir
}

func (p QMakeProject) Type() string {
	return "qmake"
}

func newQMakeProject(dir string) Project {
	f, err := filepath.Glob(dir + "/*.pro")
	if err != nil {
		return nil
	}

	if len(f) < 1 {
		return nil
	}

	return &QMakeProject{
		dir: dir,
	}
}

func (p QMakeProject) configure(engine Engine, verbose bool) error {
	return runWithProgress("Qmake", verbose,
		"qmake",
		"CONFIG+=debug",
		p.dir,
	)
}

func (p QMakeProject) build(engine Engine, verbose bool) error {
	return runWithProgress("Build", verbose,
		"pvs-studio-analyzer",
		"trace",
		"--",
		"make",
		"-j", fmt.Sprintf("%d", engine.config.BuildThreads),
	)
}
