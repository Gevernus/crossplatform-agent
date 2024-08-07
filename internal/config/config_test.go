package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Создаем временный конфигурационный файл
	content := []byte(`
api_url: "https://api.example.com"
poll_interval: 60
log_level: "info"
agent_id: "test-agent"
agent_password: "test-password"
`)
	tmpfile, err := ioutil.TempFile("", "config.*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Загружаем конфигурацию
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем значения
	if cfg.APIURL != "https://api.example.com" {
		t.Errorf("Expected APIURL to be 'https://api.example.com', got '%s'", cfg.APIURL)
	}
	if cfg.PollInterval != 60 {
		t.Errorf("Expected PollInterval to be 60, got %d", cfg.PollInterval)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be 'info', got '%s'", cfg.LogLevel)
	}
	if cfg.AgentID != "test-agent" {
		t.Errorf("Expected AgentID to be 'test-agent', got '%s'", cfg.AgentID)
	}
	if cfg.AgentPassword != "test-password" {
		t.Errorf("Expected AgentPassword to be 'test-password', got '%s'", cfg.AgentPassword)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("non_existent_file.yaml")
	if err == nil {
		t.Error("Expected an error when loading non-existent file, got nil")
	}
}