package taskExecution

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ExecUnit struct {
	BinPath      string
	Args         string
	ExecCount    uint
	IsContinuous bool
	ExecTime     time.Duration
}

var processes []*exec.Cmd

func StartTasks(toExec []ExecUnit, outDirPath string) {
	var wg sync.WaitGroup
	processes = make([]*exec.Cmd, len(toExec))
	outToFile := func(filename string, c <-chan bytes.Buffer) {
		defer wg.Done()
		b, ok := <-c
		if !ok {
			return
		}
		tmp := strings.Split(filename, "/")

		filename = tmp[len(tmp)-1]

		file, err := os.Create(outDirPath + "/tmp/" + filename + "." + strconv.FormatInt(time.Now().Unix(), 10))
		if err != nil {
			log.Fatal(err)
		}

		file.Write(b.Bytes())
	}

	execOneShot := func(binPath string, args string, execCount uint) {
		defer wg.Done()
		for execCount_ := 0; execCount_ < int(execCount); execCount_++ {
			buf := make(chan bytes.Buffer, 1)
			fmt.Printf("%s Started...\n", binPath)
			toRunOneShot(binPath, args, buf)
			wg.Add(1)
			go outToFile(binPath, buf)
		}
	}

	execContinuous := func(index int, binPath string, args string, execCount uint, execTime time.Duration) {
		defer wg.Done()
		con_run := func(p chan *exec.Cmd, c chan bytes.Buffer) {
			defer wg.Done()
			toRunContinuous(binPath, args, p, c)

		}
		var prevProc *exec.Cmd
		p := make(chan *exec.Cmd, 1)
		for execCount_ := 0; execCount_ < int(execCount); execCount_++ {
			buf := make(chan bytes.Buffer, 1)
			wg.Add(1)
			go con_run(p, buf)
			fmt.Printf("%s Started...\n", binPath)
			processes[index] = <-p
			if prevProc != nil {
				err := prevProc.Process.Signal(os.Interrupt)
				if err != nil {
					intAllProcesses(err)
				}
				fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)

			}
			prevProc = processes[index]
			wg.Add(1)
			go outToFile(binPath, buf)
			time.Sleep(execTime)
		}

		err := prevProc.Process.Signal(os.Interrupt)
		if err != nil {
			intAllProcesses(err)
		}
		fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)
	}

	for index, unit := range toExec {

		if !unit.IsContinuous {
			wg.Add(1)
			go execOneShot(unit.BinPath, unit.Args, unit.ExecCount)
		} else {
			wg.Add(1)
			go execContinuous(index, unit.BinPath, unit.Args, unit.ExecCount, unit.ExecTime)
		}
	}

	wg.Wait()

}

// poka skip
func intAllProcesses(err error) {
	log.Fatal(err)
}

func toRunOneShot(binPath string, args string, c chan<- bytes.Buffer) {
	var cmd *exec.Cmd
	if args != "" {
		cmd = exec.Command(binPath, args)
	} else {
		cmd = exec.Command(binPath)
	}

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(pipe)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer

	line, err := reader.ReadString('\n')

	for err == nil {
		buffer.WriteString(line)
		line, err = reader.ReadString('\n')
	}
	cmd.Wait()
	c <- buffer
}

func toRunContinuous(binPath string, args string, p chan<- *exec.Cmd, c chan<- bytes.Buffer) {

	var buffer bytes.Buffer

	var cmd *exec.Cmd
	if args != "" {
		cmd = exec.Command(binPath, args)
	} else {
		cmd = exec.Command(binPath)
	}
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(pipe)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Procces is running as pid %d\n", cmd.Process.Pid)
	line, err := reader.ReadString('\n')

	p <- cmd

	for err == nil {
		buffer.WriteString(line)
		line, err = reader.ReadString('\n')
	}

	cmd.Wait()
	c <- buffer
}
