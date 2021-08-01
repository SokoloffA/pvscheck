package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type ReportInfo struct {
	outFile string
	notes   int
	warns   int
	errs    int
}

func (rep ReportInfo) Print() {
	fmt.Println("**************************")
	fmt.Printf("* The %s file was created\n", rep.outFile)
	fmt.Println("*")
	fmt.Println("* Ошибок:         ", rep.errs)
	fmt.Println("* Предупреждений: ", rep.warns)
	fmt.Println("* Уведомлений:    ", rep.notes)
}

func (e Engine) BuildReport() (ReportInfo, error) {
	res := ReportInfo{outFile: e.OutFile}

	in, err := os.Open(e.TasksFile)
	if err != nil {
		return res, err
	}
	defer in.Close()

	out, err := os.Create(e.OutFile)
	if err != nil {
		return res, err
	}
	defer out.Close()

	///////////////////
	reader := bufio.NewReader(in)
	var fileErr error
	for fileErr != io.EOF {
		var b []byte
		b, fileErr = reader.ReadBytes('\n')
		if fileErr != nil && fileErr != io.EOF {
			return res, fileErr
		}

		line := strings.TrimSpace(string(b))

		if line == "" {
			continue
		}

		if strings.Contains(line, "Help: The documentation for all analyzer") {
			continue
		}

		if strings.Contains(line, "\terr\t") {
			res.errs += 1
		}

		if strings.Contains(line, "\twarn\t") {
			res.warns += 1
		}

		if strings.Contains(line, "\tnote\t") {
			res.notes += 1
		}

		out.WriteString(line)
		out.WriteString("\n")
	}

	return res, nil
}
