package engine

import (
	"fmt"
	"path/filepath"
	"strings"
)

var loadProjectFuncs []func(dir string) Project

type ProjectType int

const (
	UnknownProjectType ProjectType = iota
	CMakeProjectType
	QMakeProjectType
)

type Project interface {
	Directory() string
	Type() string
	configure(e Engine, verbose bool) error
	build(e Engine, verbose bool) error
}

func loadProject(dir string) (Project, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range loadProjectFuncs {
		if res := f(dir); res != nil {
			return res, nil
		}
	}

	return nil, fmt.Errorf("%s directory has an unknown project type", dir)
}

func findProject(dir string) (Project, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	res := []Project{}
	path := strings.Split(dir, "/")

	for len(path) > 0 {
		if proj, _ := loadProject(strings.Join(path, "/")); proj != nil {
			res = append(res, proj)
		}

		path = path[:len(path)-1]
	}

	if len(res) > 0 {
		return res[len(res)-1], nil
	}

	return nil, fmt.Errorf("the project directory was not found or has an unknown type")
}
