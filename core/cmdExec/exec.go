package cmdExec

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

var stdOutSync sync.Map
var stdErrSync sync.Map

func workerFlushRead(isStdOut bool, executionID string, reader io.ReadCloser) {
	// read 1024 bytes at a time
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			reader.Close()
			continue
		}
		if n > 0 {
			if isStdOut {
				dataStream := string(buf[:n])
				obj, ok := stdOutSync.Load(executionID)
				var existing string
				if ok {
					existing = obj.(string)
				}
				stdOutSync.Store(executionID, existing+dataStream)
			} else {
				dataStream := string(buf[:n])
				obj, ok := stdErrSync.Load(executionID)
				var existing string
				if ok {
					existing = obj.(string)
				}
				stdErrSync.Store(executionID, existing+dataStream)
			}
			fmt.Print()
		}
	}
}

// Run invokes a command procedurally and waits to complete
func Run(command string, cmdAndArguments ...string) (stdOut string, stdErr string, err error) {
	executionID := time.Now().String()
	executionID += command
	for _, x := range cmdAndArguments {
		executionID += x
	}
	cmd := exec.Command(command, cmdAndArguments...)

	reader, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.New("Could not open cmd.StdoutPipe(): " + err.Error())
		return
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		err = errors.New("Could not open cmd.StderrPipe(): " + err.Error())
		return
	}
	go workerFlushRead(true, executionID, reader)
	go workerFlushRead(false, executionID, errReader)

	err = cmd.Start()
	if err != nil {
		err = errors.New("Could not cmd.Start(): " + err.Error())
		return
	}

	err = cmd.Wait()
	obj, ok := stdOutSync.Load(executionID)
	if ok {
		stdOut = obj.(string)
	}
	obj, ok = stdErrSync.Load(executionID)
	if ok {
		stdErr = obj.(string)
	}
	if err != nil {
		err = errors.New("Could not cmd.Wait(): " + err.Error())
		return
	}
	return
}
