package service_test

import (
	"testing"

	"crossplatform-agent/internal/api"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/internal/service"

	svc "github.com/kardianos/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPIClient is a mock implementation of the APIClient interface
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) SendStatus() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAPIClient) GetCommands() ([]api.Command, error) {
	args := m.Called()
	return args.Get(0).([]api.Command), args.Error(1)
}

type MockServiceWrapper struct {
	StartCalled        bool
	StopCalled         bool
	InstallCalled      bool
	UninstallCalled    bool
	StatusCalled       bool
	RunCalled          bool
	LoggerCalled       bool
	PlatformCalled     bool
	RestartCalled      bool
	StringCalled       bool
	SystemLoggerCalled bool

	StartError        error
	StopError         error
	InstallError      error
	UninstallError    error
	StatusError       error
	StatusResult      svc.Status
	RunError          error
	LoggerError       error
	LoggerResult      svc.Logger
	PlatformResult    string
	RestartError      error
	StringResult      string
	SystemLoggerError error
}

func (m *MockServiceWrapper) Start() error {
	m.StartCalled = true
	return m.StartError
}

func (m *MockServiceWrapper) Stop() error {
	m.StopCalled = true
	return m.StopError
}

func (m *MockServiceWrapper) Install() error {
	m.InstallCalled = true
	return m.InstallError
}

func (m *MockServiceWrapper) Uninstall() error {
	m.UninstallCalled = true
	return m.UninstallError
}

func (m *MockServiceWrapper) Status() (svc.Status, error) {
	m.StatusCalled = true
	return m.StatusResult, m.StatusError
}

func (m *MockServiceWrapper) Run() error {
	m.RunCalled = true
	return m.RunError
}

func (m *MockServiceWrapper) Logger(errs chan<- error) (svc.Logger, error) {
	m.LoggerCalled = true
	return m.LoggerResult, m.LoggerError
}

func (m *MockServiceWrapper) Platform() string {
	m.PlatformCalled = true
	return m.PlatformResult
}

func (m *MockServiceWrapper) Restart() error {
	m.RestartCalled = true
	return m.RestartError
}

func (m *MockServiceWrapper) String() string {
	m.StringCalled = true
	return m.StringResult
}

func (m *MockServiceWrapper) SystemLogger(errs chan<- error) (svc.Logger, error) {
	m.SystemLoggerCalled = true
	return m.LoggerResult, m.SystemLoggerError
}

// TestNew tests the New function of CrossService
func TestNew(t *testing.T) {
	cfg := &config.Config{}
	apiClient := &MockAPIClient{}
	svc, err := service.New(cfg, apiClient)

	assert.NoError(t, err)
	assert.NotNil(t, svc)
}

// func TestCrossService_Start(t *testing.T) {
// 	mockSvc := &MockServiceWrapper{}
// 	cs := &service.CrossService{
// 		Svc: mockSvc,
// 	}

// 	err := cs.Start(mockSvc)
// 	assert.NoError(t, err)
// 	assert.False(t, mockSvc.StartCalled, "Start shouldn't have been called")
// }

// TestCrossService_Stop tests the Stop function of CrossService.
func TestCrossService_Stop(t *testing.T) {
	mockSvc := &MockServiceWrapper{}
	cs := &service.CrossService{
		Svc: mockSvc,
	}

	err := cs.Stop(mockSvc)
	assert.NoError(t, err)
	assert.False(t, mockSvc.StopCalled, "Stop shouldn't have been called")
}

// TestCrossService_IsInstalled tests the IsInstalled function of CrossService.
func TestCrossService_IsInstalled(t *testing.T) {
	mockSvc := &MockServiceWrapper{}
	cs := &service.CrossService{
		Svc: mockSvc,
	}

	mockSvc.StatusResult = svc.StatusRunning
	installed, err := cs.IsInstalled()
	assert.NoError(t, err)
	assert.True(t, installed, "Service should be installed")

	mockSvc.StatusResult = svc.StatusUnknown
	installed, err = cs.IsInstalled()
	assert.NoError(t, err)
	assert.False(t, installed, "Service should not be installed")
}

// TestCrossService_IsActive tests the IsActive function of CrossService.
func TestCrossService_IsActive(t *testing.T) {
	mockSvc := &MockServiceWrapper{}
	cs := &service.CrossService{
		Svc: mockSvc,
	}

	mockSvc.StatusResult = svc.StatusRunning
	active, err := cs.IsActive()
	assert.NoError(t, err)
	assert.True(t, active, "Service should be active")

	mockSvc.StatusResult = svc.StatusStopped
	active, err = cs.IsActive()
	assert.NoError(t, err)
	assert.False(t, active, "Service should not be active")
}

// TestCrossService_Run tests the Run function of CrossService.
func TestCrossService_Run(t *testing.T) {
	mockSvc := &MockServiceWrapper{}
	cs := &service.CrossService{
		Svc: mockSvc,
	}

	err := cs.Run()
	assert.NoError(t, err)
	assert.True(t, mockSvc.RunCalled, "Run should have been called")
}

// TestCrossService_SendStatus tests the sendStatus function of CrossService.
// func TestCrossService_SendStatus(t *testing.T) {
// 	mockAPI := &MockAPIClient{}
// 	cs := &service.CrossService{
// 		API: mockAPI,
// 	}

// 	err := cs.SendStatus()
// 	assert.NoError(t, err)
// 	assert.True(t, mockAPI.SendStatusCalled, "SendStatus should have been called")
// }

// // TestCrossService_PerformTasks tests the performTasks function of CrossService.
// func TestCrossService_PerformTasks(t *testing.T) {
// 	mockAPI := &MockAPIClient{}
// 	cs := &service.CrossService{
// 		API: mockAPI,
// 	}

// 	cs.PerformTasks()
// 	assert.True(t, mockAPI.SendStatusCalled, "SendStatus should have been called")
// 	assert.True(t, mockAPI.GetCommandsCalled, "GetCommands should have been called")
// }

// // TestCrossService_ExecuteCommand tests the executeCommand function of CrossService.
// func TestCrossService_ExecuteCommand(t *testing.T) {
// 	cs := &service.CrossService{}

// 	cmd := api.Command{Command: "shutdown"}
// 	err := cs.ExecuteCommand(cmd)
// 	assert.NoError(t, err)

// 	cmd = api.Command{Command: "unknown"}
// 	err = cs.ExecuteCommand(cmd)
// 	assert.Error(t, err)
// }
