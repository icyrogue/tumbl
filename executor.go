package tumbl

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"path"
)

type executor struct {
	pipe, errPipe io.ReadCloser
	reader        *bufio.Scanner
	lastLogs      []string
	options       Options
}

func NewExecutor(options *Options) *executor {
	return &executor{
		options: *options,
	}
}

func (e *executor) Run(file string) (err error) {
	f := path.Join(e.options.Dst, file)
	cmd := exec.Command("/bin/sh", f)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	e.errPipe, err = cmd.StderrPipe()
	if err != nil {
		return err
	}
	e.reader = bufio.NewScanner(pipe)
	cmd.Start()
	return nil
}

func (e *executor) GetLogs() (logs []string, err error) {
	if e.reader == nil {
		return nil, errors.New("GetLogs: no reader was defined")
	}
	for e.reader.Scan() {
		logs = append(logs, e.reader.Text())
		e.lastLogs = logs
	}
	return logs, e.reader.Err()
}
