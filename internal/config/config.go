package config

import (
	genStruct "linux-monitoring-utility/internal/bpfScript/generalStructIPC"
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type BpfTraceConfig map[string]map[string][]string

type ConfigFile struct {
	ScriptTime            uint     `yaml:"scriptTime"`
	ProgramTime           uint     `yaml:"programTime"`
	SyscallsFileName      string   `yaml:"SyscallsFileName"`
	OutputPath            string   `yaml:"outputPath"`
	LsofBinPath           string   `yaml:"lsofBinPath"`
	RpmBinPath            string   `yaml:"rpmBinPath"`
	BpftraceBinPath       string   `yaml:"bpftraceBinPath"`
	TmpPath               string   `yaml:"tmpPath"`
	TmpDelete             bool     `yaml:"tmpDelete"`
	BPFTRACE_STRLEN       string   `yaml:"BPFTRACE_STRLEN"`
	BPFTRACE_MAP_KEYS_MAX string   `yaml:"BPFTRACE_MAP_KEYS_MAX"`
	DirToIgnore           []string `yaml:"DirToIgnore"`
}

func configValidate(configStruct *ConfigFile) error {

	if configStruct.ScriptTime >= configStruct.ProgramTime {
		err := errors.New("script time cannot be more than program time")
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
	if configStruct.SyscallsFileName != "/etc/lmu/lmuSyscalls.yaml" {
		if _, err := os.Stat(configStruct.SyscallsFileName); errors.Is(err, os.ErrNotExist) {
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

	if cliConf.ScriptTime != 0 {

		configStruct.ScriptTime = cliConf.ScriptTime
	}

	if cliConf.ProgramTime != 0 {

		configStruct.ProgramTime = cliConf.ProgramTime
	}

	if len(cliConf.SyscallsFileName) != 0 {

		configStruct.SyscallsFileName = cliConf.SyscallsFileName
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

	//Reading syscalls yaml file
	bpftraceYamlFile, err := os.ReadFile(configStruct.SyscallsFileName)
	if err != nil {
		return nil, err
	}

	var bpftraceConfig BpfTraceConfig

	err = yaml.Unmarshal(bpftraceYamlFile, &bpftraceConfig)
	if err != nil {
		return nil, err
	}

	//////////////////////////////////////////////////////////////
	var bpfSyscalls []genStruct.IpcStruct
	for typeIpc, opts := range bpftraceConfig {
		var oneIpc genStruct.IpcStruct
		oneIpc.IpcType = typeIpc
		for typeOpt, opt := range opts {
			var oneOpt genStruct.OptionStruct
			oneOpt.OptionType = typeOpt
			oneOpt.Options = append(oneOpt.Options, opt...)
			oneIpc.Option = append(oneIpc.Option, oneOpt)
		}
		bpfSyscalls = append(bpfSyscalls, oneIpc)
	}

	return bpfSyscalls, nil
	/////////////////////////////////////////////////////////////
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
	flag.StringVar(&configStruct.SyscallsFileName, "s", "", "Path to .yaml config file with syscalls")
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
