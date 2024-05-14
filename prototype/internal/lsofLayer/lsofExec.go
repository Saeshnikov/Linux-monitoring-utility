package lsofLayer

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"
)

var LsofBinPath string
var DirToIgnore []string

func LsofExec() ([]string, error) {
	cmd := exec.Command(LsofBinPath)
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
	// var arr map[string]bool
	arr := make(map[string]bool)
	var out []string
	r, err := regexp.Compile(`\s(/.*?)$`)
	if err != nil {
		return nil, err
	}
	for outScanner.Scan() {
		res := r.FindAllStringSubmatch(outScanner.Text(), -1)
		if res != nil {
			if len(res[0][1]) > 1 && len(strings.Fields(res[0][1])) == 1 && ignoreDir(res[0][1]) {
				arr[res[0][1]] = true
			}

		}
	}
	for k := range arr {
		out = append(out, k)
	}
	return out, nil
}

func ignoreDir(s string) bool {
	for _, dir := range DirToIgnore {
		if len(strings.Split(s, dir)) != 1 {
			return false
		}
	}
	return true
}
