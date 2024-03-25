package main

import (
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

	bpfScriptFile, err := bpfScript.GenerateBpfScript(syscalls)
	if err != nil {
		log.Fatal(err)
	}

	if outputPath == "" {
		os.Mkdir("out", os.FileMode(0777))
	}

	os.Mkdir("tmp", os.FileMode(0777))
	err = taskExecution.StartTasks(program_time, bpftrace_time, bpfScriptFile.Name(), outputPath, toRun)
	if err != nil {
		log.Fatal(err)
	}
}

func toRun(bpftrace_time int, fileName string, outputPath string) {
	cmdToRun := "/usr/bin/bpftrace"
	args := []string{"", fileName}
	procAttr := new(os.ProcAttr)

	// Создание временного файла для вывода bpftrace
	file, err := os.CreateTemp("./tmp", "tmp")
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
	toAnalyse(file, outputPath)
}

func toAnalyse(fileForAnalysis *os.File, outputPath string) {
	defer os.Remove(fileForAnalysis.Name())

	res, err := bpfParsing.Parse(fileForAnalysis.Name())
	if err != nil {
		log.Fatal(err)
	}
	//Через rpm -qf проверяем относится ли файл к rpm пакету
	rpmLayer.RPMlayer(res, outputPath)
}
