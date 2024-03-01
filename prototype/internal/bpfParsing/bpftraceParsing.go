package bpfParsing

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

func Parse(r *regexp.Regexp, fileName string) []string {

	var arr []string

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {

		res := r.FindAllStringSubmatch(fileScanner.Text(), -1)
		if res != nil {
			arr = append(arr, res[0][1])
			println(res[0][1])
		}

	}

	return arr
}
