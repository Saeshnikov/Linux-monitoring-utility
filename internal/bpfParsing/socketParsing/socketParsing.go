package socketParsing

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------
type Interaction interface {
	String() string
}

type ParsingData struct {
	PathOfExecutableFile1, PathOfExecutableFile2 string
	WayOfInteraction                             Interaction
}

//----------------------------------------------------------------

type SocketInfo struct {
	Ipc, Protocol, FileDescriptor string
}

func (s SocketInfo) String() string {
	return fmt.Sprintf("%s: %s, %s", s.Ipc, s.Protocol, s.FileDescriptor)
}

type socketData struct {
	pathOfExecutableFile, syscallType, protocol, fileDescriptor string
}

func Parse(fileName string) ([]ParsingData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var sockArr []socketData

	for fileScanner.Scan() {
		arr := strings.Fields(fileScanner.Text())
		if len(arr) == 4 {
			sockArr = append(sockArr, socketData{pathOfExecutableFile: arr[0], syscallType: arr[1], protocol: arr[2], fileDescriptor: arr[3]})
		}
	}

	parsingArr := findConnection(sockArr)
	return parsingArr, nil
}

func findConnection(sockArr []socketData) []ParsingData {
	var parsingArr []ParsingData
	for i := 0; i < len(sockArr); i++ {
		for j := 1; j < len(sockArr); j++ {
			if sockArr[i].fileDescriptor == sockArr[j].fileDescriptor &&
				sockArr[i].syscallType != sockArr[j].syscallType &&
				sockArr[i].protocol == sockArr[j].protocol {
				sockInfo := SocketInfo{Ipc: "by socket", Protocol: sockArr[i].protocol, FileDescriptor: sockArr[i].fileDescriptor}
				parsingArr = append(parsingArr, ParsingData{sockArr[i].pathOfExecutableFile, sockArr[j].pathOfExecutableFile, sockInfo})
				sockArr[j].fileDescriptor = "" //!!!!!!!!!!!!1
			}
		}
	}
	return parsingArr
}
