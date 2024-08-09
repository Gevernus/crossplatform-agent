package service

import (
	"crossplatform-agent/internal/config"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock implementation of the Config struct
type MockConfig struct {
	ConfigPath string
}

// MockExec is a mock for exec.Command
type MockExec struct {
	mock.Mock
}

func (m *MockExec) Command(name string, arg ...string) *exec.Cmd {
	args := m.Called(name, arg)
	return args.Get(0).(*exec.Cmd)
}

// Setup test environment
func setupTest(t *testing.T) (*Service, *MockExec) {
	cfg := &config.Config{
		ConfigPath: "/mock/config/path",
		// Add other necessary fields from your Config struct
	}
	mockExec := new(MockExec)
	service := &Service{
		cfg: cfg,
	}

	return service, mockExec
}

func TestInstall(t *testing.T) {
	s, mockExec := setupTest(t)

	// Mock OS-specific installation
	switch runtime.GOOS {
	case "windows":
		mockExec.On("Command", "sc", mock.Anything).Return(exec.Command("echo", "success"))
		mockExec.On("Command", "reg", mock.Anything).Return(exec.Command("echo", "success"))
	case "darwin":
		mockExec.On("Command", "mkdir", mock.Anything).Return(exec.Command("echo", "success"))
		mockExec.On("Command", "launchctl", mock.Anything).Return(exec.Command("echo", "success"))
	case "linux":
		mockExec.On("Command", "mkdir", mock.Anything).Return(exec.Command("echo", "success"))
		mockExec.On("Command", "systemctl", mock.Anything).Return(exec.Command("echo", "success"))
	}

	err := s.Install()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestUninstall(t *testing.T) {
	s, mockExec := setupTest(t)

	// Mock OS-specific uninstallation
	switch runtime.GOOS {
	case "windows":
		mockExec.On("Command", "sc", mock.Anything).Return(exec.Command("echo", "success"))
	case "darwin":
		mockExec.On("Command", "launchctl", mock.Anything).Return(exec.Command("echo", "success"))
	case "linux":
		mockExec.On("Command", "systemctl", mock.Anything).Return(exec.Command("echo", "success"))
	}

	err := s.Uninstall()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestStart(t *testing.T) {
	s, mockExec := setupTest(t)

	// Mock OS-specific start
	switch runtime.GOOS {
	case "windows":
		mockExec.On("Command", "net", "start", "CrossPlatformAgentService").Return(exec.Command("echo", "success"))
	case "darwin":
		mockExec.On("Command", "launchctl", "load", "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist").Return(exec.Command("echo", "success"))
	case "linux":
		mockExec.On("Command", "systemctl", "start", "crossplatformagent.service").Return(exec.Command("echo", "success"))
	}

	err := s.Start()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestStop(t *testing.T) {
	s, mockExec := setupTest(t)

	// Mock OS-specific stop
	switch runtime.GOOS {
	case "windows":
		mockExec.On("Command", "net", "stop", "CrossPlatformAgentService").Return(exec.Command("echo", "success"))
	case "darwin":
		mockExec.On("Command", "launchctl", "unload", "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist").Return(exec.Command("echo", "success"))
	case "linux":
		mockExec.On("Command", "systemctl", "stop", "crossplatformagent.service").Return(exec.Command("echo", "success"))
	}

	err := s.Stop()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestInstallWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	s, mockExec := setupTest(t)

	mockExec.On("Command", "sc", mock.Anything).Return(exec.Command("echo", "success"))
	mockExec.On("Command", "reg", mock.Anything).Return(exec.Command("echo", "success"))

	err := s.installWindows()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestInstallMacOS(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	s, mockExec := setupTest(t)

	mockExec.On("Command", "launchctl", mock.Anything).Return(exec.Command("echo", "success"))

	err := s.installMacOS()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestInstallLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	s, mockExec := setupTest(t)

	mockExec.On("Command", "systemctl", mock.Anything).Return(exec.Command("echo", "success"))

	err := s.installLinux()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestUninstallWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	s, mockExec := setupTest(t)

	mockExec.On("Command", "sc", "delete", "CrossPlatformAgentService").Return(exec.Command("echo", "success"))

	err := s.uninstallWindows()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestUninstallMacOS(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	s, mockExec := setupTest(t)

	mockExec.On("Command", "launchctl", "unload", "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist").Return(exec.Command("echo", "success"))

	err := s.uninstallMacOS()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}

func TestUninstallLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	s, mockExec := setupTest(t)

	mockExec.On("Command", "systemctl", "disable", "crossplatformagent.service").Return(exec.Command("echo", "success"))
	mockExec.On("Command", "systemctl", "daemon-reload").Return(exec.Command("echo", "success"))

	err := s.uninstallLinux()
	assert.NoError(t, err)

	mockExec.AssertExpectations(t)
}
