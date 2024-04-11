package taskExecution

import (
	"fmt"
	lsofLayer "linux-monitoring-utility/internal/lsofLayer"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	"os"
	"sync"
	"time"
)

func StartTasks(program_time uint, bpftrace_time uint, fileName string, outputPath string, outputMap *map[string]bool, toRun func(uint, string, string, *map[string]bool, chan *os.Process)) error {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(program_time) * time.Second)
	c := make(chan *os.Process, 1)
	errorChan := make(chan error, 1)
	var curProc *os.Process = nil
	var prevProc *os.Process = nil
	lsof_run := func() {
		fmt.Printf("Lsof started...\n")
		arr, err := lsofLayer.LsofExec()
		if err != nil {
			wg.Wait()
			errorChan <- err
			return
		}

		err = rpmLayer.RPMlayer(arr, outputPath, outputMap)
		if err != nil {
			wg.Wait()
			errorChan <- err
			return
		}
	}

	bpftrace_run := func() {
		defer wg.Done()
		toRun(bpftrace_time, fileName, outputPath, outputMap, c)
	}
	flag := false

	for {
		select {
		case err := <-errorChan:
			return err
		case <-timer:
			err := curProc.Signal(os.Interrupt)
			if err != nil {
				return err
			}
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
				err := prevProc.Signal(os.Interrupt)
				if err != nil {
					return err
				}
				fmt.Printf("Stopping previous process with PID: %d\n", prevProc.Pid)

			}
			prevProc = curProc
			time.Sleep(time.Duration(bpftrace_time) * time.Second)

		}
	}

}
