package readWriteParsing

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

type ReadWriteInfo struct {
	Ipc, FileDescriptor, PathOfOpenedFile, ReadBytes, WrittenBytes string
}

func (rw ReadWriteInfo) String() string {
	return fmt.Sprintf("%s: %s, %s, %s, %s", rw.Ipc, rw.FileDescriptor, rw.PathOfOpenedFile, rw.ReadBytes, rw.WrittenBytes)
}

type readWriteData struct {
	pathOfExecutableFile, fileDescriptor, pathOfOpenedFile, readBytes, writtenBytes string
}

func Parse(fileName string) ([]ParsingData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var rwArr []readWriteData

	for fileScanner.Scan() {
		arr := strings.Fields(fileScanner.Text())
		if len(arr) == 5 {
			rw := readWriteData{pathOfExecutableFile: arr[0], fileDescriptor: arr[1], pathOfOpenedFile: arr[2],
				readBytes: arr[3], writtenBytes: arr[4]}
			if !contains(rwArr, rw) {
				rwArr = append(rwArr, rw)
			}
		}
	}

	parsingArr := findConnection(rwArr)
	return parsingArr, nil
}

func contains(rwArr []readWriteData, rw readWriteData) bool {
	for _, s := range rwArr {
		if rw == s {
			return true
		}
	}
	return false
}

func findConnection(rwArr []readWriteData) []ParsingData {
	var parsingArr []ParsingData
	for i := 0; i < len(rwArr); i++ {
		for j := i + 1; j < len(rwArr); j++ {
			if rwArr[i].fileDescriptor == rwArr[j].fileDescriptor &&
				rwArr[i].pathOfExecutableFile != rwArr[j].pathOfExecutableFile &&
				rwArr[i].pathOfOpenedFile == rwArr[j].pathOfOpenedFile &&
				rwArr[i].readBytes != rwArr[j].readBytes { //!!
				rwInfo := ReadWriteInfo{Ipc: "by reading/writing", PathOfOpenedFile: rwArr[i].pathOfOpenedFile,
					FileDescriptor: rwArr[i].fileDescriptor}
				if rwArr[i].readBytes != "0" && rwArr[i].writtenBytes == "0" {
					rwInfo.ReadBytes, rwInfo.WrittenBytes = rwArr[i].readBytes, rwArr[j].writtenBytes
				} else if rwArr[j].readBytes != "0" && rwArr[j].writtenBytes == "0" {
					rwInfo.ReadBytes, rwInfo.WrittenBytes = rwArr[j].readBytes, rwArr[i].writtenBytes
				}
				parsingArr = append(parsingArr, ParsingData{[2]string{rwArr[i].pathOfExecutableFile, rwArr[j].pathOfExecutableFile}, rwInfo})
			}
		}
	}
	return parsingArr
}
