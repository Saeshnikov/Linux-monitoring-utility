package tests

import (
	"os"
	"testing"

	bpfScript "linux-monitoring-utility/internal/bpfScript"
	config "linux-monitoring-utility/internal/config"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	configFileName := "../tests/data/invalidCfg.yaml"
	yamlFile, err := os.ReadFile(configFileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	var bpfCfg config.ConfigFile
	err = yaml.Unmarshal(yamlFile, &bpfCfg)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = bpfScript.GenerateBpfScript(bpfCfg.BpfTraceConfig, "", 2)
	if err == nil {
		t.Fatal("BpfScriptLayer accepted invalid config.")
	}
}
