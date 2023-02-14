package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	repositoryAddress := "https://github.com/example/zsh-plugin"
	zshCustomDir := filepath.Join(os.Getenv("HOME"), ".oh-my-zsh", "custom", "plugins")
	zshrcPath := filepath.Join(os.Getenv("HOME"), ".zshrc")

	// Mocked git clone command
	cloneCalled := false
	execCommand = func(command string, args ...string) *exec.Cmd {
		if command == "git" && args[0] == "clone" {
			cloneCalled = true
			return exec.Command("echo", "Cloning repository...")
		}
		return exec.Command(command, args...)
	}

	// Mocked ~/.zshrc file content
	zshrcContent := ""
	readFile = func(filename string) ([]byte, error) {
		return []byte(zshrcContent), nil
	}
	writeFile = func(filename string, data []byte, perm os.FileMode) error {
		zshrcContent = string(data)
		return nil
	}

	// Mocked README.md file content
	readmeContent := "## Installation\ninstall-command1\ninstall-command2\n"
	openFile = func(filename string, flag int, perm os.FileMode) (*os.File, error) {
		file := &os.File{}
		file.Read = func(p []byte) (int, error) {
			return strings.NewReader(readmeContent).Read(p)
		}
		file.Close = func() error { return nil }
		return file, nil
	}

	os.Args = []string{"oh-plugin", "install", repositoryAddress}
	main()

	// Check if the git clone command was called
	if !cloneCalled {
		t.Error("git clone command was not called")
	}

	// Check if the plugin was added to the ~/.zshrc file
	if !strings.Contains(zshrcContent, "plugins=(zsh-plugin)\n") {
		t.Error("Plugin was not added to the ~/.zshrc file")
	}
}

var execCommand func(command string, args ...string) *exec.Cmd

func TestMain(m *testing.M) {
	execCommand = exec.Command
	defer func() { execCommand = exec.Command }()
	exec.Command = func(command string, args ...string) *exec.Cmd {
		return execCommand(command, args...)
	}

	os.Exit(m.Run())
}

var writeFile func(filename string, data []byte, perm os.FileMode) error

func init() {
	writeFile = ioutil.WriteFile
	ioutil.WriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return writeFile(filename, data, perm)
	}
}

var openFile func(filename string, flag int, perm os.FileMode) (*os.File, error)

func init() {
	openFile = os.OpenFile
	os.OpenFile = func(filename string, flag int, perm os.FileMode) (*os.File, error) {
		return openFile(filename, flag, perm)
	}
}

