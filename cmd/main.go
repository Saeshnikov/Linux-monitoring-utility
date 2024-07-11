package main

import (
	"bytes"
	"encoding/json"
	bpfParsing "linux-monitoring-utility/internal/bpfParsing/bpftraceParsing"
	"linux-monitoring-utility/internal/bpfParsing/namedPipesParsing"
	parsingstruct "linux-monitoring-utility/internal/bpfParsing/parsingStruct"
	"linux-monitoring-utility/internal/bpfParsing/readWriteParsing"
	"linux-monitoring-utility/internal/bpfParsing/semaphoreParsing"
	"linux-monitoring-utility/internal/bpfParsing/sharedMemParsing"
	bpfScript "linux-monitoring-utility/internal/bpfScript"
	config "linux-monitoring-utility/internal/config"
	lsofLayer "linux-monitoring-utility/internal/lsofLayer"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	taskExecution "linux-monitoring-utility/internal/taskExecution"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var programConfig config.ConfigFile
var pathToTmp string

func main() {

	syscalls, err := config.ConfigRead(&programConfig)
	if err != nil {
		log.Fatal(err)
	}

	// os.Setenv("BPFTRACE_STRLEN", programConfig.BPFTRACE_STRLEN)
	os.Setenv("BPFTRACE_MAP_KEYS_MAX", programConfig.BPFTRACE_MAP_KEYS_MAX)
	os.Setenv("BPFTRACE_LOG_SIZE", strconv.Itoa(10000000))

	lsofLayer.LsofBinPath = programConfig.LsofBinPath
	lsofLayer.DirToIgnore = programConfig.DirToIgnore
	bpfParsing.DirToIgnore = programConfig.DirToIgnore
	rpmLayer.RpmBinPath = programConfig.RpmBinPath
	inodeStr, err := exec.Command("ls", "-id", "/").Output()
	if err != nil {
		log.Fatal(err)
	}
	inodeStr = inodeStr[:len(inodeStr)-3]
	inodeInt, err := strconv.Atoi(string(inodeStr))
	if err != nil {
		log.Fatal(err)
	}
	bpfScriptFiles, err := bpfScript.GenerateBpfScript(syscalls, programConfig.OutputPath, inodeInt)
	if err != nil {
		log.Fatal(err)
	}

	if programConfig.OutputPath == "." {
		err = os.MkdirAll("out", os.FileMode(0777))

		if err != nil {
			log.Fatal(err)
		}
	}
	pathToTmp = programConfig.TmpPath
	os.MkdirAll(pathToTmp+"/tmp", os.FileMode(0777))

	if programConfig.TmpDelete {
		defer os.RemoveAll(pathToTmp + "/tmp")
	}

	// outputMap, err := rpmLayer.FindAllPackages()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var bpfCommands []taskExecution.ExecUnit
	for _, i := range bpfScriptFiles {
		dir := i.Name()[:len(i.Name())-3]
		bpf := taskExecution.NewExecUnitContinuousF(programConfig.BpftraceBinPath,
			[]string{i.Name()},
			uint(programConfig.ProgramTime/programConfig.ScriptTime),
			time.Duration(programConfig.ScriptTime)*time.Second, pathToTmp+"/tmp/"+dir)
		bpfCommands = append(bpfCommands, *bpf)
	}
	err = taskExecution.StartTasks(bpfCommands...)

	// bpf := taskExecution.NewExecUnitContinuousF(programConfig.BpftraceBinPath,
	// 	[]string{"sc.bt"},
	// 	uint(programConfig.ProgramTime/programConfig.ScriptTime),
	// 	time.Duration(programConfig.ScriptTime)*time.Second, pathToTmp+"/tmp/"+"fsorw")
	// err = taskExecution.StartTasks(*bpf)
	if err != nil {
		log.Fatal(err)
	}

	parsedData, err := toParse(pathToTmp + "/tmp")
	if err != nil {
		log.Fatal(err)
	}
	// analysedData, err := toAnalyse(parsedData)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	file, err := os.Create("out/" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString("PACKAGE1 \t\tPACKAGE2 \t\tINTERACTION\n")
	for _, data := range parsedData {
		file.WriteString(data.PathsOfExecutableFiles[0] + "\t , \t" + data.PathsOfExecutableFiles[1] + "\t : \t" + data.WayOfInteraction.String() + "\n")
	}
	// err = exportToJson(programConfig.OutputPath, outputMap)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func exportToJson(filePath string, outputMap map[string]bool) error {

	entriesArr := make([]string, 0)

	for entry := range outputMap {
		entriesArr = append(entriesArr, entry)
	}
	jsonArray, err := json.Marshal(entriesArr)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(filePath + "/result.json")
	if err != nil {
		return err
	}
	_, err = outputFile.Write(jsonArray)
	if err != nil {
		return err
	}
	return nil
}

func toParse(directory string) ([]parsingstruct.ParsingData, error) {
	var parsingInfo []parsingstruct.ParsingData
	files, err := os.ReadDir(directory + "/fsorw/")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				data, err := readWriteParsing.Parse(directory + "/fsorw/" + file.Name())
				if err != nil {
					return nil, err
				}
				parsingInfo = append(parsingInfo, data...)
			}
		}
	}

	files, err = os.ReadDir(directory + "/namedpipe/")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				data, err := namedPipesParsing.Parse(directory + "/namedpipe/" + file.Name())
				if err != nil {
					return nil, err
				}
				parsingInfo = append(parsingInfo, data...)
			}
		}
	}

	files, err = os.ReadDir(directory + "/semaphore/")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				data, err := semaphoreParsing.Parse(directory + "/semaphore/" + file.Name())
				if err != nil {
					return nil, err
				}
				parsingInfo = append(parsingInfo, data...)
			}
		}
	}

	files, err = os.ReadDir(directory + "/sharedMem/")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				data, err := sharedMemParsing.Parse(directory + "/sharedMem/" + file.Name())
				if err != nil {
					return nil, err
				}
				parsingInfo = append(parsingInfo, data...)
			}
		}
	}

	return parsingInfo, nil
}

// not available yet
func toAnalyse(data []parsingstruct.ParsingData) ([]parsingstruct.ParsingData, error) {
	newParsingData := make([]parsingstruct.ParsingData, len(data))
	for n, unit := range data {
		var ch1 chan chan bytes.Buffer
		var ch2 chan chan bytes.Buffer
		tmp1 := taskExecution.NewExecUnitOneShotC(programConfig.RpmBinPath, []string{"-qf", unit.PathsOfExecutableFiles[0]}, 1, ch1)
		tmp2 := taskExecution.NewExecUnitOneShotC(programConfig.RpmBinPath, []string{"-qf", unit.PathsOfExecutableFiles[1]}, 1, ch2)
		go func(n int, i parsingstruct.Interaction) {
			newParsingData[n].WayOfInteraction = i
			newPaths := make([]string, 2)
			c1 := <-ch1
			c2 := <-ch2
			select {
			case b := <-c1:
				newPaths[0] = b.String()
			case b := <-c2:
				newPaths[1] = b.String()
			}
			newParsingData[n].PathsOfExecutableFiles = [2]string(newPaths)
		}(n, unit.WayOfInteraction)
		err := taskExecution.StartTasks(tmp1, tmp2)
		if err != nil {
			return nil, err
		}
	}
	return newParsingData, nil
}
