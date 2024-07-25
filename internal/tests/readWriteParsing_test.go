package tests

import (
	"reflect"
	"testing"

	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"linux-monitoring-utility/internal/bpfParsing/readWriteParsing"
)

func TestParseReadWrite(t *testing.T) {
	ParsedDataExample := []parsingstruct.ParsingData{
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/pipe/fifoclient.out", "/home/anna/Desktop/bpftrace/pipe/fifoserver.out"},
			WayOfInteraction: readWriteParsing.ReadWriteInfo{Ipc: "by reading/writing",
				PathOfOpenedFile: "MYFIFO", ReadBytes: "4", WrittenBytes: "14"}},
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name        string
		args        args
		expectedArr []parsingstruct.ParsingData
	}{
		{"Basic test", args{"./data/testReadWrite.txt"}, ParsedDataExample},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedData, err := readWriteParsing.Parse(tt.args.fileName)
			if err != nil {
				t.Error("Incorrectly opened file\n")
			} else {
				if !isEqualReadWrite(tt.expectedArr, parsedData) {
					t.Error("The resulting array did not match the expected one\n")
				}
			}
		})
	}
}

func isEqualReadWrite(parsedDataExp []parsingstruct.ParsingData, parsedDataGot []parsingstruct.ParsingData) bool {
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
