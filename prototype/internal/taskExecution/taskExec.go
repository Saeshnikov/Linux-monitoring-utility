package taskExecution

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func StartTasks(program_time uint, bpftrace_time uint, fileName string, toRun func(string, chan *exec.Cmd), toRunLsof func()) error {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(program_time) * time.Second)
	c := make(chan *exec.Cmd, 1)
	var curProc *exec.Cmd = nil
	var prevProc *exec.Cmd = nil
	lsof_run := func() {
		defer wg.Done()
		fmt.Printf("Lsof started...\n")
		toRunLsof()
	}

	bpftrace_run := func() {
		defer wg.Done()
		toRun(fileName, c)
	}
	flag := false

	for {
		select {
		case <-timer:
			err := curProc.Process.Signal(os.Interrupt)
			if err != nil {
				return err
			}
			fmt.Printf("Stopping previous process with PID: %d\n", curProc.Process.Pid)
			wg.Wait()
			return nil
		default:
			wg.Add(1)
			go bpftrace_run()
			curProc = <-c
			if !flag {
				wg.Add(1)
				go lsof_run()
				flag = true
			}
			if prevProc != nil {
				err := prevProc.Process.Signal(os.Interrupt)
				if err != nil {
					return err
				}
				fmt.Printf("Stopping previous process with PID: %d\n", prevProc.Process.Pid)

			}
			prevProc = curProc
			time.Sleep(time.Duration(bpftrace_time) * time.Second)

		}
	}

}
