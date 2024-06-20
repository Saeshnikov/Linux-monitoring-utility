package taskExecution

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var processes []*exec.Cmd

var hotExit chan bool
var mutex sync.RWMutex

func StartTasks(outDirPath string, toExec ...ExecUnit) error {
	var wg sync.WaitGroup
	processes = make([]*exec.Cmd, len(toExec))
	//function that writing programs output to {outDirPath}/tmp/{binary_filename.timestamp}
	outToFile := func(outDirPath string, c <-chan bytes.Buffer, errChan chan<- error) {
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

		n, err := file.Write(b.Bytes())
		if err != nil {
			errChan <- err
			return
		}
		fmt.Println(n)
	}

	//function that execute one shot for multiple times (doesn't wait for ending previous)
	execOneShotF := func(index int, binPath string, outDirPath string, args string, execCount uint, errChan chan error) {
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
			go outToFile(outDirPath, buf, errChan_)
			select {
			case err := <-errChan:
				errChan <- err
				return
			default:
			}
		}
	}

	execOneShotC := func(index int, binPath string, outChan chan chan bytes.Buffer, args string, execCount uint, errChan chan error) {
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
			outChan <- buf
			select {
			case err := <-errChan:
				errChan <- err
				return
			default:
			}
		}
	}

	//function that execute contionuous for some time multiple times (starts new time is out, then ends previous)
	execContinuousF := func(index int, binPath string, outDir string, args string, execCount uint, execTime time.Duration, errChan chan error) {
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
			go outToFile(outDir, buf, errChan_)
		}

		err := prevProc.Process.Signal(os.Interrupt)
		if err != nil {
			errChan <- err
			return
		}
		fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)
	}

	execContinuousC := func(index int, binPath string, outChan chan chan bytes.Buffer, args string, execCount uint, execTime time.Duration, errChan chan error) {
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
			outChan <- buf
		}

		err := prevProc.Process.Signal(os.Interrupt)
		if err != nil {
			errChan <- err
			return
		}
		fmt.Printf("Stopping %s process with PID: %d\n", binPath, prevProc.Process.Pid)
	}

	var last_con int
	last_con = -1
	errChan := make(chan error, len(toExec))
	for index, unit := range toExec {

		switch v := unit.(type) {
		case execUnitOneShotF:
		case execUnitContinuousF:
			wg.Add(1)
			os.MkdirAll(unit.(execUnitContinuousF).getDir(), os.FileMode(0777))
			go execContinuousF(index, unit.getBinPath(), unit.(execUnitContinuousF).getDir(), unit.getArgs(), unit.getExecCount(), unit.(execUnitContinuousF).getExecTime(), errChan)
			last_con = index
		case execUnitOneShotC:
		case execUnitContinuousC:
			wg.Add(1)
			go execContinuousC(index, unit.getBinPath(), unit.(execUnitContinuousC).getChan(), unit.getArgs(), unit.getExecCount(), unit.(execUnitContinuousF).getExecTime(), errChan)
			last_con = index
		default:
			return fmt.Errorf("recieved unexpected type : %s", v)
		}

	}

	if last_con != -1 {
	wait_for_con:
		for {
			select {
			case err := <-errChan:
				return err
			default:
				mutex.RLock()
				if processes[last_con] != nil {
					mutex.RUnlock()
					break wait_for_con
				}
				mutex.RUnlock()
			}
		}
	}

	for index, unit := range toExec {

		switch v := unit.(type) {
		case execUnitOneShotF:
			wg.Add(1)
			os.MkdirAll(unit.(execUnitOneShotF).getDir(), os.FileMode(0777))
			go execOneShotF(index, unit.getBinPath(), unit.(execUnitOneShotF).getDir(), unit.getArgs(), unit.getExecCount(), errChan)
		case execUnitContinuousF:
		case execUnitOneShotC:
			wg.Add(1)
			go execOneShotC(index, unit.getBinPath(), unit.(execUnitOneShotC).getChan(), unit.getArgs(), unit.getExecCount(), errChan)
		case execUnitContinuousC:
		default:
			return fmt.Errorf("recieved unexpected type : %s", v)
		}
	}
	waitCh := make(chan bool, 1)
	go func() {
		wg.Wait()
		waitCh <- true
	}()

	select {
	case <-waitCh:
		return nil
	case err := <-errChan:
		if hotExit != nil && <-hotExit {
			fmt.Println("HOT EXIT!")
			<-waitCh
			return nil
		} else {
			IntAllProcesses()
			<-waitCh
			return err
		}

	}
}

func IntAllProcesses() error {
	hotExit = make(chan bool, 1)
	hotExit <- true
	mutex.RLock()
	for _, cmd := range processes {
		if cmd != nil {
			cmd.Process.Signal(os.Interrupt)

		}
	}
	mutex.RUnlock()
	close(hotExit)
	return nil
}

func toRunOneShot(binPath string, args string, c chan<- bytes.Buffer, p chan<- *exec.Cmd, errChan chan<- error) {
	var cmd *exec.Cmd
	if args != "" {
		cmd = exec.Command(binPath, args)
	} else {
		cmd = exec.Command(binPath)
	}

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err
		return
	}

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

func toRunContinuous(binPath string, args string, p chan<- *exec.Cmd, c chan<- bytes.Buffer, errChan chan<- error) {

	var buffer bytes.Buffer

	var cmd *exec.Cmd
	if args != "" {
		cmd = exec.Command(binPath, args)
	} else {
		cmd = exec.Command(binPath)
	}
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err
		return
	}

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
