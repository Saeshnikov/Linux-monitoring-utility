package RPMAnalysis

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"os/exec"
)

func RPMlayer(usedPackages []string) {
	allPackages := findAllPackages()
	findUsedPackages(usedPackages, allPackages)
	findUnusedPackages(allPackages)
}

func findAllPackages() map[string]int {
	var allPackages map[string]int
	mCache, err := os.Open("m_cache")
	if err == nil {
		d := gob.NewDecoder(mCache)
		// Decoding the serialized data
		err = d.Decode(&allPackages)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		allPackages = make(map[string]int)
		cmd := exec.Command("/usr/bin/rpm", "-qa")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		outScanner := bufio.NewScanner(stdout)
		for outScanner.Scan() {
			allPackages[outScanner.Text()] = 0
		}

		b := new(bytes.Buffer)
		e := gob.NewEncoder(b)
		// Encoding the map
		err = e.Encode(allPackages)
		if err != nil {
			log.Fatal(err)
		}
		mCache, err := os.Create("m_cache")
		if err != nil {
			log.Fatal(err)
		}
		defer mCache.Close()
		mCache.Write(b.Bytes())
	}
	return allPackages
}

func findUnusedPackages(allPackages map[string]int) {
	file, err := os.Create("unusedPackages.txt")
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString("Not used packages:\n")
	for key, value := range allPackages {
		if value == 0 {
			file.WriteString(key + "\n")
		}
	}
}

func findUsedPackages(usedPackages []string, allPackages map[string]int) {
	var packageName string

	for _, fileName := range usedPackages {
		cmd := exec.Command("/usr/bin/rpm", "-qf", fileName)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		outScanner := bufio.NewScanner(stdout)
		outScanner.Split(bufio.ScanWords)
		cnt := 0
		for outScanner.Scan() {
			packageName = outScanner.Text()
			cnt++
			if cnt > 1 {
				break
			}
		}
		if cnt == 1 {
			if _, ok := allPackages[packageName]; ok {
				allPackages[packageName]++
			} else {
				log.Fatal("Found New Package") //
			}
		}
	}
}
