package taskExecution

import (
	"bytes"
	"time"
)

type ExecUnit interface {
	getBinPath() string
	getArgs() string
	getExecCount() uint
}

// to channel
type execUnitOneShotC struct {
	binPath   string
	args      string
	execCount uint
	ch        chan chan bytes.Buffer
}

type execUnitContinuousC struct {
	execUnitOneShotC
	execTime time.Duration
}

func NewExecUnitContinuousC(binPath string, args string, execCount uint, execTime time.Duration, ch chan chan bytes.Buffer) *execUnitContinuousC {
	execUnitOneShotC := execUnitOneShotC{binPath: binPath, args: args, execCount: execCount, ch: ch}
	return &execUnitContinuousC{execUnitOneShotC: execUnitOneShotC, execTime: execTime}
}

func NewExecUnitOneShotC(binPath string, args string, execCount uint, ch chan chan bytes.Buffer) *execUnitOneShotC {
	return &execUnitOneShotC{binPath: binPath, args: args, execCount: execCount, ch: ch}
}

func (t execUnitOneShotC) getBinPath() string {
	return t.binPath
}

func (t execUnitOneShotC) getArgs() string {
	return t.args
}

func (t execUnitOneShotC) getExecCount() uint {
	return t.execCount
}

func (t execUnitOneShotC) getChan() chan chan bytes.Buffer {
	return t.ch
}

func (t execUnitContinuousC) getExecTime() time.Duration {
	return t.execTime
}

// to file
type execUnitOneShotF struct {
	binPath   string
	args      string
	execCount uint
	dir       string
}

type execUnitContinuousF struct {
	execUnitOneShotF
	execTime time.Duration
}

func NewExecUnitContinuousF(binPath string, args string, execCount uint, execTime time.Duration, dir string) *execUnitContinuousF {
	execUnitOneShotF := execUnitOneShotF{binPath: binPath, args: args, execCount: execCount, dir: dir}
	return &execUnitContinuousF{execUnitOneShotF: execUnitOneShotF, execTime: execTime}
}

func NewExecUnitOneShotF(binPath string, args string, execCount uint, dir string) *execUnitOneShotF {
	return &execUnitOneShotF{binPath: binPath, args: args, execCount: execCount, dir: dir}
}

func (t execUnitOneShotF) getBinPath() string {
	return t.binPath
}

func (t execUnitOneShotF) getArgs() string {
	return t.args
}

func (t execUnitOneShotF) getExecCount() uint {
	return t.execCount
}
func (t execUnitOneShotF) getDir() string {
	return t.dir
}

func (t execUnitContinuousF) getExecTime() time.Duration {
	return t.execTime
}
