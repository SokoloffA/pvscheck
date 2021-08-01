package engine

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

type Engine struct {
	ProjectDir   string
	BuildDir     string
	DataDir      string
	TmpCfgFile   string
	TmpRulesFile string
	LogFile      string
	TasksFile    string
	OutFile      string
	config       Config
	project      Project
}

func New(dir string) (Engine, error) {
	var err error
	res := Engine{}

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return res, err
		}

		res.project, err = findProject(dir)
	} else {
		res.project, err = loadProject(dir)
	}

	if err != nil {
		return res, err
	}

	if res.config, err = NewConfig(); err != nil {
		return res, err
	}

	res.ProjectDir = res.project.Directory()
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

func (e Engine) ProjectType() string {
	if e.project != nil {
		return e.project.Type()
	}
	return "Unknown"
}

func (e Engine) Prepare(clear bool) error {

	if err := os.RemoveAll(e.BuildDir); err != nil {
		return err
	}

	os.MkdirAll(e.BuildDir, os.ModePerm)
	os.MkdirAll(e.DataDir, os.ModePerm)

	e.createTempConfig()

	return nil
}

func (e Engine) createTempConfig() error {
	f, err := os.Create(e.TmpCfgFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// In MacOS os.chdir(/var/folders/XXX) actually changes the directory to /private/var/folders/XXX
	absBuildDir, err := filepath.Abs(e.BuildDir)
	if err != nil {
		return err
	}

	lines := []string{
		fmt.Sprintf("analysis-mode=%d", e.config.PvsMode()),
		fmt.Sprintf("sourcetree-root=%s", e.ProjectDir),

		fmt.Sprintf("exclude-path=%s", e.BuildDir),
		fmt.Sprintf("exclude-path=%s", absBuildDir),
		fmt.Sprintf("exclude-path=%s", "/opt"),
		fmt.Sprintf("exclude-path=%s", "/usr"),
	}

	rulesFiles, err := e.searchRulesFiles()
	if err != nil {
		return err
	}

	if len(rulesFiles) > 0 {
		if err := e.createTempRulesFile(rulesFiles); err != nil {
			return err
		}
		lines = append(lines, fmt.Sprintf("rules-config=%s", e.TmpRulesFile))
	}

	// Write lines ..............................
	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (e Engine) searchRulesFiles() ([]string, error) {
	do := func(dir string) []string {
		res := []string{}

		path := strings.Split(dir, "/")
		for len(path) > 0 {

			if fileExists(dir + "/.pvsconfig") {
				res = append(res, dir+"/.pvsconfig")
			}

			path = path[:len(path)-1]
		}

		return res
	}

	dir, err := filepath.Abs(e.ProjectDir)
	if err != nil {
		return []string{}, err
	}

	res := do(dir)

	if fileExists(os.Getenv("HOME") + "/.config/PVS-Studio/.pvsconfig") {
		res = append(res, os.Getenv("HOME")+"/.config/PVS-Studio/.pvsconfig")
	}

	// Reverse results
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return res, nil
}

func (e Engine) createTempRulesFile(rulesFiles []string) error {
	out, err := os.Create(e.TmpRulesFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range rulesFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		out.Write(data)
		out.WriteString("\n")
	}

	return nil
}

func (e Engine) Configure(verbose bool) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(e.BuildDir); err != nil {
		return err
	}

	return e.project.configure(e, verbose)
}

func (e Engine) Build(verbose bool) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cur)

	if err := os.Chdir(e.BuildDir); err != nil {
		return err
	}

	return e.project.build(e, verbose)
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
