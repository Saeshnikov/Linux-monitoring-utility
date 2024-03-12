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
	"time"
)

func main() {
	bpftrace_time, program_time, syscalls, err := config.ConfigRead()

	if err != nil {
		log.Fatal(err)
	}

	bpfScriptFile := bpfScript.GenerateBpfScript(syscalls)
	taskExecution.StartTasks(program_time, bpftrace_time, bpfScriptFile.Name(), toRun)
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
	fmt.Printf("Script started...\n")
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
		time.Sleep(time.Duration(bpftrace_time) * time.Second)
		process.Signal(os.Interrupt)
		fmt.Printf("Script stoped...\n")
	}
	toAnalyse(file)
}

func toAnalyse(fileForAnalysis *os.File) {

	defer os.Remove(fileForAnalysis.Name())

	res, err := bpfParsing.Parse(fileForAnalysis.Name())
	if err != nil {
		log.Fatal(err)
	}
	//Через rpm -qf проверяем относится ли файл к rpm пакету
	RPMAnalysis.RPMlayer(res)
}
