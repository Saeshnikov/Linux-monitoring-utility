package semaphoreParsing

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

type SemaphoreInfo struct {
	Ipc, Id string
}

func (s SemaphoreInfo) String() string {
	return fmt.Sprintf("%s: %s", s.Ipc, s.Id)
}

type semaphoreData struct {
	pathOfExecutableFile, key, id string
}

func Parse(fileName string) ([]ParsingData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var semArr []semaphoreData

	for fileScanner.Scan() {
		arr := strings.Fields(fileScanner.Text())
		if len(arr) == 3 {
			semArr = append(semArr, semaphoreData{pathOfExecutableFile: arr[0], key: arr[1], id: arr[2]})
		}
	}

	parsingArr := findConnection(semArr)
	return parsingArr, nil
}

func findConnection(semArr []semaphoreData) []ParsingData {
	var parsingArr []ParsingData
	for i := 0; i < len(semArr); i++ {
		for j := 1; j < len(semArr); j++ {
			if semArr[i].id == semArr[j].id &&
				semArr[i].pathOfExecutableFile != semArr[j].pathOfExecutableFile {
				semInfo := SemaphoreInfo{Ipc: "by semaphore", Id: semArr[i].id}
				parsingArr = append(parsingArr, ParsingData{semArr[i].pathOfExecutableFile, semArr[j].pathOfExecutableFile, semInfo})
				semArr[j].id = "" //!!!!!!!!!!!!1
			}
		}
	}
	return parsingArr
}
