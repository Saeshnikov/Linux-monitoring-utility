package config

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	Syscalls []string `yaml:"Syscalls"`
}

func ConfigRead() (uint, uint, []string, string, string, string, error) {
	var config ConfigFile

	//Reading command line arguments
	scriptTime := flag.Uint("t", 3600, "One bpftrace script working time") //BPFtrace script working time
	programTime := flag.Uint("T", 86400, "Program working time")           //Program working time
	configFileName := flag.String("c", "/etc/lmuConf.yaml", "Path to .yaml config file")
	outputPath := flag.String("o", ".", "Path to the result")
	lsofBinPath := flag.String("lsof", "/usr/bin/lsof", "Path to the lsof binary")
	bpfTraceBinPath := flag.String("bpf", "/usr/bin/bpftrace", "Path to the Bpftrace binary")

	if *scriptTime >= *programTime {
		err := errors.New("Script time cannot be more than program time.")
		return 0, 0, nil, "", "", "", err
	}

	flag.Parse()

	//Opening config file with syscalls
	yamlFile, err := os.ReadFile(*configFileName)
	if err != nil {
		return 0, 0, nil, "", "", "", err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return 0, 0, nil, "", "", "", err
	}

	return *scriptTime, *programTime, config.Syscalls, *outputPath, *lsofBinPath, *bpfTraceBinPath, nil
}
