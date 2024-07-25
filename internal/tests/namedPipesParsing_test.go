package tests

import (
	"reflect"
	"testing"

	"linux-monitoring-utility/internal/bpfParsing/namedPipesParsing"
	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
)

func TestParsePipes(t *testing.T) {
	ParsedDataExample := []parsingstruct.ParsingData{
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/pipe/fifo2.out", "/home/anna/Desktop/bpftrace/pipe/fifo1.out"},
			WayOfInteraction: namedPipesParsing.NamedPipesInfo{Ipc: "by named pipes", Name: "myfifo2"}},
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name        string
		args        args
		expectedArr []parsingstruct.ParsingData
	}{
		{"Basic test", args{"./data/testNamedPipes.txt"}, ParsedDataExample},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedData, err := namedPipesParsing.Parse(tt.args.fileName)
			if err != nil {
				t.Error("Incorrectly opened file\n")
			} else {
				if !isEqualPipes(tt.expectedArr, parsedData) {
					t.Error("The resulting array did not match the expected one\n")
				}
			}
		})
	}
}

func isEqualPipes(parsedDataExp []parsingstruct.ParsingData, parsedDataGot []parsingstruct.ParsingData) bool {
	if len(parsedDataExp) != len(parsedDataGot) {
		return false
	} else {
		var counter int
		for pExp := range parsedDataExp {
			for pGot := range parsedDataGot {
				if reflect.DeepEqual(pExp, pGot) {
					counter++
				}
			}
		}
		return counter == len(parsedDataExp)
	}
}
