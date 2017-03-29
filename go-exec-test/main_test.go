package main

import (
	"os"
	"os/exec"
	"testing"
)

var mockSuccessCmd = exec.Command(os.Args[0], "-test.run=TestExecSuccess")
var mockErrorCmd = &exec.Cmd{Process: &os.Process{}}
var testCmd = mockSuccessCmd

func init() {
	cmd = func(_ string, _ ...string) *exec.Cmd { return testCmd }
}

func TestMain(_ *testing.T) {
	main()
	testCmd = mockErrorCmd
	main()
}

func TestExecSuccess(_ *testing.T) {}
