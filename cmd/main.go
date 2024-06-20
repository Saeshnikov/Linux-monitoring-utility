package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	bpfParsing "linux-monitoring-utility/internal/bpfParsing"
	bpfScript "linux-monitoring-utility/internal/bpfScript"
	config "linux-monitoring-utility/internal/config"
	lsofLayer "linux-monitoring-utility/internal/lsofLayer"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	taskExecution "linux-monitoring-utility/internal/taskExecution"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var programConfig config.ConfigFile
var pathToTmp string

func main() {

	syscalls, err := config.ConfigRead(&programConfig)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("BPFTRACE_STRLEN", programConfig.BPFTRACE_STRLEN)
	os.Setenv("BPFTRACE_MAP_KEYS_MAX", programConfig.BPFTRACE_MAP_KEYS_MAX)

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

	if programConfig.OutputPath == "" {
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

	outputMap, err := rpmLayer.FindAllPackages()
	if err != nil {
		log.Fatal(err)
	}

	lsof := taskExecution.NewExecUnitOneShotF("/usr/bin/lsof", "", 1, pathToTmp+"/tmp/lsof")

	var bpfCommands []taskExecution.ExecUnit
	bpfCommands = append(bpfCommands, *lsof)
	for _, i := range bpfScriptFiles {
		dir := i.Name()[:len(i.Name())-3]
		bpf := taskExecution.NewExecUnitContinuousF(programConfig.BpftraceBinPath,
			i.Name(),
			uint(programConfig.ProgramTime/programConfig.ScriptTime),
			time.Duration(programConfig.ScriptTime)*time.Second, pathToTmp+"/tmp/"+dir)
		bpfCommands = append(bpfCommands, *bpf)
	}
	fmt.Println("a")
	err = taskExecution.StartTasks(pathToTmp, bpfCommands...)
	if err != nil {
		log.Fatal(err)
	}

	err = toAnalyse(pathToTmp+"/tmp/", programConfig.OutputPath, &outputMap)
	if err != nil {
		log.Fatal(err)
	}

	err = exportToJson(programConfig.OutputPath, outputMap)
	if err != nil {
		log.Fatal(err)
	}
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

func toAnalyse(directory string, dirPath string, outputMap *map[string]bool) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			if strings.Split(file.Name(), ".")[0] == "lsof" {
				var res []string
				filePath := filepath.Join(directory, file.Name())
				f, err := os.Open(filePath)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				scanner := bufio.NewScanner(f)
				res, err = lsofLayer.LsofParsing(scanner)
				if err != nil {
					return err
				}

				fmt.Print("File with name: ", file.Name(), " to analyse... ")
				rpmLayer.RPMlayer(res, dirPath, outputMap)
				fmt.Print(" DONE\n")

			} else {
				filePath := filepath.Join(directory, file.Name())
				res, err := bpfParsing.Parse(filePath)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Print("File with name: ", file.Name(), " to analyse... ")
				rpmLayer.RPMlayer(res, dirPath, outputMap)
				fmt.Print(" DONE\n")
			}
		}

	}
	return nil
}
