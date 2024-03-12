package taskExecution

import (
	"sync"
	"time"
)

func StartTasks(program_time int, bpftrace_time int, fileName string, toRun func(int, string)) {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(program_time) * time.Second)

	bpftrace_run := func() {
		defer wg.Done()
		toRun(bpftrace_time, fileName)
	}

	for {
		select {
		case <-timer:
			wg.Wait()
			return
		default:
			wg.Add(1)
			go bpftrace_run()
			time.Sleep(time.Duration(bpftrace_time)*time.Second - 2*time.Second)
		}
	}

}
