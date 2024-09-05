package tests

import (
	"reflect"
	"testing"

	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"linux-monitoring-utility/internal/bpfParsing/sharedMemParsing"
)

func TestParseSharedMem(t *testing.T) {
	ParsedDataExample := []parsingstruct.ParsingData{
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/IPC/shmwrite", "/home/anna/Desktop/bpftrace/socket/socket_client"},
			WayOfInteraction: sharedMemParsing.SharedMemInfo{Ipc: "by semaphore", Key: "ffffffff", Id: "3080239", Type: "systemV"}},
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name        string
		args        args
		expectedArr []parsingstruct.ParsingData
	}{
		{"Basic test", args{"./data/sharedMemParsingTest.txt"}, ParsedDataExample},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedData, err := sharedMemParsing.Parse(tt.args.fileName)
			if err != nil {
				t.Error("Incorrectly opened file\n")
			} else {
				if !isEqualSharedMem(tt.expectedArr, parsedData) {
					t.Error("The resulting array did not match the expected one\n")
				}
			}
		})
	}
}

func isEqualSharedMem(parsedDataExp []parsingstruct.ParsingData, parsedDataGot []parsingstruct.ParsingData) bool {
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
