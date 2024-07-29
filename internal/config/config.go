package config

import (
	genStruct "linux-monitoring-utility/internal/bpfScript/generalStructIPC"

	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	ScriptTime            uint                  `yaml:"scriptTime"`
	ProgramTime           uint                  `yaml:"programTime"`
	RpmTasks              uint                  `yaml:"rpmTasks"`
	OutputPath            string                `yaml:"outputPath"`
	LsofBinPath           string                `yaml:"lsofBinPath"`
	RpmBinPath            string                `yaml:"rpmBinPath"`
	BpftraceBinPath       string                `yaml:"bpftraceBinPath"`
	TmpPath               string                `yaml:"tmpPath"`
	TmpDelete             bool                  `yaml:"tmpDelete"`
	BPFTRACE_STRLEN       string                `yaml:"BPFTRACE_STRLEN"`
	BPFTRACE_MAP_KEYS_MAX string                `yaml:"BPFTRACE_MAP_KEYS_MAX"`
	DirToIgnore           []string              `yaml:"DirToIgnore"`
	BpfTraceConfig        []genStruct.IpcStruct `yaml:"Syscalls"`
}

func configValidate(configStruct *ConfigFile) error {

	if configStruct.ScriptTime > configStruct.ProgramTime {
		err := errors.New("script time cannot be more than program time")
		return err
	}

	if configStruct.ScriptTime == configStruct.ProgramTime {
		err := errors.New("you should use -T flag only")
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

	//Checking existing of rpm bin file
	if configStruct.RpmBinPath != "/usr/bin/rpm" {
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

func ConfigRead(configStruct *ConfigFile) ([]genStruct.IpcStruct, error) {

	var cliConf ConfigFile
	var configFileName string
	err := cliRead(&configFileName, &cliConf)
	if err != nil {
		return nil, err
	}
	err = configFileRead(configFileName, configStruct)
	if err != nil {
		return nil, err
	}

	//Checking if only -T flag provided
	if cliConf.ScriptTime == 0 {

		configStruct.RpmTasks = 1
	} else {
		configStruct.ScriptTime = cliConf.ScriptTime
		if cliConf.RpmTasks != 0 {
			configStruct.RpmTasks = cliConf.RpmTasks
		}
	}

	if cliConf.ProgramTime != 0 {

		configStruct.ProgramTime = cliConf.ProgramTime
	}
	if len(cliConf.OutputPath) != 0 {

		configStruct.OutputPath = cliConf.OutputPath
	}
	if len(cliConf.LsofBinPath) != 0 {

		configStruct.LsofBinPath = cliConf.LsofBinPath
	}
	if len(cliConf.RpmBinPath) != 0 {

		configStruct.RpmBinPath = cliConf.RpmBinPath
	}
	if len(cliConf.BpftraceBinPath) != 0 {

		configStruct.BpftraceBinPath = cliConf.BpftraceBinPath
	}
	if len(cliConf.TmpPath) != 0 {

		configStruct.TmpPath = cliConf.TmpPath
	}

	if !cliConf.TmpDelete {
		configStruct.TmpDelete = cliConf.TmpDelete
	}

	//Validating config struct
	err = configValidate(configStruct)
	if err != nil {
		return nil, err
	}

	return configStruct.BpfTraceConfig, nil
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

func cliRead(configFileName *string, configStruct *ConfigFile) error {
	flag.StringVar(configFileName, "cfg", "/etc/lmu/lmuConfig.yaml", "Path to the .yaml config")

	//Reading command line arguments
	flag.UintVar(&configStruct.ScriptTime, "t", 0, "One bpftrace script working time") //BPFtrace script working time
	flag.UintVar(&configStruct.ProgramTime, "T", 0, "Program working time")            //Program working time
	flag.UintVar(&configStruct.RpmTasks, "rpmC", 1, "Number of rpm tasks running at the same time")
	flag.StringVar(&configStruct.OutputPath, "o", "", "Path to the result")
	flag.StringVar(&configStruct.LsofBinPath, "lsof", "", "Path to the lsof binary")
	flag.StringVar(&configStruct.BpftraceBinPath, "bpf", "", "Path to the Bpftrace binary")
	flag.StringVar(&configStruct.RpmBinPath, "rpm", "", "Path to the rpm binary")
	flag.StringVar(&configStruct.TmpPath, "tmp", "", "Path to the tmp folder")
	flag.BoolVar(&configStruct.TmpDelete, "tmpRM", true, "Delete tmp folder or not")

	//Parsing CLI parameters
	flag.Parse()

	return nil
}
