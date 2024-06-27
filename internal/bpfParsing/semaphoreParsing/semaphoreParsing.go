package semaphoreParsing

import (
	"bufio"
	"fmt"
	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"os"
	"strings"
)

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

func Parse(fileName string) ([]parsingstruct.ParsingData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var semArr []semaphoreData

	for fileScanner.Scan() {
		arr := strings.Fields(fileScanner.Text())
		if len(arr) == 3 {
			sem := semaphoreData{pathOfExecutableFile: arr[0], key: arr[1], id: arr[2]}
			if !contains(semArr, sem) {
				semArr = append(semArr, sem)
			}
		}
	}

	parsingArr := findConnection(semArr)
	return parsingArr, nil
}

func contains(semArr []semaphoreData, sem semaphoreData) bool {
	for _, s := range semArr {
		if sem == s {
			return true
		}
	}
	return false
}

func findConnection(semArr []semaphoreData) []parsingstruct.ParsingData {
	var parsingArr []parsingstruct.ParsingData
	for i := 0; i < len(semArr); i++ {
		for j := i + 1; j < len(semArr); j++ {
			if semArr[i].id == semArr[j].id &&
				semArr[i].pathOfExecutableFile != semArr[j].pathOfExecutableFile {
				semInfo := SemaphoreInfo{Ipc: "semaphore", Id: semArr[i].id}
				parsingArr = append(parsingArr, parsingstruct.ParsingData{PathsOfExecutableFiles: [2]string{semArr[i].pathOfExecutableFile, semArr[j].pathOfExecutableFile}, WayOfInteraction: semInfo})
			}
		}
	}
	return parsingArr
}
