package taskExecution

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

func toRunOneShot(binPath string, args []string, c chan<- bytes.Buffer, p chan<- *exec.Cmd, errChan chan<- error) {
	var cmd *exec.Cmd
	if args != nil {
		cmd = exec.Command(binPath, args...)
	} else {
		cmd = exec.Command(binPath)
	}

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err
		return
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		errChan <- err
		return
	}

	go func(pipe io.ReadCloser) {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')

		for err == nil {
			fmt.Print(line)
			line, err = reader.ReadString('\n')
		}
	}(errPipe)

	reader := bufio.NewReader(pipe)
	if err := cmd.Start(); err != nil {
		errChan <- err
		return
	}

	p <- cmd

	var buffer bytes.Buffer

	line, err := reader.ReadString('\n')

	for err == nil {
		buffer.WriteString(line)
		line, err = reader.ReadString('\n')
	}
	cmd.Wait()
	c <- buffer
}

func toRunContinuous(binPath string, args []string, p chan<- *exec.Cmd, c chan<- bytes.Buffer, errChan chan<- error) {

	var buffer bytes.Buffer

	var cmd *exec.Cmd
	if args != nil {
		cmd = exec.Command(binPath, args...)
	} else {
		cmd = exec.Command(binPath)
	}
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err
		return
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		errChan <- err
		return
	}

	go func(pipe io.ReadCloser) {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')

		for err == nil {
			fmt.Print(line)
			line, err = reader.ReadString('\n')
		}
	}(errPipe)

	reader := bufio.NewReader(pipe)
	if err := cmd.Start(); err != nil {
		errChan <- err
		return
	}
	p <- cmd

	fmt.Printf("Procces is running as pid %d\n", cmd.Process.Pid)
	line, err := reader.ReadString('\n')

	for err == nil {
		buffer.WriteString(line)
		line, err = reader.ReadString('\n')
	}

	cmd.Wait()
	c <- buffer
}

// function that writing programs output to {outDirPath}/tmp/{binary_filename.timestamp}
func outToFile(outDirPath string, c <-chan bytes.Buffer, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	b, ok := <-c
	if !ok {
		return
	}

	file, err := os.Create(outDirPath + "/" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		errChan <- err
		return
	}

	_, err = file.Write(b.Bytes())
	if err != nil {
		errChan <- err
		return
	}
}

// function that execute one shot for multiple times (doesn't wait for ending previous)
func execOneShotF(index int, binPath string, outDirPath string, args []string, execCount uint, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	ons_run := func(p chan *exec.Cmd, c chan bytes.Buffer, errChan chan error) {
		defer wg.Done()
		toRunOneShot(binPath, args, c, p, errChan)
	}

	for execCount_ := 0; execCount_ < int(execCount); execCount_++ {

		buf := make(chan bytes.Buffer, 1)
		p := make(chan *exec.Cmd, 1)
		errChan_ := make(chan error, 1)

		wg.Add(1)
		go ons_run(p, buf, errChan_)
		var p_ *exec.Cmd

		select {
		case p_ = <-p:
		case err := <-errChan_:
			errChan <- err
			return
		}
		fmt.Printf("%s Started...\n", binPath)
		mutex.Lock()
		processes[index] = p_
		mutex.Unlock()
		wg.Add(1)
		go outToFile(outDirPath, buf, errChan_, wg)
		select {
		case err := <-errChan:
			errChan <- err
			return
		default:
		}
	}
}

