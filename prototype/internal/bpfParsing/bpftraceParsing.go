package bpfParsing

import (
	"bufio"
	"os"
	"regexp"
	"strings"
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
			if len(res[0][1]) > 1 && len(strings.Split(res[0][1], "/proc/")) == 1 && len(strings.Split(res[0][1], " /dev/")) == 1 && len(strings.Fields(res[0][1])) == 1 {
				arr = append(arr, res[0][1])
			}
		}
	}

	return arr, nil

}
