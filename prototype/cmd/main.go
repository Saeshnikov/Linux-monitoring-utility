package main

import (
	"encoding/json"
	"fmt"
	bpfParsing "linux-monitoring-utility/internal/bpfParsing"
	bpfScript "linux-monitoring-utility/internal/bpfScript"
	config "linux-monitoring-utility/internal/config"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	taskExecution "linux-monitoring-utility/internal/taskExecution"
	"log"
	"os"
	"time"
)

func main() {
	bpftrace_time, program_time, syscalls, outputPath, lsofBinPath, bpfTraceBinPath, err := config.ConfigRead()
	if err != nil {
		log.Fatal(err)
	}

	bpfScriptFile, err := bpfScript.GenerateBpfScript(syscalls, outputPath)
	if err != nil {
		log.Fatal(err)
	}

	if outputPath == "" {
		os.Mkdir("out", os.FileMode(0777))
	}

	os.Mkdir("tmp", os.FileMode(0777))

	outputMap, err := rpmLayer.FindAllPackages()
	if err != nil {
		log.Fatal(err)
	}

	err = taskExecution.StartTasks(program_time, bpftrace_time, bpfScriptFile.Name(), outputPath, lsofBinPath, bpfTraceBinPath, &outputMap, toRun)
	if err != nil {
		log.Fatal(err)
	}

	err = exportToJson(outputPath, outputMap)
	if err != nil {
		log.Fatal(err)
	}
}

func toRun(bpftrace_time uint, fileName string, outputPath string, bpfTraceBinPath string, outputMap *map[string]bool, c chan *os.Process) {
	args := []string{"", fileName}
	procAttr := new(os.ProcAttr)

	// Создание временного файла для вывода bpftrace
	file, err := os.CreateTemp("./tmp/", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	procAttr.Files = []*os.File{os.Stdin, file, os.Stderr}
	// Запуск bpftrace
	if process, err := os.StartProcess(bpfTraceBinPath, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", bpfTraceBinPath, err.Error())
	} else {
		c <- process
		fmt.Printf("%s running as pid %d\n", bpfTraceBinPath, process.Pid)
	}
	time.Sleep(time.Duration(bpftrace_time) * time.Second)
	toAnalyse(file, outputPath, outputMap)
}

func toAnalyse(fileForAnalysis *os.File, outputPath string, outputMap *map[string]bool) {
	defer os.Remove(fileForAnalysis.Name())

	res, err := bpfParsing.Parse(fileForAnalysis.Name())
	if err != nil {
		log.Fatal(err)
	}
	//Через rpm -qf проверяем относится ли файл к rpm пакету
	rpmLayer.RPMlayer(res, outputPath, outputMap)
}

func exportToJson(filePath string, outputMap map[string]bool) error {
	entriesArr := make([]string, 0)

	for entry, _ := range outputMap {
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
	outputFile.Write(jsonArray)
	return nil
}
