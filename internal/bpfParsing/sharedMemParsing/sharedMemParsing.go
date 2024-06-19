package sharedMemParsing

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

type SharedMemInfo struct {
	Ipc, Key, Id, Type string
}

func (s SharedMemInfo) String() string {
	return fmt.Sprintf("%s: %s, %s, %s", s.Ipc, s.Key, s.Id, s.Type)
}

type sharedMemData struct {
	pathOfExecutableFile, key, id, typeIpc string
}

func Parse(fileName string) ([]ParsingData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var memArr []sharedMemData

	for fileScanner.Scan() {
		arr := strings.Fields(fileScanner.Text())
		if len(arr) == 4 { // sections of posix are not being taken into account now!
			mem := sharedMemData{pathOfExecutableFile: arr[0], key: arr[1], id: arr[2], typeIpc: arr[3]}
			if !contains(memArr, mem) {
				memArr = append(memArr, mem)
			}
		}
	}

	parsingArr := findConnection(memArr)
	return parsingArr, nil
}

func contains(memArr []sharedMemData, mem sharedMemData) bool {
	for _, m := range memArr {
		if mem == m {
			return true
		}
	}
	return false
}

func findConnection(memArr []sharedMemData) []ParsingData {
	var parsingArr []ParsingData
	for i := 0; i < len(memArr); i++ {
		for j := i + 1; j < len(memArr); j++ {
			if memArr[i].id == memArr[j].id &&
				memArr[i].pathOfExecutableFile != memArr[j].pathOfExecutableFile {
				memInfo := SharedMemInfo{Ipc: "sharedMemory", Id: memArr[i].id, Key: memArr[i].key, Type: memArr[i].typeIpc}
				parsingArr = append(parsingArr, ParsingData{[2]string{memArr[i].pathOfExecutableFile, memArr[j].pathOfExecutableFile}, memInfo})
			}
		}
	}
	return parsingArr
}
