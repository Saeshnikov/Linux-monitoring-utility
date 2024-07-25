package tests

import (
	"reflect"
	"testing"

	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"linux-monitoring-utility/internal/bpfParsing/semaphoreParsing"
)

func TestParseSemaphores(t *testing.T) {
	ParsedDataExample := []parsingstruct.ParsingData{
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/IPC/semtest", "/home/anna/Desktop/bpftrace/IPC/semtest1"},
			WayOfInteraction: semaphoreParsing.SemaphoreInfo{Ipc: "by semaphore", Id: "2"}},
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/IPC/semtest", "/home/anna/Desktop/bpftrace/IPC/semtest2"},
			WayOfInteraction: semaphoreParsing.SemaphoreInfo{Ipc: "by semaphore", Id: "2"}},
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/IPC/semtest2", "/home/anna/Desktop/bpftrace/IPC/semtest1"},
			WayOfInteraction: semaphoreParsing.SemaphoreInfo{Ipc: "by semaphore", Id: "2"}},
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name        string
		args        args
		expectedArr []parsingstruct.ParsingData
	}{
		{"Basic test", args{"./data/semaphoreParsingTest.txt"}, ParsedDataExample},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedData, err := semaphoreParsing.Parse(tt.args.fileName)
			if err != nil {
				t.Error("Incorrectly opened file\n")
			} else {
				if !isEqualSemaphores(tt.expectedArr, parsedData) {
					t.Error("The resulting array did not match the expected one\n")
				}
			}
		})
	}
}

func isEqualSemaphores(parsedDataExp []parsingstruct.ParsingData, parsedDataGot []parsingstruct.ParsingData) bool {
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
