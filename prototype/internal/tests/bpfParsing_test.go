package tests

import (
	"linux-monitoring-utility/internal/bpfParsing"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := os.CreateTemp("./", "tmp")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(file.Name())
	arr_test := []string{"a", "b", "c", "d", "e"}
	for _, s := range arr_test {
		file.WriteString("@filename[" + s + "]\n")
	}

	arr_res, err := bpfParsing.Parse(file.Name())
	if err != nil {
		t.Fatal(err.Error())
	}
	if reflect.DeepEqual(arr_res, arr_test) {
		return
	}
	t.Fatal("Test failed")
}
