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

	parsingArr := packData(rwArr)
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

func packData(rwArr []readWriteData) []ParsingData {
	var parsingArr []ParsingData
	for i := 0; i < len(rwArr); i++ {
		rwInfo := ReadWriteInfo{PathOfOpenedFile: rwArr[i].pathOfOpenedFile,
			FileDescriptor: rwArr[i].fileDescriptor, ReadBytes: rwArr[i].readBytes, WrittenBytes: rwArr[i].writtenBytes}
		if rwArr[i].readBytes != "0" && rwArr[i].writtenBytes == "0" {
			rwInfo.Ipc = "reading"
			rwInfo.ReadBytes, rwInfo.WrittenBytes = rwArr[i].readBytes, "-"
		} else if rwArr[i].readBytes == "0" && rwArr[i].writtenBytes != "0" {
			rwInfo.Ipc = "writing"
			rwInfo.ReadBytes, rwInfo.WrittenBytes = "-", rwArr[i].writtenBytes
		}
		rwArr[i].pathOfExecutableFile = strings.TrimPrefix(rwArr[i].pathOfExecutableFile, "/snapshot")
		parsingArr = append(parsingArr, ParsingData{[2]string{rwArr[i].pathOfExecutableFile, "-"}, rwInfo})
	}
	return parsingArr
}
