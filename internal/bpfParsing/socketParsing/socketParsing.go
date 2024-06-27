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
	PathsOfExecutableFiles [2]string
	WayOfInteraction       Interaction
}
//----------------------------------------------------------------

type SocketInfo struct {
	Ipc, Protocol string
}

func (s SocketInfo) String() string {
	return fmt.Sprintf("%s: %s", s.Ipc, s.Protocol)
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
			sock := socketData{pathOfExecutableFile: arr[0], syscallType: arr[1], protocol: arr[2], fileDescriptor: arr[3]}
			if !contains(sockArr, sock) {
				sockArr = append(sockArr, sock)
			}
		}
	}

	parsingArr := findConnection(sockArr)
	return parsingArr, nil
}

func contains(sockArr []socketData, sock socketData) bool {
	for _, s := range sockArr {
		if sock == s {
			return true
		}
	}
	return false
}

func findConnection(sockArr []socketData) []ParsingData {
	var parsingArr []ParsingData
	for i := 0; i < len(sockArr); i++ {
		for j := i + 1; j < len(sockArr); j++ {
			if sockArr[i].fileDescriptor == sockArr[j].fileDescriptor &&
				sockArr[i].syscallType != sockArr[j].syscallType &&
				sockArr[i].protocol == sockArr[j].protocol {
				sockInfo := SocketInfo{Ipc: "socket", Protocol: sockArr[i].protocol}
				parsingArr = append(parsingArr, ParsingData{[2]string{sockArr[i].pathOfExecutableFile, sockArr[j].pathOfExecutableFile}, sockInfo})
			}
		}
	}
	return parsingArr
}
