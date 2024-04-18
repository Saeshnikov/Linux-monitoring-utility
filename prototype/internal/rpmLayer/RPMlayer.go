package rpmLayer

import (
	"bufio"
	"os/exec"
)

func RPMlayer(usedFiles []string, dirPath string, outputMap *map[string]bool) error {
	usedPackages, err := FindUsedPackages(usedFiles)
	if err != nil {
		return err
	}
	err = FindUnusedPackages(usedPackages, dirPath, outputMap)
	if err != nil {
		return err
	}
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

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		outScanner := bufio.NewScanner(stdout)
		outScanner.Split(bufio.ScanWords)
		for outScanner.Scan() {
			usedPackages[outScanner.Text()] = true
		}
	}
	return usedPackages, nil
}

func FindUnusedPackages(usedPackages map[string]bool, dirPath string, outputMap *map[string]bool) {
	for packageName := range usedPackages {
		if _, ok := allPackages[packageName]; ok {
			delete(*outputMap, packageName)
		}
	}
}
