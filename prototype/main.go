package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var TIME_BPFTRACE = 5
var m map[string]int

func main() {

	m_init()

	/*
		bpftrace is working
	*/

	cmdToRun := "/usr/bin/bpftrace"
	args := []string{"", "script.bt"}
	procAttr := new(os.ProcAttr)
	// Временный файл под вывод bpftrace (я НАДЕЮСЬ можно без этого)
	file, err := os.CreateTemp(".", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	procAttr.Files = []*os.File{os.Stdin, file, os.Stderr}
	//Запуск bpftrace
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
		//Ждем TIME_BPFTRACE секунд, после заканчиваем работу bpftrace
		time.Sleep(time.Duration(TIME_BPFTRACE) * time.Second)
		process.Signal(os.Interrupt)
		process.Wait()
	}

	/*
		reading bpftrace out
	*/
	//Открываем reader файла
	file, err = os.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)
	r, err := regexp.Compile(`@filename\[(.*?)\]`)
	if err != nil {
		log.Fatal(err)
	}
	//Читаем пути файлам (тут проблемки конеч)
	for fileScanner.Scan() {

		res := r.FindAllStringSubmatch(fileScanner.Text(), -1)
		if res != nil {

			/*
				rpm analysis
			*/
			//Через rpm -qf проверяем относится ли файл к rpm пакету
			cmd := exec.Command("/usr/bin/rpm", "-qf", res[0][1])
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
			var pkg string
			for outScanner.Scan() {
				pkg = outScanner.Text()
				cnt++
				if cnt > 1 {
					break
				}
			}
			if cnt == 1 {
				if _, ok := m[pkg]; ok {
					m[pkg]++
				} else {
					log.Fatal("Found New Package") //
				}
			}

		}
	}
	m_out()
}

func m_init() {
	m_cache, err := os.Open("m_cache")
	if err == nil {
		d := gob.NewDecoder(m_cache)
		// Decoding the serialized data
		err = d.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		m = make(map[string]int)
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
			m[outScanner.Text()] = 0
		}

		b := new(bytes.Buffer)
		e := gob.NewEncoder(b)
		// Encoding the map
		err = e.Encode(m)
		if err != nil {
			log.Fatal(err)
		}
		m_cache, err := os.Create("m_cache")
		if err != nil {
			log.Fatal(err)
		}
		defer m_cache.Close()
		m_cache.Write(b.Bytes())
	}

}

func m_out() {
	file, err := os.Create("out.txt")
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString("Not used:\n")
	for k, v := range m {
		if v == 0 {
			file.WriteString(k + "\n")
		}
	}
}
