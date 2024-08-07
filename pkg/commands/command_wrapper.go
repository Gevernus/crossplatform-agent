package commands

import (
	"os/exec"
)

// CommandRunner defines a function type for running commands
type CommandRunner func(name string, arg ...string) *exec.Cmd

// DefaultCommandRunner is the default implementation using exec.Command
var DefaultCommandRunner CommandRunner = exec.Command
