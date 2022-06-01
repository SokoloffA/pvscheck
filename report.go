package main

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

func saveReport(inFile, outFile string) (ReportInfo, error) {
	res := ReportInfo{outFile: outFile}

	in, err := os.Open(inFile)
	if err != nil {
		return res, err
	}
	defer in.Close()

	out, err := os.Create(outFile)
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

		// Workaround for https://github.com/viva64/pvs-studio-cmake-examples/issues/18
		if strings.Contains(line, "V1042") {
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

func (rep ReportInfo) print() {
	fmt.Println("**************************")
	fmt.Printf("* The %s file was created\n", rep.outFile)
	fmt.Println("*")
	fmt.Println("* Ошибок:         ", rep.errs)
	fmt.Println("* Предупреждений: ", rep.warns)
	fmt.Println("* Уведомлений:    ", rep.notes)
}
