package service

import (
	"crossplatform-agent/internal/api"
	"crossplatform-agent/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPIClient is a mock implementation of the api.Client
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

// MockTrayManager is a mock implementation of the tray.TrayManager
type MockTrayManager struct {
	mock.Mock
}

func (m *MockTrayManager) Run() error {
	m.Called()
	return nil
}

func (m *MockTrayManager) Stop() error {
	m.Called()
	return nil
}

func TestService_Run(t *testing.T) {
	cfg := &config.Config{
		APIURL:        "http://localhost",
		PollInterval:  1,
		LogLevel:      "debug",
		AgentID:       "test-agent",
		AgentPassword: "test-pass",
	}

	mockAPI := new(MockAPIClient)
	mockTray := new(MockTrayManager)

	mockAPI.On("SendStatus", "active").Return(nil) // Adjusted to match expected calls
	mockAPI.On("GetCommands").Return([]api.Command{}, nil)
	mockTray.On("Run").Return(nil)
	// mockTray.On("Stop").Return(nil)

	service := New(cfg, mockAPI, mockTray)

	go func() {
		time.Sleep(2 * time.Second)
		service.StopService()
	}()

	// err := service.Run()
	// assert.NoError(t, err)

	mockAPI.AssertExpectations(t)
	mockTray.AssertExpectations(t)
}

func TestService_StopService(t *testing.T) {
	mockTray := new(MockTrayManager)
	service := &Service{
		tray: mockTray,
	}
	service.wg.Add(1)

	// mockTray.On("Stop").Return(nil)

	go func() {
		time.Sleep(1 * time.Second)
		service.StopService()
	}()

	service.wg.Wait()
	mockTray.AssertExpectations(t)
}

func TestService_sendStatus(t *testing.T) {
	cfg := &config.Config{
		APIURL:        "http://localhost",
		PollInterval:  1,
		LogLevel:      "debug",
		AgentID:       "test-agent",
		AgentPassword: "test-pass",
	}

	mockAPI := new(MockAPIClient)
	mockAPI.On("SendStatus", "active").Return(nil)

	service := New(cfg, mockAPI, nil)

	err := service.sendStatus()
	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
}

func TestService_processCommands(t *testing.T) {
	cfg := &config.Config{
		APIURL:        "http://localhost",
		PollInterval:  1,
		LogLevel:      "debug",
		AgentID:       "test-agent",
		AgentPassword: "test-pass",
	}

	mockAPI := new(MockAPIClient)
	commands := []api.Command{
		{Command: "shutdown"},
	}

	mockAPI.On("GetCommands").Return(commands, nil)
	// mockAPI.On("SendStatus", "active").Return(nil)

	service := New(cfg, mockAPI, nil)

	err := service.processCommands()
	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
}

// func TestService_executeCommand(t *testing.T) {
// 	service := &Service{}
// 	cmd := api.Command{Action: "shutdown"}

// 	err := service.executeCommand(cmd)
// 	assert.NoError(t, err)
// }
