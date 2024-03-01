package taskExecution

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func StartTasks(execTime int, fileName string, toRun func(int, string) *os.File, toAnalyse func(*os.File)) {

	var wg sync.WaitGroup

	timer := time.After(time.Duration(execTime) * time.Second)

	// канал для сигнала о завершении runBpftrace
	done := make(chan *os.File)

	// горутина для запуска bpftrace
	bpftrace_run := func() {
		defer wg.Done()

		for {
			select {
			case <-timer:
				close(done)
				fmt.Fprintln(os.Stdout, "Final!")
				os.Exit(1)
			default:
				file := toRun(execTime, fileName)
				done <- file
			}
		}
	}

	// горутина для анализа вывода bpftrace
	bpftrace_analyze := func() {
		defer wg.Done()

		for {
			result := <-done
			toAnalyse(result)
			fmt.Fprintln(os.Stdout, "Written.")
		}

	}

	wg.Add(2)
	go bpftrace_run()
	go bpftrace_analyze()
	wg.Wait()
}
