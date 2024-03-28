package test

import (
	"linux-monitoring-utility/internal/bpfParsing"
	"log"
	"os"
	"reflect"
)

func TestParsing() bool {
	file, err := os.CreateTemp("./", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	arr_test := []string{"a", "b", "c", "d", "e"}
	for _, s := range arr_test {
		file.WriteString("@filename[" + s + "]\n")
	}

	arr_res, err := bpfParsing.Parse(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	if reflect.DeepEqual(arr_res, arr_test) {
		return true
	}
	return false
}
