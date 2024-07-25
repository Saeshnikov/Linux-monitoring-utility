package main

import (
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
	rpmanalysis "linux-monitoring-utility/internal/rpmAnalysis"
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
		programConfig.OutputPath = "./out"
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

	if err != nil {
		log.Fatal(err)
	}

	parsedData, err := toParse(pathToTmp + "/tmp")
	if err != nil {
		log.Fatal(err)
	}

	analysedData, err := rpmanalysis.ToAnalyse(parsedData, programConfig.RpmBinPath, 4)
	if err != nil {
		log.Fatal(err)
	}

	exportToJson(programConfig.OutputPath, analysedData)

	// file, err := os.Create("out/" + strconv.FormatInt(time.Now().Unix(), 10))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// file.WriteString("PACKAGE1 \t\tPACKAGE2 \t\tINTERACTION\n")
	// for _, data := range analysedData {
	// 	file.WriteString(data.PathsOfExecutableFiles[0] + "\t , \t" + data.PathsOfExecutableFiles[1] + "\t : \t" + data.WayOfInteraction.String() + "\n")
	// }
}

func exportToJson(filePath string, data []parsingstruct.ParsingData) error {

	jsonArray, err := json.Marshal(data)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(filePath + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".json")
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
