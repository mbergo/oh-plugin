package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestInstallPlugin(t *testing.T) {
	// Prepare test inputs and expected outputs
	input := "https://github.com/some-plugin"
	expectedOutput := "Plugin installed successfully"

	// Set up a fake file system for testing
	oldExecCommand := exec.Command
	defer func() { exec.Command = oldExecCommand }()
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{
			Stdout: &bytes.Buffer{},
			Stderr: &bytes.Buffer{},
		}
	}

	oldReadFile := ioutil.ReadFile
	defer func() { ioutil.ReadFile = oldReadFile }()
	ioutil.ReadFile = func(filename string) ([]byte, error) {
		return []byte(""), nil
	}

	oldWriteFile := ioutil.WriteFile
	defer func() { ioutil.WriteFile = oldWriteFile }()
	ioutil.WriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return nil
	}

	oldOpenFile := os.OpenFile
	defer func() { os.OpenFile = oldOpenFile }()
	os.OpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return &os.File{}, nil
	}

	oldNewScanner := bufio.NewScanner
	defer func() { bufio.NewScanner = oldNewScanner }()
	bufio.NewScanner = func(r io.Reader) *bufio.Scanner {
		return &bufio.Scanner{}
	}

	// Call the function being tested
	output := installPlugin(input)

	// Compare the actual output to the expected output
	if output != expectedOutput {
		t.Errorf("Test failed: expected output %q, got %q", expectedOutput, output)
	}
}

