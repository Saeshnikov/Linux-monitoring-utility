package tests

import (
	"bufio"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestRPMLayer(t *testing.T) {
	var allPackages = make(map[string]bool)
	var outputMap = make(map[string]bool)
	usedPackages, err := readFile("./usedPackages.txt")

	cmd := exec.Command("/usr/bin/rpm", "-qa")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err.Error())
	}
	outScanner := bufio.NewScanner(stdout)
	for outScanner.Scan() {
		allPackages[outScanner.Text()] = true
		outputMap[outScanner.Text()] = true
	}

	for _, packageName := range usedPackages {
		if _, ok := allPackages[packageName]; ok {
			delete(allPackages, packageName)
		}
	}
	usedFiles, err := readFile("./usedFiles.txt")
	rpmLayer.RPMlayer(usedFiles, "", outputMap)
	eq := reflect.DeepEqual(allPackages, outputMap)
	if eq {
		t.Fatal("Test passed")
	} else {
		t.Fatal("Test failed")
	}
}

func readFile(filePath string) ([]string, error) {
	var arr []string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		arr = append(arr, fileScanner.Text())
	}
	return arr, nil
}