func execOneShotC(index int, binPath string, outChan chan chan bytes.Buffer, args []string, execCount uint, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	ons_run := func(p chan *exec.Cmd, c chan bytes.Buffer, errChan chan error) {
		defer wg.Done()
		toRunOneShot(binPath, args, c, p, errChan)
	}

	for execCount_ := 0; execCount_ < int(execCount); execCount_++ {

		buf := make(chan bytes.Buffer, 1)
		p := make(chan *exec.Cmd, 1)
		errChan_ := make(chan error, 1)

		wg.Add(1)
		go ons_run(p, buf, errChan_)
		var p_ *exec.Cmd

		select {
		case p_ = <-p:
		case err := <-errChan_:
			errChan <- err
			return
		}
		fmt.Printf("%s Started...\n", binPath)
		mutex.Lock()
		processes[index] = p_
		mutex.Unlock()
		outChan <- buf

		select {
		case err := <-errChan:
			errChan <- err
			return
		default:
		}
	}
}

// function that execute contionuous for some time multiple times (starts new time is out, then ends previous)
func execContinuousF(index int, binPath string, outDir string, args []string, execCount uint, execTime time.Duration, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	con_run := func(p chan *exec.Cmd, c chan bytes.Buffer, errChan chan error) {
		defer wg.Done()
		toRunContinuous(binPath, args, p, c, errChan)

	}

	var prevProc *exec.Cmd
	p := make(chan *exec.Cmd, 1)
	errChan_ := make(chan error, 1)
	for execCount_ := 0; execCount_ < int(execCount); execCount_++ {
		buf := make(chan bytes.Buffer, 1)
		wg.Add(1)
		go con_run(p, buf, errChan_)
		var p_ *exec.Cmd

		select {
		case p_ = <-p:
		case err := <-errChan_:
			errChan <- err
			return
		}

		fmt.Printf("%s Started...\n", binPath)

		mutex.Lock()
		processes[index] = p_
		mutex.Unlock()
		if prevProc != nil {

			err := prevProc.Process.Signal(os.Interrupt)
			if err != nil {
				p_.Process.Signal(os.Interrupt)
				errChan <- err
				return
			}
			fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)

		}

		mutex.RLock()
		prevProc = processes[index]
		mutex.RUnlock()

		timer := time.After(execTime)
		select {
		case <-buf:
			errChan <- fmt.Errorf("unexpected end of execution %s (%d)", binPath, prevProc.Process.Pid)
			return
		case err := <-errChan_:
			errChan <- err
			return
		case <-timer:
		}
		wg.Add(1)
		go outToFile(outDir, buf, errChan_, wg)
	}

	err := prevProc.Process.Signal(os.Interrupt)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)
}

func execContinuousC(index int, binPath string, outChan chan chan bytes.Buffer, args []string, execCount uint, execTime time.Duration, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	con_run := func(p chan *exec.Cmd, c chan bytes.Buffer, errChan chan error) {
		defer wg.Done()
		toRunContinuous(binPath, args, p, c, errChan)

	}

	var prevProc *exec.Cmd
	p := make(chan *exec.Cmd, 1)
	errChan_ := make(chan error, 1)
	for execCount_ := 0; execCount_ < int(execCount); execCount_++ {
		buf := make(chan bytes.Buffer, 1)
		wg.Add(1)
		go con_run(p, buf, errChan_)
		var p_ *exec.Cmd

		select {
		case p_ = <-p:
		case err := <-errChan_:
			errChan <- err
			return
		}

		fmt.Printf("%s Started...\n", binPath)

		mutex.Lock()
		processes[index] = p_
		mutex.Unlock()
		if prevProc != nil {

			err := prevProc.Process.Signal(os.Interrupt)
			if err != nil {
				p_.Process.Signal(os.Interrupt)
				errChan <- err
				return
			}
			fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)

		}

		mutex.RLock()
		prevProc = processes[index]
		mutex.RUnlock()

		timer := time.After(execTime)
		select {
		case <-buf:
			errChan <- fmt.Errorf("unexpected end of execution %s (%d)", binPath, prevProc.Process.Pid)
			return
		case err := <-errChan_:
			errChan <- err
			return
		case <-timer:
		}
		outChan <- buf
	}

	err := prevProc.Process.Signal(os.Interrupt)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)
}
