package test

import (
	"linux-monitoring-utility/internal/bpfParsing"
	"os"
	"reflect"
)

func TestParsing() (bool, error) {
	file, err := os.CreateTemp("./", "tmp")
	if err != nil {
		return false, err
	}
	defer os.Remove(file.Name())
	arr_test := []string{"a", "b", "c", "d", "e"}
	for _, s := range arr_test {
		file.WriteString("@filename[" + s + "]\n")
	}

	arr_res, err := bpfParsing.Parse(file.Name())
	if err != nil {
		return false, err
	}
	if reflect.DeepEqual(arr_res, arr_test) {
		return true, nil
	}
	return false, nil
}
