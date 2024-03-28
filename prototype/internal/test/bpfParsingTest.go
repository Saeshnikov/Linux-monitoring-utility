package test

import (
	"bufio"
	"fmt"
	"linux-monitoring-utility/internal/bpfParsing"
	"log"
	"os"
	"reflect"
)

func ParsingTest() bool {
	file, err := os.CreateTemp("./", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	arr_test := []string{"a", "b", "c", "d", "e"}
	for _, s := range arr_test {
		file.WriteString("@filename[" + s + "]\n")
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {

		fmt.Println(fileScanner.Text())
	}

	arr_res, err := bpfParsing.Parse(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(arr_res)
	fmt.Println(arr_test)
	if reflect.DeepEqual(arr_res, arr_test) {
		fmt.Println("wow!")
		return true
	}
	return false
}
