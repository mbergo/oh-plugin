package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: oh-plugin install <repository-address>")
		os.Exit(1)
	}

	repositoryAddress := os.Args[2]
	zshCustomDir := filepath.Join(os.Getenv("HOME"), ".oh-my-zsh", "custom", "plugins")
	zshrcPath := filepath.Join(os.Getenv("HOME"), ".zshrc")

	// Clone the plugin repository to the custom plugins directory
	cloneCmd := exec.Command("git", "clone", repositoryAddress, zshCustomDir)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to clone repository: %s\n", string(output))
		os.Exit(1)
	}

	// Read the README.md file from the repository
	readmePath := filepath.Join(zshCustomDir, "README.md")
	readme, err := os.Open(readmePath)
	if err != nil {
		fmt.Printf("Failed to open README.md file: %s\n", err)
		os.Exit(1)
	}
	defer readme.Close()

	// Scan the README.md file for the installation instructions
	var installCommands []string
	scanner := bufio.NewScanner(readme)
	inInstallSection := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "## Install") || strings.Contains(line, "## Installation") {
			inInstallSection = true
		} else if inInstallSection && strings.HasPrefix(line, "##") {
			break
		} else if inInstallSection {
			installCommands = append(installCommands, line)
		}
	}

	// Execute the installation commands
	for _, command := range installCommands {
		parts := strings.Fields(command)
		execCmd := exec.Command(parts[0], parts[1:]...)
		if output, err := execCmd.CombinedOutput(); err != nil {
			fmt.Printf("Failed to execute command '%s': %s\n", command, string(output))
			os.Exit(1)
		}
	}

	// Extract the plugin name from the repository address
	parts := strings.Split(repositoryAddress, "/")
	if len(parts) < 2 {
		fmt.Println("Failed to extract plugin name from repository address")
		os.Exit(1)
	}
	pluginName := parts[len(parts)-1]
		// Add the plugin name to the ~/.zshrc file
	zshrc, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Failed to open ~/.zshrc file: %s\n", err)
		os.Exit(1)
	}
	defer zshrc.Close()

	if _, err = zshrc.WriteString("\nplugins=(" + pluginName + ")\n"); err != nil {
		fmt.Printf("Failed to write to ~/.zshrc file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Plugin successfully installed!")
}

