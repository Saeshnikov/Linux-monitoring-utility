package rpmanalysis

import (
	"bytes"
	"fmt"
	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"linux-monitoring-utility/internal/bpfParsing/readWriteParsing"
	"linux-monitoring-utility/internal/taskExecution"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// not available yet
func ToAnalyse(data []parsingstruct.ParsingData, rpmBinPath string, count int) ([]parsingstruct.ParsingData, error) {
	var wg sync.WaitGroup
	var newParsingData []parsingstruct.ParsingData
	ch := make([]chan chan bytes.Buffer, count)
	for i := 0; i <= (len(data)-1)*2; i += count {
		var unit []taskExecution.ExecUnit
		s := make([]string, count)
		for j := 0; j < count; j++ {
			if (i+j)/2 < len(data) {
				ch[j] = make(chan chan bytes.Buffer, 1)
				if data[(i+j)/2].PathsOfExecutableFiles[(i+j)%2] != "-" {
					unit = append(unit, *taskExecution.NewExecUnitOneShotC(rpmBinPath, []string{"-qf", data[(i+j)/2].PathsOfExecutableFiles[(i+j)%2]}, 1, ch[j]))
					wg.Add(1)
					go func(j int) {
						defer wg.Done()
						c := <-ch[j]
						buf := <-c
						b := buf.String()
						if b == "" || len(strings.Split(b, " ")) > 1 {
							return
						}
						s[j] = b[:len(b)-1]
					}(j)
				} else {
					s[j] = ""
				}
			}
		}
		fmt.Println(strconv.Itoa(i) + " ===========")
		err := taskExecution.StartTasks(unit...)
		if err != nil {
			return nil, err
		}
		wg.Wait()
		for j := 0; j < count; j++ {
			if (i+j)/2 < len(data) {
				if (i+j)%2 == 0 {
					newParsingData = append(newParsingData, parsingstruct.ParsingData{})
					newParsingData[len(newParsingData)-1].WayOfInteraction = data[(i+j)/2].WayOfInteraction
				}
				newParsingData[len(newParsingData)-1].PathsOfExecutableFiles[(i+j)%2] = s[j]
				if (i+j)%2 == 1 && (newParsingData[len(newParsingData)-1].PathsOfExecutableFiles[0] == "" || (newParsingData[len(newParsingData)-1].PathsOfExecutableFiles[1] == "" && reflect.TypeOf(newParsingData[len(newParsingData)-1].WayOfInteraction) != reflect.TypeOf(readWriteParsing.ReadWriteInfo{}))) {
					newParsingData = newParsingData[:len(newParsingData)-1]
				}
			}
		}
	}
	return newParsingData, nil
}
