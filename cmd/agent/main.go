package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"crossplatform-agent/internal/api"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/internal/service"
	"crossplatform-agent/internal/tray"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Assuming the config file is located in the same directory as the executable
	configPath := filepath.Join(filepath.Dir(execPath), "config.yaml")
	mode := flag.String("mode", "gui", "run mode: service or gui")
	action := flag.String("action", "", "action to perform: install, uninstall, start, stop")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	setLogLevel(cfg.LogLevel)
	apiClient := api.NewClient(cfg.APIURL, cfg.UUID, cfg.DeviceID)
	trayManager := tray.NewTrayManager(cfg)

	// Pass these instances to the Service constructor
	svc := service.New(cfg, apiClient, trayManager)
	trayManager.SetStartServiceCallback(svc.Start)
	trayManager.SetStopServiceCallback(svc.Stop)
	trayManager.SetOnExitCallback(svc.StopService)

	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	switch *action {
	case "install":
		err = svc.Install()
	case "uninstall":
		err = svc.Uninstall()
	case "start":
		err = svc.Start()
	case "stop":
		err = svc.Stop()
	case "":
		switch *mode {
		case "service":
			err = svc.RunAsService()
		case "gui":
			err = svc.RunAsGUI()
		default:
			err = fmt.Errorf("unknown mode: %s", *mode)
		}
	default:
		err = fmt.Errorf("unknown action: %s", *action)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setLogLevel(level string) error {
	parsedLevel, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLevel)
	return nil
}
