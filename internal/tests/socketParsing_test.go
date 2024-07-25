package tests

import (
	"reflect"
	"testing"

	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"linux-monitoring-utility/internal/bpfParsing/socketParsing"
)

func TestParseSockets(t *testing.T) {
	ParsedDataExample := []parsingstruct.ParsingData{
		{PathsOfExecutableFiles: [2]string{"/snapshot/usr/bin/VBoxClient", "/snapshot/usr/bin/VBoxClient/a"},
			WayOfInteraction: socketParsing.SocketInfo{Ipc: "by socket", Protocol: "UNIX"}},
		{PathsOfExecutableFiles: [2]string{"/home/anna/Desktop/bpftrace/socket/socket_server", "/home/anna/Desktop/bpftrace/socket/socket_client"},
			WayOfInteraction: socketParsing.SocketInfo{Ipc: "by socket", Protocol: "INET"}},
		{PathsOfExecutableFiles: [2]string{"/snapshot/usr/bin/VBoxClient/c", "/snapshot/usr/bin/VBoxClient/a"},
			WayOfInteraction: socketParsing.SocketInfo{Ipc: "by socket", Protocol: "UNIX"}},
		{PathsOfExecutableFiles: [2]string{"/snapshot/usr/bin/VBoxClient/b", "/snapshot/usr/bin/VBoxClient/a"},
			WayOfInteraction: socketParsing.SocketInfo{Ipc: "by socket", Protocol: "UNIX"}},
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name        string
		args        args
		expectedArr []parsingstruct.ParsingData
	}{
		{"Basic test", args{"./data/socketParsingTest.txt"}, ParsedDataExample},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedData, err := socketParsing.Parse(tt.args.fileName)
			if err != nil {
				t.Error("Incorrectly opened file\n")
			} else {
				if !isEqualSockets(tt.expectedArr, parsedData) {
					t.Error("The resulting array did not match the expected one\n")
				}
			}
		})
	}
}

func isEqualSockets(parsedDataExp []parsingstruct.ParsingData, parsedDataGot []parsingstruct.ParsingData) bool {
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
