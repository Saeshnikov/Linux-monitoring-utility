package tests

import (
	"bufio"
	bpfParsing "linux-monitoring-utility/internal/bpfParsing/bpftraceParsing"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {

	arr_res, err := bpfParsing.Parse("./data/bpfParsingIn.txt")
	if err != nil {
		t.Fatal(err.Error())
	}

	outputFile, err := os.Open("./data/bpfParsingOut.txt")
	if err != nil {
		t.Fatal(err.Error())
	}

	fileScanner := bufio.NewScanner(outputFile)
	var arr_test []string
	for fileScanner.Scan() {
		arr_test = append(arr_test, fileScanner.Text())
	}

	if reflect.DeepEqual(arr_res, arr_test) {
		return
	}
	t.Fatal("Test failed")
}
