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
	"sync"
	"time"
)

var TIME_BPFTRACE = 100
var m map[string]int

func main() {

	var wg sync.WaitGroup
	m_init()

	timer := time.After(time.Duration(TIME_BPFTRACE) * time.Second)

	// канал для сигнала о завершении runBpftrace
	done := make(chan struct{})

	// горутина для запуска bpftrace
	bpftrace_run := func() {
		defer wg.Done()

		for {
			select {
			case <-timer:
				close(done)
				fmt.Fprintln(os.Stdout, "Final!")
				os.Remove("tmp.txt")
				os.Exit(1)
			default:
				runBpftrace()
				done <- struct{}{}
			}
		}
	}

	// горутина для анализа вывода bpftrace
	bpftrace_analyze := func() {
		defer wg.Done()

		for {
			<-done
			analyzeBpftraceOutput()
			m_out()
			fmt.Fprintln(os.Stdout, "Written.")
		}

	}

	wg.Add(2)
	go bpftrace_run()
	go bpftrace_analyze()
	wg.Wait()
}

func runBpftrace() {

	cmdToRun := "/usr/bin/bpftrace"
	args := []string{"", "ex.bt"}
	procAttr := new(os.ProcAttr)

	// Создание временного файла для вывода bpftrace
	file, err := os.Create("tmp.txt")
	if err != nil {
		log.Fatal(err)
	}
	procAttr.Files = []*os.File{os.Stdin, file, os.Stderr}

	// Запуск bpftrace
	if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
		fmt.Printf("ERROR Unable to run %s: %s\n", cmdToRun, err.Error())
	} else {
		fmt.Printf("%s running as pid %d\n", cmdToRun, process.Pid)
		time.Sleep(time.Duration(TIME_BPFTRACE) * time.Second / 10) // скрипт перезапускается 10 раз
		process.Signal(os.Interrupt)
	}

}

func analyzeBpftraceOutput() {
	file, err := os.Open("tmp.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	r, err := regexp.Compile(`@filename\[(.*?)\]`)
	if err != nil {
		log.Fatal(err)
	}

	for fileScanner.Scan() {

		res := r.FindAllStringSubmatch(fileScanner.Text(), -1)
		if res != nil {

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
