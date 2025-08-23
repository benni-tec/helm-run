package src

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Runner string

const (
	Shell              Runner = "sh"
	Bash                      = "bash"
	Fish                      = "fish"
	Zsh                       = "zsh"
	Powershell                = "pwsh"
	WindowsPowershell         = "powershell"
	WindowsCommandLine        = "cmd"
)

var Runners = []Runner{
	Shell,
	Bash,
	Fish,
	Zsh,
	Powershell,
	WindowsPowershell,
	WindowsCommandLine,
}

var ShellRunners = []Runner{
	Shell,
	Bash,
	Fish,
	Zsh,
}

var PowershellRunners = []Runner{
	Powershell,
	WindowsPowershell,
}

func (r Runner) Run(name, script string, scriptArgs []string) error {
	// create a temporary file
	extension, err := r.Extension()
	if err != nil {
		return err
	}

	script = r.PrefixScript() + "\n\n" + script

	tmp, err := createTempWithContent("helm-run_"+CurrentContext.Chart.Name+"_"+name+"-*"+extension, script)
	if err != nil {
		return err
	}

	// build the Runner arguments
	cmd := exec.Command(string(r), r.args(tmp, scriptArgs)...)

	if CurrentContext.Verbose {
		fmt.Printf("# Script File: %s\n", tmp)
		fmt.Printf("# Execute: %s\n", cmd.String())
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createTempWithContent(pattern, content string) (string, error) {
	tmp, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}

	defer tmp.Close()

	// write the script to the file
	_, err = tmp.WriteString(content)
	if err != nil {
		return tmp.Name(), err
	}

	// make it executable
	err = os.Chmod(tmp.Name(), 0755)
	if err != nil {
		return tmp.Name(), err
	}

	return tmp.Name(), nil
}

// TODO: this might need be adjusted for different runners
func (r Runner) args(tmp string, scriptArgs []string) []string {
	var runnerArgs []string
	if strings.TrimSpace(CurrentContext.RunnerArgs) != "" {
		runnerArgs = append(runnerArgs, CurrentContext.RunnerArgs)
	}

	// prefix to execute a file
	switch r {
	case WindowsCommandLine:
		runnerArgs = append(runnerArgs, "/c")
	case Powershell:
		runnerArgs = append(runnerArgs, "-File")
	}

	runnerArgs = append(runnerArgs, tmp)
	runnerArgs = append(runnerArgs, scriptArgs...)

	return runnerArgs
}

func (r Runner) PrefixScript() string {
	switch r {
	case Shell:
		return "#!/bin/sh"
	case Bash:
		return "#!/bin/bash"
	case Fish:
		return "#!/bin/fish"
	case Zsh:
		return "#!/bin/zsh"
	case WindowsCommandLine:
		return "@echo off"
	}

	return ""
}

func (r Runner) Extension() (string, error) {
	switch r {
	case Shell:
		fallthrough
	case Bash:
		fallthrough
	case Fish:
		fallthrough
	case Zsh:
		return ".sh", nil

	case Powershell:
		fallthrough
	case WindowsPowershell:
		return ".ps1", nil

	case WindowsCommandLine:
		return ".bat", nil
	}

	return "", errors.New("Runner.Extension: unknown Runner " + string(r))
}
