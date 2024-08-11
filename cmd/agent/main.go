package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	logrus "github.com/sirupsen/logrus"

	"crossplatform-agent/internal/api"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/internal/service"
	"crossplatform-agent/internal/tray"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Redirect standard log output to logrus
	log.SetOutput(logrus.StandardLogger().Writer())
	log.SetFlags(0)

	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Determine the directory where the executable is located
	execDir := filepath.Dir(execPath)

	// Assuming the config file is located in the same directory as the executable
	configPath := filepath.Join(execDir, "config.yaml")
	mode := flag.String("mode", "gui", "run mode: service or gui")

	logrus.Debug("Arguments passed:", os.Args)
	flag.Parse()
	logrus.Debug("Mode after parsing:", *mode)

	cfg, err := config.Load(configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	setLogLevel(cfg.LogLevel)
	apiClient := api.NewClient(cfg.APIURL, cfg.UUID, cfg.DeviceID)

	// Pass these instances to the Service constructor
	svc, err := service.New(cfg, apiClient)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	switch *mode {
	case "service":
		err = svc.Run()
	default:
		// initLogFiles("tray", execDir)
		installed, errr := svc.IsInstalled()
		if errr != nil {
			logrus.Fatalf("Failed to check if service is installed: %v", err)
		}

		if !installed {
			logrus.Println("Service not installed, installing now...")
			err := svc.Install()
			if err != nil {
				logrus.Fatalf("Failed to install service: %v", err)
			}
			logrus.Println("Service installed successfully")
		}

		trayManager := tray.NewTrayManager(cfg, svc)
		err = trayManager.Run()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setLogLevel(level string) error {
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(parsedLevel)
	return nil
}
