package taskExecution

import (
	lsofLayer "linux-monitoring-utility/internal/lsofLayer"
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	"sync"
	"time"
)

func StartTasks(program_time int, bpftrace_time int, fileName string, toRun func(int, string)) error {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(program_time) * time.Second)

	bpftrace_run := func() {
		defer wg.Done()
		toRun(bpftrace_time, fileName)
	}
	wg.Add(1)
	go bpftrace_run()

	arr, err := lsofLayer.LsofExec()
	if err != nil {
		wg.Wait()
		return err
	}

	err = rpmLayer.RPMlayer(arr)
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
