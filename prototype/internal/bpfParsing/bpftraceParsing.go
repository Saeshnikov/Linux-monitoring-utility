package bpfParsing

import (
	"bufio"
	"os"
	"regexp"
)

func Parse(fileName string) ([]string, error) {

	var arr []string
	r, err := regexp.Compile(`@filename\[(.*?)\]`)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {

		res := r.FindAllStringSubmatch(fileScanner.Text(), -1)
		if res != nil {
			arr = append(arr, res[0][1])
		}
	}

	return arr, nil

}
