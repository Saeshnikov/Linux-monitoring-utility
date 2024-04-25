package tests

import (
	"linux-monitoring-utility/internal/lsofLayer"
	"os"
	"testing"
)

func estLsofExec(t *testing.T) {
	file, err := os.Create("/usr/tmp")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(file.Name())
	arr_res, err := lsofLayer.LsofExec("usr/bin/lsof")
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, s := range arr_res {
		if s == file.Name() {
			return
		}
	}
	t.Fatal("Test failed")
}
