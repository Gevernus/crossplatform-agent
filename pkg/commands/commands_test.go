package commands

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCmd is a mock implementation of exec.Cmd
type MockCmd struct {
	mock.Mock
	Path string
	Args []string
}

func (m *MockCmd) Run() error {
	args := m.Called()
	return args.Error(0)
}

var mockCmdInstance *MockCmd

// mockCommand is a helper function to create a mocked command
func mockCommand(name string, arg ...string) *exec.Cmd {
	mockCmdInstance = new(MockCmd)
	mockCmdInstance.On("Run").Return(nil)
	mockCmdInstance.Path = name
	mockCmdInstance.Args = append([]string{name}, arg...)
	cmd := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
	}
	// Use the cmd variable to avoid "declared and not used" error
	_ = cmd
	return cmd
}

// mockOS is a helper function to mock the OS getter
func mockOS(os string) OSGetter {
	return func() string {
		return os
	}
}

func TestShutdown(t *testing.T) {
	// Save the real DefaultCommandRunner and DefaultOSGetter
	realCommandRunner := DefaultCommandRunner
	realOSGetter := DefaultOSGetter

	// Replace DefaultCommandRunner and DefaultOSGetter with mock functions during the test
	DefaultCommandRunner = mockCommand
	defer func() {
		DefaultCommandRunner = realCommandRunner
		DefaultOSGetter = realOSGetter
	}()

	tests := []struct {
		os  string
		cmd []string
	}{
		{os: "windows", cmd: []string{"shutdown", "/s", "/t", "0"}},
		{os: "darwin", cmd: []string{"shutdown", "-h", "now"}},
		{os: "linux", cmd: []string{"shutdown", "-h", "now"}},
	}

	for _, test := range tests {
		t.Run(test.os, func(t *testing.T) {
			DefaultOSGetter = mockOS(test.os)
			err := Shutdown()
			assert.NoError(t, err)
			assert.Equal(t, test.cmd[0], mockCmdInstance.Path)
			assert.ElementsMatch(t, test.cmd, mockCmdInstance.Args)
			mockCmdInstance.AssertExpectations(t)
		})
	}
}

func TestShutdownNetwork(t *testing.T) {
	// Save the real DefaultCommandRunner and DefaultOSGetter
	realCommandRunner := DefaultCommandRunner
	realOSGetter := DefaultOSGetter

	// Replace DefaultCommandRunner and DefaultOSGetter with mock functions during the test
	DefaultCommandRunner = mockCommand
	defer func() {
		DefaultCommandRunner = realCommandRunner
		DefaultOSGetter = realOSGetter
	}()

	tests := []struct {
		os   string
		cmds [][]string
	}{
		{os: "windows", cmds: [][]string{{"netsh", "interface", "set", "interface", "name=*", "admin=disabled"}}},
		{os: "darwin", cmds: [][]string{{"networksetup", "-setairportpower", "en0", "off"}, {"networksetup", "-setnetworkserviceenabled", "Ethernet", "off"}}},
		{os: "linux", cmds: [][]string{{"nmcli", "networking", "off"}}},
	}

	for _, test := range tests {
		t.Run(test.os, func(t *testing.T) {
			DefaultOSGetter = mockOS(test.os)
			for _, cmd := range test.cmds {
				mockCmdInstance = new(MockCmd)
				mockCmdInstance.On("Run").Return(nil)
				DefaultCommandRunner = func(name string, arg ...string) *exec.Cmd {
					mockCmdInstance.Path = name
					mockCmdInstance.Args = append([]string{name}, arg...)
					cmd := &exec.Cmd{
						Path: name,
						Args: append([]string{name}, arg...),
					}
					// Use the cmd variable to avoid "declared and not used" error
					_ = cmd
					return cmd
				}
				err := ShutdownNetwork()
				assert.NoError(t, err)
				assert.Equal(t, cmd[0], mockCmdInstance.Path)
				assert.ElementsMatch(t, cmd, mockCmdInstance.Args)
				mockCmdInstance.AssertExpectations(t)
			}
		})
	}
}
