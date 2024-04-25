package tests

import (
	"bufio"
	"fmt"
	"linux-monitoring-utility/internal/lsofLayer"
	"os"
	"reflect"
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

	for fileScanner.Scan() {

		fmt.Println(fileScanner.Text())
	}

	outputFile, err := os.Open("./data/lsofParsingOut.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	fileScanner = bufio.NewScanner(outputFile)
	var arr_test []string
	for fileScanner.Scan() {
		arr_test = append(arr_test, fileScanner.Text())

	}
	if reflect.DeepEqual(arr_res, arr_test) {
		return
	}
	t.Fatal("Test failed")
}
