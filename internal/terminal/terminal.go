package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ShellEnvVar        = "SHELL"
	DefaultShell       = "/bin/bash"
	SubShellVar        = "PSENV_SUBSHELL"
	SubShellVarEnabled = "PSENV_SUBSHELL=1"
)

type SubShell struct {
	CommandArgs []string
}

func GetShell() string {
	shell := os.Getenv(ShellEnvVar)
	if shell == "" {
		shell = DefaultShell
	}
	return shell
}

func GetIsPsenvSubShell() bool {
	return os.Getenv(SubShellVar) == "1"
}

func NewSubShell(envVars []string, args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	if GetIsPsenvSubShell() {
		return nil, fmt.Errorf("cannot create a subshell from a subshell. Exit the current subshell by typing exit and try again")
	}

	shell := GetShell()
	if len(args) > 0 {
		cmd = exec.Command(shell, "-c", strings.Join(args, " "))
	} else {
		cmd = exec.Command(shell)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(), envVars...)
	cmd.Env = append(cmd.Env, SubShellVarEnabled)
	return cmd, nil
}
