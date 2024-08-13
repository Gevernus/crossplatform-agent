package utils

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

// TestOpenDirectory tests the OpenDirectory function
func TestOpenDirectory(t *testing.T) {
	originalOS := currentOS
	defer func() { currentOS = originalOS }()

	tests := []struct {
		name        string
		os          string
		path        string
		expectError bool
	}{
		{"Windows", "windows", "C:\\test", false},
		{"Darwin", "darwin", "/test", false},
		{"Linux", "linux", "/test", false},
		{"Unsupported", "unsupported", "/test", true},
	}
	u := &Impl{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS = tt.os

			// Mock execCommand
			oldExecCommand := execCommand
			defer func() { execCommand = oldExecCommand }()

			execCommand = func(name string, arg ...string) *exec.Cmd {
				return helperCommand(t, "echo", name, strings.Join(arg, " "))
			}

			err := u.OpenDirectory(tt.path)

			if tt.expectError && err == nil {
				t.Errorf("expected an error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestOpenFile tests the OpenFile function
func TestOpenFile(t *testing.T) {
	originalExecCommand := execCommand
	originalOS := currentOS
	defer func() {
		execCommand = originalExecCommand
		currentOS = originalOS
	}()

	tests := []struct {
		name        string
		goos        string
		path        string
		expectedCmd string
		expectError bool
	}{
		{"Windows Success", "windows", "C:\\test.txt", "rundll32 url.dll,FileProtocolHandler C:\\test.txt", false},
		{"Darwin Success", "darwin", "/test.txt", "open /test.txt", false},
		{"Linux Success", "linux", "/test.txt", "xdg-open /test.txt", false},
		{"Unsupported OS", "freebsd", "/test.txt", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOS = tt.goos

			var cmd *exec.Cmd
			execCommand = func(name string, arg ...string) *exec.Cmd {
				cmd = helperCommand(t, "echo", name)
				cmd.Args = append(cmd.Args, arg...)
				return cmd
			}
			u := &Impl{}
			err := u.OpenFile(tt.path)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported platform")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCmd, strings.Join(cmd.Args[3:], " "), "Unexpected command")
			}
		})
	}
}
