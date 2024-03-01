package main

import (
	"fmt"
	bpfParsing "linux-monitoring-utility/internal/bpfParsing"
	bpfScript "linux-monitoring-utility/internal/bpfScript"
	config "linux-monitoring-utility/internal/config"
	RPMAnalysis "linux-monitoring-utility/internal/rpmLayer"
	taskExecution "linux-monitoring-utility/internal/taskExecution"
	"log"
	"os"
	"regexp"
	"time"
)

func main() {
	var bpftrace_time int
	var syscalls []string

	bpftrace_time, syscalls = config.ConfigRead()
	bpfScriptFile := bpfScript.GenerateBpfScript(syscalls)
	taskExecution.StartTasks(bpftrace_time, bpfScriptFile.Name(), toRun, toAnalyse)
}

func toRun(bpfTime int, fileName string) *os.File {
	cmdToRun := "/usr/bin/bpftrace"
	args := []string{"", fileName}
	procAttr := new(os.ProcAttr)

	// Создание временного файла для вывода bpftrace
	file, err := os.CreateTemp(".", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	procAttr.Files = []*os.File{os.Stdin, file, os.Stderr}

	// Запуск bpftrace
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
		time.Sleep(time.Duration(bpfTime) * time.Second / 10) // скрипт перезапускается 10 раз
		process.Signal(os.Interrupt)
	}
	return file
}

func toAnalyse(fileForAnalysis *os.File) {

	defer os.Remove(fileForAnalysis.Name())

	r, err := regexp.Compile(`@filename\[(.*?)\]`)
	if err != nil {
		log.Fatal(err)
	}
	res := bpfParsing.Parse(r, fileForAnalysis.Name())

	//Через rpm -qf проверяем относится ли файл к rpm пакету
	RPMAnalysis.RPMlayer(res)
}
