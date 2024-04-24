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
	defer os.RemoveAll("tmp")
	outputMap, err := rpmLayer.FindAllPackages()
	if err != nil {
		log.Fatal(err)
	}

	err = taskExecution.StartTasks(program_time, bpftrace_time, bpfScriptFile.Name(), toRun, toRunLsof)
	if err != nil {
		log.Fatal(err)
	}

	err = toAnalyse("./tmp/", outputPath, &outputMap)
	if err != nil {
		log.Fatal(err)
	}

	err = exportToJson(outputPath, outputMap)
	if err != nil {
		log.Fatal(err)
	}
}

func toRunLsof() {
	arr, err := lsofLayer.LsofExec()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("./tmp/" + "lsofOutput")
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range arr {
		file.WriteString(line + "\n")
	}
}

func toRun(fileName string, c chan *exec.Cmd) {
	file, err := os.Create("./tmp/" + time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cmd := exec.Command("/usr/bin/bpftrace", fileName)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(pipe)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Procces is running as pid %d\n", cmd.Process.Pid)
	line, err := reader.ReadString('\n')

	c <- cmd

	for err == nil {
		file.WriteString(line)
		line, err = reader.ReadString('\n')
	}

	cmd.Wait()

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
	outputFile.Write(jsonArray)
	return nil
}

func toAnalyse(directory string, dirPath string, outputMap *map[string]bool) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			if file.Name() == "lsofOutput" {
				var res []string
				filePath := filepath.Join(directory, file.Name())
				f, err := os.Open(filePath)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line := scanner.Text()
					res = append(res, line)
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
