package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"crossplatform-agent/internal/config"
)

// Mock structures
type MockConfig struct {
	mock.Mock
}

type MockService struct {
	mock.Mock
}

type MockTrayManager struct {
	mock.Mock
}

func (m *MockConfig) Load(path string) (*config.Config, error) {
	args := m.Called(path)
	return args.Get(0).(*config.Config), args.Error(1)
}

func (m *MockService) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) IsInstalled() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *MockService) Install() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTrayManager) Run() error {
	args := m.Called()
	return args.Error(0)
}

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel logrus.Level
		expectError   bool
	}{
		{"Debug", "debug", logrus.DebugLevel, false},
		{"Info", "info", logrus.InfoLevel, false},
		{"Warning", "warning", logrus.WarnLevel, false},
		{"Error", "error", logrus.ErrorLevel, false},
		{"Invalid", "invalid", logrus.InfoLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setLogLevel(tt.level)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLevel, logrus.GetLevel())
			}
		})
	}
}
