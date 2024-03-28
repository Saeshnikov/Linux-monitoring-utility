package lsofLayer

import (
	"bufio"
	"os/exec"
	"regexp"
)

func LsofExec() ([]string, error) {
	cmd := exec.Command("/usr/bin/lsof")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	outScanner := bufio.NewScanner(stdout)
	var arr []string
	r, err := regexp.Compile(`(/.*?)$`)
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
