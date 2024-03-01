package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func ConfigRead() (int, []string) {
	//Reading command line arguments
	progArgs := os.Args[1:]

	if len(progArgs) < 1 || len(progArgs) > 2 {
		fmt.Println("Error! Invalid number of arguments.")
		os.Exit(1)
	}

	var bpftrace_time_str = progArgs[0] //BPFtrace working time

	//Converting bpftrace working time to int
	bpftrace_time, err := strconv.Atoi(bpftrace_time_str)
	if err != nil {
		panic(err)
	}

	var syscalls []string
	if len(progArgs) == 2 {
		var configFilePath = progArgs[1] //Syscalls to trace
		//Opening config file with syscalls
		file, err := os.Open(configFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		//Reading config file
		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			syscalls = append(syscalls, fileScanner.Text())
		}
	}

	return bpftrace_time, syscalls
}
