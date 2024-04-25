package lsofLayer

import (
	"bufio"
	"os/exec"
	"regexp"
)

func LsofExec(lsofBinPath string) ([]string, error) {
	cmd := exec.Command(lsofBinPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	outScanner := bufio.NewScanner(stdout)
	arr, err := LsofParsing(outScanner)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func LsofParsing(outScanner *bufio.Scanner) ([]string, error) {
	var arr []string
	r, err := regexp.Compile(`\s(/.*?)$`)
	if err != nil {
		return nil, err
	}
	for outScanner.Scan() {
		res := r.FindAllStringSubmatch(outScanner.Text(), -1)
		if res != nil {
			arr = append(arr, res[0][1])
		}
	}
	return arr, nil
}
