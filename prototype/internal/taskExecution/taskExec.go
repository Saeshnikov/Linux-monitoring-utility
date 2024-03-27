package taskExecution

import (
	lsofLayer "linux-monitoring-utility/internal/lsofLayer"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	"sync"
	"time"
)

func StartTasks(program_time int, bpftrace_time int, fileName string, outputPath string, outputMap *map[string]bool, toRun func(int, string, string, *map[string]bool)) error {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(program_time) * time.Second)

	bpftrace_run := func() {
		defer wg.Done()
		toRun(bpftrace_time, fileName, outputPath, outputMap)
	}
	wg.Add(1)
	go bpftrace_run()

	arr, err := lsofLayer.LsofExec()
	if err != nil {
		wg.Wait()
		return err
	}

	err = rpmLayer.RPMlayer(arr, outputPath, outputMap)
	if err != nil {
		wg.Wait()
		return err
	}

	for {
		select {
		case <-timer:
			wg.Wait()
			return nil
		default:
			wg.Add(1)
			time.Sleep(time.Duration(bpftrace_time)*time.Second - 2*time.Second)
			go bpftrace_run()
		}
	}

}
