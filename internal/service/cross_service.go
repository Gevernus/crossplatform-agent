package service

import (
	"crossplatform-agent/internal/api"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/pkg/commands"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"

	log "github.com/sirupsen/logrus"
)

var logFile *os.File

type APIClient interface {
	SendStatus() error
	GetCommands() ([]api.Command, error)
}

type TrayManager interface {
	Run() error
}

type CrossService struct {
	cfg *config.Config
	api APIClient
	Svc ServiceWrapper
}

func New(cfg *config.Config, apiClient APIClient) (*CrossService, error) {
	svc := &CrossService{
		cfg: cfg,
		api: apiClient,
	}

	// Configure the service.Config for kardianos/service
	svcConfig := &service.Config{
		Name:        "CrossPlatformAgentService",
		DisplayName: "Cross Platform Agent Service",
		Description: "This service handles cross-platform agent operations.",
		Option: service.KeyValue{
			"RunAtLoad": true,
			"KeepAlive": true,
		},
		Arguments: []string{
			"-mode",
			"service",
		},
	}

	// Create the kardianos/service.Service object
	serviceObj, err := service.New(svc, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	svc.Svc = serviceObj

	return svc, nil
}

func (s *CrossService) Start(svc service.Service) error {
	// Start should not block. Do the actual work in a goroutine.
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to determine executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)

	logFilePath := filepath.Join(execDir, "CrossPlatformAgentService.log")

	// Open log files
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)

	go s.runMainLoop()
	return nil
}

// Stop stops the service
func (s *CrossService) Stop(svc service.Service) error {
	log.Println("Stopping service...")
	if logFile != nil {
		logFile.Close()
	}
	return nil
}

func (s *CrossService) StartService() error {
	err := service.Control(s.Svc, "start")
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}
	return nil
}

func (s *CrossService) StopService() error {
	err := service.Control(s.Svc, "stop")
	if err != nil {
		return fmt.Errorf("failed to stop service: %v", err)
	}
	return nil
}

func (s *CrossService) Install() error {
	err := s.Svc.Install()
	if err != nil {
		return fmt.Errorf("failed to install service: %v", err)
	}
	log.Info("Service installation completed successfully")
	return nil
}

func (s *CrossService) IsInstalled() (bool, error) {
	status, err := s.Svc.Status()
	if err != nil {
		if err == service.ErrNotInstalled {
			// The service is not installed
			return false, nil
		}
		// An unexpected error occurred
		return false, err
	}

	// If no error and we got a status, the service is installed
	return status != service.StatusUnknown, nil
}

func (s *CrossService) IsActive() (bool, error) {
	status, err := s.Svc.Status()
	if err != nil {
		if err == service.ErrNotInstalled {
			return false, fmt.Errorf("service is not installed")
		}
		return false, err
	}

	// Check if the service status is running
	return status == service.StatusRunning, nil
}

func (s *CrossService) Run() error {
	if err := s.Svc.Run(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *CrossService) runMainLoop() error {
	log.Println("Service started")
	s.performTasks()
	ticker := time.NewTicker(time.Duration(s.cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.performTasks()
	}

	return nil
}

func (s *CrossService) performTasks() {
	if err := s.sendStatus(); err != nil {
		log.Error("Failed to send status:", err)
	}

	if err := s.processCommands(); err != nil {
		log.Error("Failed to process commands:", err)
	}
}

func (s *CrossService) sendStatus() error {
	log.Debug("Sending status update")
	return s.api.SendStatus()
}

func (s *CrossService) processCommands() error {
	log.Debug("Getting commands from server")
	commands, err := s.api.GetCommands()
	if err != nil {
		return err
	}
	log.Debugf("Got: %d commands", len(commands))
	for _, cmd := range commands {
		log.Debug("Processing command:", cmd.Command)
		s.logCommand(cmd.String())
		if err := s.executeCommand(cmd); err != nil {
			log.Error("Failed to execute command:", cmd.Command, err)
		}
	}

	return nil
}

func (s *CrossService) executeCommand(cmd api.Command) error {
	switch cmd.Command {
	case "shutdown":
		log.Info("Executing shutdown command")
		if err := commands.Shutdown(); err != nil {
			log.Error("Failed to execute shutdown command:", err)
			return err
		}
	case "shutdownNetwork":
		log.Info("Executing network shutdown command")
		if err := commands.ShutdownNetwork(); err != nil {
			log.Error("Failed to execute network shutdown command:", err)
			return err
		}
	default:
		log.Warn("Unknown command received:", cmd.Command)
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
	return nil
}

func (s *CrossService) logCommand(command string) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to determine executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)

	logFilePath := filepath.Join(execDir, "Command_history.log")
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Failed to open log file:", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("%s: %s\n", timestamp, command)); err != nil {
		log.Error("Failed to write to log file:", err)
	}
}
