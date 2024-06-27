package namedPipesParsing

import (
	"bufio"
	"fmt"
	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"os"
	"strings"
)

// ---------------------------------------------------------------

//----------------------------------------------------------------

type NamedPipesInfo struct {
	Ipc, Name string
}

func (s NamedPipesInfo) String() string {
	return fmt.Sprintf("%s: %s", s.Ipc, s.Name)
}

type namedPipesData struct {
	pathOfExecutableFile, fileDescriptor, name string
}

func Parse(fileName string) ([]parsingstruct.ParsingData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var pipeArr []namedPipesData

	for fileScanner.Scan() {
		arr := strings.Fields(fileScanner.Text())
		if len(arr) == 3 {
			pipe := namedPipesData{pathOfExecutableFile: arr[0], fileDescriptor: arr[1], name: arr[2]}
			if !contains(pipeArr, pipe) {
				pipeArr = append(pipeArr, pipe)
			}
		}
	}

	parsingArr := findConnection(pipeArr)
	return parsingArr, nil
}

func contains(pipeArr []namedPipesData, pipe namedPipesData) bool {
	for _, p := range pipeArr {
		if pipe == p {
			return true
		}
	}
	return false
}

func findConnection(pipeArr []namedPipesData) []parsingstruct.ParsingData {
	var parsingArr []parsingstruct.ParsingData
	for i := 0; i < len(pipeArr); i++ {
		for j := i + 1; j < len(pipeArr); j++ {
			if pipeArr[i].name == pipeArr[j].name &&
				pipeArr[i].fileDescriptor == pipeArr[j].fileDescriptor &&
				pipeArr[i].pathOfExecutableFile != pipeArr[j].pathOfExecutableFile {
				pipeInfo := NamedPipesInfo{Ipc: "namedPipes", Name: pipeArr[i].name}
				parsingArr = append(parsingArr, parsingstruct.ParsingData{[2]string{pipeArr[i].pathOfExecutableFile, pipeArr[j].pathOfExecutableFile}, pipeInfo})
			}
		}
	}
	return parsingArr
}
