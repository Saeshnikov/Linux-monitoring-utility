package tests

import (
	"bufio"
	"linux-monitoring-utility/internal/lsofLayer"
	"os"
	"testing"
)

func TestLsofParsing(t *testing.T) {
	inputFile, err := os.Open("./data/lsofParsingIn.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	fileScanner := bufio.NewScanner(inputFile)
	arr_res, err := lsofLayer.LsofParsing(fileScanner)
	if err != nil {
		t.Fatal(err.Error())
	}
	var m = make(map[string]bool)
	for _, s := range arr_res {
		m[s] = true
	}

	outputFile, err := os.Open("./data/lsofParsingOut.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	fileScanner = bufio.NewScanner(outputFile)

	for fileScanner.Scan() {
		if !m[fileScanner.Text()] {
			t.Fatal("Test failed")
		}

	}

}
