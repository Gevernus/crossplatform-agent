package service

import (
	"crossplatform-agent/internal/api"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/pkg/commands"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type APIClient interface {
	SendStatus(status string) error
	GetCommands() ([]api.Command, error)
}

type TrayManager interface {
	Run() error
}

type Service struct {
	cfg  *config.Config
	api  APIClient
	tray TrayManager
	wg   sync.WaitGroup
}

func New(cfg *config.Config, apiClient APIClient, trayManager TrayManager) *Service {
	return &Service{
		cfg:  cfg,
		api:  apiClient,
		tray: trayManager,
	}
}

func (s *Service) RunAsService() error {
	log.Info("Starting in service mode")
	return s.runMainLoop()
}

func (s *Service) RunAsGUI() error {
	log.Info("Starting in GUI mode")
	s.Start()
	return s.runTray()
}

func (s *Service) runMainLoop() error {
	ticker := time.NewTicker(time.Duration(s.cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.sendStatus(); err != nil {
				log.Error("Failed to send status:", err)
			}

			if err := s.processCommands(); err != nil {
				log.Error("Failed to process commands:", err)
			}
		}
	}
}

func (s *Service) runTray() error {
	defer s.wg.Done()
	return s.tray.Run()
}

func (s *Service) StopService() {
	log.Info("Stopping service")
	// Implement any cleanup or shutdown logic here
	s.wg.Done()
}

func (s *Service) sendStatus() error {
	log.Debug("Sending status update")
	return s.api.SendStatus("active")
}

func (s *Service) processCommands() error {
	commands, err := s.api.GetCommands()
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		log.Debug("Processing command:", cmd.Action)
		s.logCommand(cmd.String())
		if err := s.executeCommand(cmd); err != nil {
			log.Error("Failed to execute command:", cmd.Action, err)
		}
	}

	return nil
}

func (s *Service) executeCommand(cmd api.Command) error {
	switch cmd.Action {
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
		log.Warn("Unknown command received:", cmd.Action)
		return fmt.Errorf("unknown command: %s", cmd.Action)
	}
	return nil
}

func (s *Service) logCommand(command string) {
	logFilePath := filepath.Join(s.cfg.LogDir, "command_history.log")
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
