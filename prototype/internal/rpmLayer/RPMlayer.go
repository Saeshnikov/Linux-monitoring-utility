package rpmLayer

import (
	"bufio"
	"os"
	"os/exec"
	"time"
)

func RPMlayer(usedFiles []string, dirPath string) error {
	allPackages, err := findAllPackages()
	if err != nil {
		return err
	}
	usedPackages, err := findUsedPackages(usedFiles)
	if err != nil {
		return err
	}
	err_ := findUnusedPackages(allPackages, usedPackages, dirPath)
	if err_ != nil {
		return err_
	}
	return nil
}

func findAllPackages() (map[string]bool, error) {
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

func findUsedPackages(usedFiles []string) (map[string]bool, error) {
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

func findUnusedPackages(allPackages map[string]bool, usedPackages map[string]bool, dirPath string) error {
	filePath := "./out/"
	if dirPath != "" {
		// filePath = filepath.Join(dirPath, "/out/")
		filePath = dirPath + "/out/"
	}
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath + time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	for packageName := range usedPackages {
		if _, ok := allPackages[packageName]; ok {
			delete(allPackages, packageName)
		}
	}

	file.WriteString("Not used packages:\n")
	for packageName := range allPackages {
		file.WriteString(packageName + "\n")
	}
	return nil
}
