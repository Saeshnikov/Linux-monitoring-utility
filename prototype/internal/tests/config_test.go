package test

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	configFileName := "../tests/data/invalidCfg.yaml"
	yamlFile, err := os.ReadFile(*configFileName)
	if err != nil {
		return 0, 0, nil, "", err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		t.fatal(err.Error())
	}

	bpfScriptFile, err := bpfScript.GenerateBpfScript(syscalls, "")
	if err == nil {
		t.Fatal("BpfScriptLayer accepted invalid config file.")
	}
	return
}
