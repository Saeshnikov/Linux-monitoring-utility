package taskExecution

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func StartTasks(program_time uint, bpftrace_time uint, fileName string, toRun func(string) *os.Process, toRunLsof func()) error {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(program_time) * time.Second)
	c := make(chan *os.Process, 1)
	var curProc *os.Process = nil
	var prevProc *os.Process = nil
	lsof_run := func() {
		fmt.Printf("Lsof started...\n")
		toRunLsof()
	}

	bpftrace_run := func() {
		defer wg.Done()
		c <- toRun(fileName)
	}
	flag := false

	for {
		select {
		case <-timer:
			curProc.Signal(os.Interrupt)
			fmt.Printf("Stopping previous process with PID: %d\n", curProc.Pid)
			wg.Wait()
			return nil
		default:
			wg.Add(1)
			go bpftrace_run()
			curProc = <-c
			if !flag {
				go lsof_run()
				flag = true
			}
			if prevProc != nil {
				prevProc.Signal(os.Interrupt)
				fmt.Printf("Stopping previous process with PID: %d\n", prevProc.Pid)

			}
			prevProc = curProc
			time.Sleep(time.Duration(bpftrace_time) * time.Second)

		}
	}

}
