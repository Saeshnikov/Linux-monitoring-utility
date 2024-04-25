package config

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type BpftraceConfig struct {
	Syscalls []string `yaml:"Syscalls"`
}

type ConfigFile struct {
	ScriptTime      uint   `yaml:"scriptTime"`
	ProgramTime     uint   `yaml:"programTime"`
	ConfigFileName  string `yaml:"configFileName"`
	OutputPath      string `yaml:"outputPath"`
	LsofBinPath     string `yaml:"lsofBinPath"`
	RpmBinPath      string `yaml:"rpmBinPath"`
	BpftraceBinPath string `yaml:"bpftraceBinPath"`
	TmpPath         string `yaml:"tmpPath"`
	TmpDelete       bool   `yaml:"tmpDelete"`
}

func configValidate(configStruct *ConfigFile) error {
	if configStruct.ScriptTime >= configStruct.ProgramTime {
		err := errors.New("Script time cannot be more than program time")
		return err
	}

	//Checking existing of bpftrace bin path
	if configStruct.BpftraceBinPath != "/usr/bin/bpftrace" {
		if _, err := os.Stat(configStruct.BpftraceBinPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return err
			}
		}
	}

	//Checking existing of syscalls file
	if configStruct.ConfigFileName != "/etc/lmuConf.yaml" {
		if _, err := os.Stat(configStruct.ConfigFileName); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return err
			}
		}
	}

	//Checking existing of rpm bin file
	if configStruct.ConfigFileName != "/usr/bin/rpm" {
		if _, err := os.Stat(configStruct.RpmBinPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return err
			}
		}
	}

	//Checking existing of lsof bin path
	if configStruct.LsofBinPath != "/usr/bin/lsof" {
		if _, err := os.Stat(configStruct.LsofBinPath); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ConfigRead(configStruct *ConfigFile) ([]string, error) {
	cfgFileFlag := false
	//Checking if cfg flag provided
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "cfg" {
			cfgFileFlag = true
		}
	})

	if cfgFileFlag {
		configFileName := flag.String("cfg", "../configs/defaultCfg.yaml", "Path to the .yaml config")
		flag.Parse()
		err := configFileRead(*configFileName, configStruct)
		if err != nil {
			return nil, err
		}
	} else {
		err := cliRead(configStruct)
		if err != nil {
			return nil, err
		}
	}

	//Validating config struct
	err := configValidate(configStruct)
	if err != nil {
		return nil, err
	}

	//Reading syscalls yaml file
	var config BpftraceConfig

	bpftraceYamlFile, err := os.ReadFile(configStruct.ConfigFileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bpftraceYamlFile, &config)
	if err != nil {
		return nil, err
	}

	return config.Syscalls, nil
}

func configFileRead(configFileName string, configStruct *ConfigFile) error {
	//Reading config file
	configYamlFile, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	//Parsing config file
	err = yaml.Unmarshal(configYamlFile, &configStruct)
	if err != nil {
		return err
	}

	return nil
}

func cliRead(configStruct *ConfigFile) error {

	//Reading command line arguments
	configStruct.ScriptTime = *flag.Uint("t", 3600, "One bpftrace script working time") //BPFtrace script working time
	configStruct.ProgramTime = *flag.Uint("T", 86400, "Program working time")           //Program working time
	configStruct.ConfigFileName = *flag.String("c", "/etc/lmuConf.yaml", "Path to .yaml config file with syscalls")
	configStruct.OutputPath = *flag.String("o", ".", "Path to the result")
	configStruct.LsofBinPath = *flag.String("lsof", "/usr/bin/lsof", "Path to the lsof binary")
	configStruct.BpftraceBinPath = *flag.String("bpf", "/usr/bin/bpftrace", "Path to the Bpftrace binary")
	configStruct.RpmBinPath = *flag.String("rpm", "/usr/bin/rpm", "Path to the rpm binary")
	configStruct.TmpPath = *flag.String("tmp", "/dist", "Path to the tmp folder")
	configStruct.TmpDelete = *flag.Bool("tmpRM", true, "Delete tmp folder or not")

	//Parsing CLI parameters
	flag.Parse()

	return nil
}
