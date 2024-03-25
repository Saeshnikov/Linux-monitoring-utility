package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	Syscalls []string `yaml:"Syscalls"`
}

func ConfigRead() (int, int, []string, string, error) {
	var config ConfigFile

	//Reading command line arguments
	scriptTime := flag.Int("t", 3600, "One bpftrace script working time") //BPFtrace script working time
	programTime := flag.Int("T", 86400, "Program working time")           //Program working time
	configFileName := flag.String("c", "../configs/defaultConf.yaml", "Path to .yaml config file")
	outputPath := flag.String("o", ".", "Path to the result")

	flag.Parse()

	//Opening config file with syscalls
	yamlFile, err := os.ReadFile(*configFileName)
	if err != nil {
		return 0, 0, nil, "", err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return 0, 0, nil, "", err
	}

	return *scriptTime, *programTime, config.Syscalls, *outputPath, nil
}
