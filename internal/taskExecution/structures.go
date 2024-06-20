package taskExecution

import "time"

type ExecUnit interface {
	getBinPath() string
	getArgs() string
	getExecCount() uint
}

type execUnitOneShot struct {
	binPath   string
	args      string
	execCount uint
}

type execUnitContinuous struct {
	execUnitOneShot
	execTime time.Duration
}

func NewExecUnitContinuous(binPath string, args string, execCount uint, execTime time.Duration) *execUnitContinuous {
	ExecUnitOneShot := execUnitOneShot{binPath: binPath, args: args, execCount: execCount}
	return &execUnitContinuous{execUnitOneShot: ExecUnitOneShot, execTime: execTime}
}

func NewExecUnitOneShot(binPath string, args string, execCount uint) *execUnitOneShot {
	return &execUnitOneShot{binPath: binPath, args: args, execCount: execCount}
}

func (t execUnitOneShot) getBinPath() string {
	return t.binPath
}

func (t execUnitOneShot) getArgs() string {
	return t.args
}

func (t execUnitOneShot) getExecCount() uint {
	return t.execCount
}

func (t execUnitContinuous) getExecTime() time.Duration {
	return t.execTime
}
