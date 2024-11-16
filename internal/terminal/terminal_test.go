package terminal

import (
	"os"
	"testing"
)

func TestShellEnvVarNotSet(t *testing.T) {
	os.Unsetenv(ShellEnvVar)
	expected := DefaultShell
	result := GetShell()
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestShellEnvVarSet(t *testing.T) {
	os.Setenv(ShellEnvVar, "/bin/zsh")
	expected := "/bin/zsh"
	result := GetShell()
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
	os.Unsetenv(ShellEnvVar)
}

func TestPsenvSubShellEnabled(t *testing.T) {
	os.Setenv(SubShellVar, "1")
	result := GetIsPsenvSubShell()
	if !result {
		t.Errorf("expected true, got %v", result)
	}
	os.Unsetenv(SubShellVar)
}

func TestPsenvSubShellDisabled(t *testing.T) {
	os.Unsetenv(SubShellVar)
	result := GetIsPsenvSubShell()
	if result {
		t.Errorf("expected false, got %v", result)
	}
}

func TestNewSubShellFromSubShell(t *testing.T) {
	os.Setenv(SubShellVar, "1")
	_, err := NewSubShell([]string{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	os.Unsetenv(SubShellVar)
}

func TestNewSubShellWithArgs(t *testing.T) {
	os.Unsetenv(SubShellVar)
	cmd, err := NewSubShell([]string{"VAR=value"}, "echo", "hello")
	if err != nil {
		t.Errorf("expected nil, got error %v", err)
	}
	expected := "echo hello"
	if cmd.Args[2] != expected {
		t.Errorf("expected %v, got %v", expected, cmd.Args[2])
	}
}

func TestNewSubShellWithoutArgs(t *testing.T) {
	os.Unsetenv(SubShellVar)
	cmd, err := NewSubShell([]string{"VAR=value"})
	if err != nil {
		t.Errorf("expected nil, got error %v", err)
	}
	expected := GetShell()
	if cmd.Path != expected {
		t.Errorf("expected %v, got %v", expected, cmd.Path)
	}
}
