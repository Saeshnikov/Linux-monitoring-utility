package tests

import (
	"os"
	"testing"

	bpfScript "linux-monitoring-utility/internal/bpfScript"
	config "linux-monitoring-utility/internal/config"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	bpfConfigFileName := "../tests/data/invalidBpfCfg.yaml"
	yamlFile, err := os.ReadFile(bpfConfigFileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	var bpfCfg config.BpftraceConfig
	err = yaml.Unmarshal(yamlFile, &bpfCfg)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = bpfScript.GenerateBpfScript(bpfCfg.Syscalls, "")
	if err == nil {
		t.Fatal("BpfScriptLayer accepted invalid config file.")
	}
}
