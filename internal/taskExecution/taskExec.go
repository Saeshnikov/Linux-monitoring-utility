package taskExecution

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var processes []*exec.Cmd

var hotExit chan bool
var mutex sync.RWMutex

func StartTasks(toExec ...ExecUnit) error {
	var wg sync.WaitGroup
	processes = make([]*exec.Cmd, len(toExec))

	last_con := -1
	errChan := make(chan error, len(toExec))
	for index, unit := range toExec {

		switch v := unit.(type) {
		case execUnitOneShotF:
		case execUnitContinuousF:
			wg.Add(1)
			os.MkdirAll(unit.(execUnitContinuousF).getDir(), os.FileMode(0777))
			go execContinuousF(index, unit.getBinPath(), unit.(execUnitContinuousF).getDir(), unit.getArgs(), unit.getExecCount(), unit.(execUnitContinuousF).getExecTime(), errChan, &wg)
			last_con = index
		case execUnitOneShotC:
		case execUnitContinuousC:
			wg.Add(1)
			go execContinuousC(index, unit.getBinPath(), unit.(execUnitContinuousC).getChan(), unit.getArgs(), unit.getExecCount(), unit.(execUnitContinuousF).getExecTime(), errChan, &wg)
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
			go execOneShotF(index, unit.getBinPath(), unit.(execUnitOneShotF).getDir(), unit.getArgs(), unit.getExecCount(), errChan, &wg)
		case execUnitContinuousF:
		case execUnitOneShotC:
			wg.Add(1)
			go execOneShotC(index, unit.getBinPath(), unit.(execUnitOneShotC).getChan(), unit.getArgs(), unit.getExecCount(), errChan, &wg)
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
