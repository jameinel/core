// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package exec

import (
	"bytes"
	"os/exec"
	"syscall"

	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("juju.util.exec")

// Parameters for RunCommands.  Commands contains one or more commands to be
// executed using '/bin/bash -s'.  If WorkingDir is set, this is passed
// through to bash.  Similarly if the Environment is specified, this is used
// for executing the command.
type RunParams struct {
	Commands    string
	WorkingDir  string
	Environment []string
}

// ExecResponse contains the return code and output generated by executing a
// command.
type ExecResponse struct {
	Code   int
	Stdout []byte
	Stderr []byte
}

// RunCommands executes the Commands specified in the RunParams using
// '/bin/bash -s', passing the commands through as stdin, and collecting
// stdout and stderr.  If a non-zero return code is returned, this is
// collected as the code for the response and this does not classify as an
// error.
func RunCommands(run RunParams) (*ExecResponse, error) {
	ps := exec.Command("/bin/bash", "-s")
	if run.Environment != nil {
		ps.Env = run.Environment
	}
	if run.WorkingDir != "" {
		ps.Dir = run.WorkingDir
	}
	ps.Stdin = bytes.NewBufferString(run.Commands)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	ps.Stdout = stdout
	ps.Stderr = stderr

	err := ps.Start()
	if err == nil {
		err = ps.Wait()
	}
	result := &ExecResponse{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	if ee, ok := err.(*exec.ExitError); ok && err != nil {
		status := ee.ProcessState.Sys().(syscall.WaitStatus)
		if status.Exited() {
			// A non-zero return code isn't considered an error here.
			result.Code = status.ExitStatus()
			err = nil
		}
		logger.Infof("run result: %v", ee)
	}
	return result, err
}
