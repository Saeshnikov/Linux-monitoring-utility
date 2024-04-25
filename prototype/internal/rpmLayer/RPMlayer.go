package rpmLayer

import (
	"bufio"
	"os/exec"
	"strings"
)

func RPMlayer(usedFiles []string, dirPath string, outputMap *map[string]bool) error {
	usedPackages, err := FindUsedPackages(usedFiles)
	if err != nil {
		return err
	}

	FindUnusedPackages(usedPackages, dirPath, outputMap)

	return nil
}

func FindAllPackages() (map[string]bool, error) {
	var allPackages = make(map[string]bool)
	cmd := exec.Command("/usr/bin/rpm", "-qa")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	outScanner := bufio.NewScanner(stdout)
	for outScanner.Scan() {
		allPackages[outScanner.Text()] = true
	}
	return allPackages, nil
}

func FindUsedPackages(usedFiles []string) (map[string]bool, error) {

	var usedPackages = make(map[string]bool)
	for _, fileName := range usedFiles {
		cmd := exec.Command("/usr/bin/rpm", "-qf", fileName)

		pipe, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(pipe)
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		line, err := reader.ReadString('\n')

		for err == nil {
			if len(strings.Fields(line)) == 1 {
				usedPackages[strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "\n", "")] = true
			}

			line, err = reader.ReadString('\n')
		}
		cmd.Wait()
	}
	return usedPackages, nil
}

func FindUnusedPackages(usedPackages map[string]bool, dirPath string, outputMap *map[string]bool) {
	for packageName := range usedPackages {
		delete(*outputMap, packageName)
	}
}
