package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectType int

const (
	UnknownProjectType ProjectType = iota
	CMakeProjectType
	QMakeProjectType
)

type Project struct {
	ProjectDir   string
	ProjectType  ProjectType
	BuildDir     string
	DataDir      string
	TmpCfgFile   string
	TmpRulesFile string
	LogFile      string
	TasksFile    string
	OutFile      string
}

func newProject(dir string) (Project, error) {
	var err error
	res := Project{}

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return res, err
		}

		res.ProjectDir, res.ProjectType, err = findSrcProject(dir)
	} else {
		res.ProjectDir, res.ProjectType, err = loadSrcProject(dir)
	}

	if err != nil {
		return res, err
	}

	//if res.config, err = NewConfig(); err != nil {
	//return res, err
	//}

	dir1 := filepath.Base(filepath.Dir(res.ProjectDir))
	dir2 := filepath.Base(res.ProjectDir)
	hash := md5.Sum([]byte(res.ProjectDir))
	res.BuildDir = fmt.Sprintf("%spvs-%s-%s_[%x]", os.TempDir(), dir1, dir2, hash)

	res.DataDir = res.BuildDir + "/.PVS-Studio"
	res.TmpCfgFile = res.DataDir + "/pvscheck.cfg"
	res.TmpRulesFile = res.DataDir + "/.pvsconfig"
	res.LogFile = res.DataDir + "/PVS.log"
	res.TasksFile = res.DataDir + "/PVS.tasks"
	res.OutFile = res.ProjectDir + "/PVS.tasks"

	return res, nil
}

func loadSrcProject(dir string) (string, ProjectType, error) {
	// Cmake project ............................
	isCmakeProject := func(dir string) bool {
		return fileExists(dir + "/CMakeLists.txt")
	}

	// QMake project ............................
	isQmakeProject := func(dir string) bool {
		f, err := filepath.Glob(dir + "/*.pro")
		return err == nil && len(f) > 0
	}

	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", UnknownProjectType, err
	}

	if isCmakeProject(dir) {
		return dir, CMakeProjectType, nil
	}

	if isQmakeProject(dir) {
		return dir, QMakeProjectType, nil
	}

	return "", UnknownProjectType, fmt.Errorf("%s directory has an unknown project type", dir)
}

func findSrcProject(dir string) (string, ProjectType, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", UnknownProjectType, err
	}

	type Info struct {
		d string
		t ProjectType
	}
	found := []Info{}
	path := strings.Split(dir, "/")

	for len(path) > 0 {
		d, t, _ := loadSrcProject(strings.Join(path, "/"))
		if t != UnknownProjectType {
			found = append(found, Info{d, t})
		}

		path = path[:len(path)-1]
	}

	if len(found) > 0 {
		res := found[len(found)-1]
		return res.d, res.t, nil
	}

	return "", UnknownProjectType, fmt.Errorf("the project directory was not found or has an unknown type")
}
