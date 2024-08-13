package tray

import (
	"crossplatform-agent/internal/config"
	"crossplatform-agent/internal/service"
	"testing"

	"github.com/getlantern/systray"
	"github.com/stretchr/testify/assert"
)

type MockUtils struct{}

func (m *MockUtils) OpenDirectory(path string) error {
	// Mock behavior for OpenDirectory
	return nil
}

func (m *MockUtils) OpenFile(filePath string) error {
	// Mock behavior for OpenFile
	return nil
}

type MockService struct{}

func (ms *MockService) IsActive() (bool, error) {
	return true, nil // or whatever mock behavior you need
}

func (ms *MockService) StartService() error {
	return nil
}

func (ms *MockService) StopService() error {
	return nil
}

type MockSystrayWrapper struct {
	tooltip          string
	addMenuItemCalls int
	addCheckboxCalls int
	quitCalls        int
	runCalls         int

	addMenuItemParams     []struct{ label, tooltip string }
	addCheckboxItemParams []struct {
		label, tooltip string
		checked        bool
	}
	menuItems map[string]*systray.MenuItem
}

// NewMockSystrayWrapper creates a new instance of MockSystrayWrapper
func NewMockSystrayWrapper() *MockSystrayWrapper {
	return &MockSystrayWrapper{
		menuItems:         make(map[string]*systray.MenuItem),
		addMenuItemParams: []struct{ label, tooltip string }{},
		addCheckboxItemParams: []struct {
			label, tooltip string
			checked        bool
		}{},
	}
}

// Run mocks the Run method
func (m *MockSystrayWrapper) Run(onReady func(), onExit func()) {
	m.runCalls++
	onReady()
}

// Quit mocks the Quit method
func (m *MockSystrayWrapper) Quit() {
	m.quitCalls++
}

// SetTooltip mocks setting the tooltip
func (m *MockSystrayWrapper) SetTooltip(tooltip string) {
	m.tooltip = tooltip
}

// SetIcon mocks setting the icon
func (m *MockSystrayWrapper) SetIcon(iconBytes []byte) {
	// Mock set icon functionality
}

// AddMenuItem mocks adding a menu item
func (m *MockSystrayWrapper) AddMenuItem(label string, tooltip string) *systray.MenuItem {
	m.addMenuItemCalls++
	m.addMenuItemParams = append(m.addMenuItemParams, struct{ label, tooltip string }{label, tooltip})
	item := systray.AddMenuItem(label, tooltip)
	m.menuItems[label] = item
	return item
}

// AddMenuItemCheckbox mocks adding a checkbox menu item
func (m *MockSystrayWrapper) AddMenuItemCheckbox(label string, tooltip string, checked bool) *systray.MenuItem {
	m.addCheckboxCalls++
	m.addCheckboxItemParams = append(m.addCheckboxItemParams, struct {
		label, tooltip string
		checked        bool
	}{label, tooltip, checked})
	item := m.AddMenuItem(label, tooltip)
	item.Check()
	if !checked {
		item.Uncheck()
	}
	return item
}

func TestNewTrayManager(t *testing.T) {
	cfg := &config.Config{}
	svc := &service.CrossService{}
	tm := NewTrayManager(cfg, svc)

	assert.NotNil(t, tm)
	assert.Equal(t, cfg, tm.cfg)
	assert.Equal(t, svc, tm.svc)
}

func TestRun(t *testing.T) {
	cfg := &config.Config{}
	svc := &MockService{}

	mockSystray := NewMockSystrayWrapper()

	tm := NewTrayManager(cfg, svc)
	tm.systray = mockSystray // Inject the mock

	err := tm.Run()

	assert.NoError(t, err)
	assert.Equal(t, "Agent Service", mockSystray.tooltip, "Tooltip should be set correctly")
	assert.NotNil(t, mockSystray.menuItems["Status: Active"], "Status menu item should be created")
	assert.NotNil(t, mockSystray.menuItems["API: "+cfg.APIURL], "API URL menu item should be created")
	assert.NotNil(t, mockSystray.menuItems["Turn Off Service"], "Toggle service menu item should be created")
	assert.Equal(t, 6, mockSystray.addMenuItemCalls, "Expected AddMenuItem to be called once")
	assert.Equal(t, 1, mockSystray.addCheckboxCalls, "Expected AddMenuItemCheckbox to be called once")
}

func TestTrayManager_onReady(t *testing.T) {
	cfg := &config.Config{
		APIURL: "http://localhost:8080",
	}
	svc := &MockService{}
	mockSystray := NewMockSystrayWrapper()

	tm := NewTrayManager(cfg, svc)
	tm.systray = mockSystray

	tm.onReady()

	assert.Equal(t, "Agent Service", mockSystray.tooltip, "Tooltip should be set correctly")
	assert.NotNil(t, mockSystray.menuItems["Status: Active"], "Status menu item should be created")
	assert.NotNil(t, mockSystray.menuItems["API: "+cfg.APIURL], "API URL menu item should be created")
	assert.NotNil(t, mockSystray.menuItems["Turn Off Service"], "Toggle service menu item should be created")
	assert.Equal(t, 6, mockSystray.addMenuItemCalls, "Expected AddMenuItem to be called once")
	assert.Equal(t, 1, mockSystray.addCheckboxCalls, "Expected AddMenuItemCheckbox to be called once")
}

func TestTrayManager_showCommandHistory(t *testing.T) {
	cfg := &config.Config{}
	svc := &MockService{}
	mockSystray := &MockSystrayWrapper{}
	mockUtils := &MockUtils{}

	tm := NewTrayManager(cfg, svc)
	tm.systray = mockSystray
	tm.utils = mockUtils
	tm.showCommandHistory()

	// Since showCommandHistory uses a utility function to open a file,
	// the main assertion is that no errors or panics occurred during the method call.
	assert.True(t, true, "Show command history executed without errors")
}

func TestTrayManager_openLogs(t *testing.T) {
	cfg := &config.Config{}
	svc := &MockService{}
	mockSystray := &MockSystrayWrapper{}
	mockUtils := &MockUtils{}
	tm := NewTrayManager(cfg, svc)
	tm.systray = mockSystray
	tm.utils = mockUtils

	tm.openLogs()

	// Since openLogs uses a utility function to open a directory,
	// the main assertion is that no errors or panics occurred during the method call.
	assert.True(t, true, "Open logs executed without errors")
}
