package regression

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

type Executor struct {
	command  string
	args     []string
	out      string
	executed bool

	// metrics
	rusage *syscall.Rusage
	wall   time.Duration
}

var ErrNotRun = fmt.Errorf("command still was not executed")
var ErrRusageNotAvailable = fmt.Errorf("rusage information not available")

func NewExecutor(command string, args ...string) (*Executor, error) {
	return &Executor{
		command: command,
		args:    args,
	}, nil
}

func (e *Executor) Run() error {
	defer func() { e.executed = true }()

	cmd := exec.Command(e.command, e.args...)

	start := time.Now()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	e.wall = time.Since(start)

	e.out = string(out)

	rusage, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	if ok {
		e.rusage = rusage
	}

	return nil
}

func (e *Executor) Out() (string, error) {
	if !e.executed {
		return "", ErrNotRun
	}

	return e.out, nil
}

func (e *Executor) Rusage() (*syscall.Rusage, error) {
	if !e.executed {
		return nil, ErrNotRun
	}

	if e.rusage == nil {
		return nil, ErrRusageNotAvailable
	}

	return e.rusage, nil
}

func (e *Executor) Wall() (time.Duration, error) {
	if !e.executed {
		return 0 * time.Second, ErrNotRun
	}

	return e.wall, nil
}
