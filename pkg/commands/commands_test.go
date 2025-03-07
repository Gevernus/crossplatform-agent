package commands

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var helperCommandPath string

func init() {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// If we're the helper process, run the requested command.
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		os.Exit(2)
	}
	cmd, args := args[0], args[1:]
	switch cmd {
	case "echo":
		for _, arg := range args {
			os.Stdout.WriteString(arg + " ")
		}
		os.Exit(0)
	// Add more commands as needed
	default:
		os.Exit(2)
	}
}

func helperCommand(t *testing.T, cmd string, args ...string) *exec.Cmd {
	t.Helper()

	if helperCommandPath == "" {
		var err error
		helperCommandPath, err = os.Executable()
		if err != nil {
			t.Fatal(err)
		}
	}

	command := exec.Command(helperCommandPath, append([]string{"--", cmd}, args...)...)
	command.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return command
}

func TestShutdown(t *testing.T) {
	originalOS := currentOS
	originalExecCommand := execCommand
	defer func() {
		currentOS = originalOS
		execCommand = originalExecCommand
	}()

	tests := []struct {
		name        string
		os          string
		expectedCmd string
		expectError bool
	}{
		{"Windows", "windows", "shutdown /s /t 0", false},
		{"Darwin", "darwin", "shutdown -h now", false},
		{"Linux", "linux", "shutdown -h now", false},
		{"Unsupported", "freebsd", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS = tt.os

			var cmd *exec.Cmd
			execCommand = func(name string, arg ...string) *exec.Cmd {
				cmd = helperCommand(t, "echo", name)
				cmd.Args = append(cmd.Args, arg...)
				return cmd
			}

			err := Shutdown()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported operating system")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmd, strings.Join(cmd.Args[3:], " "), "Unexpected command")
			}
		})
	}
}

func TestShutdownNetwork(t *testing.T) {
	originalOS := currentOS
	originalExecCommand := execCommand
	defer func() {
		currentOS = originalOS
		execCommand = originalExecCommand
	}()

	tests := []struct {
		name         string
		os           string
		expectedCmds []string
		expectError  bool
	}{
		{
			"Windows",
			"windows",
			[]string{"netsh interface set interface name=* admin=disabled"},
			false,
		},
		{
			"Darwin",
			"darwin",
			[]string{
				"networksetup -setairportpower en0 off",
				"networksetup -setnetworkserviceenabled Ethernet off",
			},
			false,
		},
		{
			"Linux",
			"linux",
			[]string{"nmcli networking off"},
			false,
		},
		{
			"Unsupported",
			"freebsd",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS = tt.os

			var commands []string
			execCommand = func(name string, arg ...string) *exec.Cmd {
				cmd := helperCommand(t, "echo", name)
				cmd.Args = append(cmd.Args, arg...)
				commands = append(commands, strings.Join(cmd.Args[3:], " "))
				return cmd
			}

			err := ShutdownNetwork()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported operating system")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmds, commands, "Unexpected commands")
			}
		})
	}
}
