package bpfParsing

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var DirToIgnore []string

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
			if len(res[0][1]) > 1 && ignoreDir(res[0][1]) && len(strings.Fields(res[0][1])) == 1 {
				arr = append(arr, res[0][1])
			}
		}
	}

	return arr, nil

}

func ignoreDir(s string) bool {
	for _, dir := range DirToIgnore {
		if len(strings.Split(s, dir)) != 1 {
			return false
		}
	}
	return true
}
