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
	bpftrace_time, program_time, syscalls, outputPath, err := config.ConfigRead()
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

	err = taskExecution.StartTasks(program_time, bpftrace_time, bpfScriptFile.Name(), outputPath, &outputMap, toRun)
	if err != nil {
		log.Fatal(err)
	}

	err = exportToJson(outputPath, outputMap)
	if err != nil {
		log.Fatal(err)
	}
}

func toRun(bpftrace_time uint, fileName string, outputPath string, outputMap *map[string]bool) {
	cmdToRun := "/usr/bin/bpftrace"
	args := []string{"", fileName}
	procAttr := new(os.ProcAttr)

	// Создание временного файла для вывода bpftrace
	file, err := os.CreateTemp("./tmp/", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	procAttr.Files = []*os.File{os.Stdin, file, os.Stderr}

	// Запуск bpftrace
	fmt.Printf("Script started...\n")
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
		time.Sleep(time.Duration(bpftrace_time) * time.Second)
		process.Signal(os.Interrupt)
		fmt.Printf("Script stoped...\n")
	}
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
