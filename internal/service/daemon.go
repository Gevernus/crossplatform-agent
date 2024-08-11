package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func (s *Service) Install() error {
	log.Infof("Starting service installation on OS: %s", runtime.GOOS)

	// Get absolute path for config
	absConfigPath, err := filepath.Abs(s.cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for config: %v", err)
	}
	s.cfg.ConfigPathAbs = absConfigPath

	var installErr error
	switch runtime.GOOS {
	case "windows":
		installErr = s.installWindows()
	case "darwin":
		installErr = s.installMacOS()
	case "linux":
		installErr = s.installLinux()
	default:
		installErr = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	if installErr != nil {
		log.Error("Service installation failed:", installErr)
		return installErr
	}
	log.Info("Service installation completed successfully")
	return nil
}

func (s *Service) Uninstall() error {
	log.Info("Starting service uninstallation")
	var err error
	switch runtime.GOOS {
	case "windows":
		err = s.uninstallWindows()
	case "darwin":
		err = s.uninstallMacOS()
	case "linux":
		err = s.uninstallLinux()
	default:
		err = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	if err != nil {
		log.Error("Service uninstallation failed:", err)
		return err
	}
	log.Info("Service uninstallation completed successfully")
	return nil
}

func (s *Service) Start() error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("net", "start", "CrossPlatformAgentService").Run()
	case "darwin":
		return exec.Command("launchctl", "load", "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist").Run()
	case "linux":
		return exec.Command("systemctl", "start", "crossplatformagent.service").Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (s *Service) Stop() error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("net", "stop", "CrossPlatformAgentService").Run()
	case "darwin":
		return exec.Command("launchctl", "unload", "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist").Run()
	case "linux":
		return exec.Command("systemctl", "stop", "crossplatformagent.service").Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (s *Service) installWindows() error {
	logDir := filepath.Join(os.Getenv("ProgramData"), "CrossPlatformAgent", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "crossplatformagent.log")
	errFile := filepath.Join(logDir, "crossplatformagent.err")
	if err := os.WriteFile(logFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	if err := os.WriteFile(errFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create error log file: %v", err)
	}

	cmd := exec.Command("sc", "create", "CrossPlatformAgentService",
		"binPath=", fmt.Sprintf("%s -config %s", os.Args[0], s.cfg.ConfigPath),
		"start=", "demand",
		"DisplayName=", "Cross Platform Agent Service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Windows service: %v", err)
	}

	regCmd := exec.Command("reg", "add", `HKLM\SYSTEM\CurrentControlSet\Services\CrossPlatformAgentService\Parameters`,
		"/v", "AppStdout", "/t", "REG_EXPAND_SZ", "/d", logFile, "/f")
	if err := regCmd.Run(); err != nil {
		return fmt.Errorf("failed to set up stdout logging: %v", err)
	}

	regCmd = exec.Command("reg", "add", `HKLM\SYSTEM\CurrentControlSet\Services\CrossPlatformAgentService\Parameters`,
		"/v", "AppStderr", "/t", "REG_EXPAND_SZ", "/d", errFile, "/f")
	if err := regCmd.Run(); err != nil {
		return fmt.Errorf("failed to set up stderr logging: %v", err)
	}

	return nil
}

func (s *Service) uninstallWindows() error {
	return exec.Command("sc", "delete", "CrossPlatformAgentService").Run()
}

func (s *Service) installMacOS() error {
	logDir := "/var/log/crossplatformagent"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "crossplatformagent.log")
	errFile := filepath.Join(logDir, "crossplatformagent.err")
	if err := os.WriteFile(logFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	if err := os.WriteFile(errFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create error log file: %v", err)
	}

	plistPath := "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist"

	tmpl, err := template.New("plist").Parse(macOSPlistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse plist template: %v", err)
	}

	file, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %v", err)
	}
	defer file.Close()

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	err = tmpl.Execute(file, struct {
		ExecPath   string
		ConfigPath string
		LogFile    string
		ErrFile    string
	}{
		ExecPath:   execPath,
		ConfigPath: s.cfg.ConfigPathAbs,
		LogFile:    logFile,
		ErrFile:    errFile,
	})
	if err != nil {
		return fmt.Errorf("failed to write plist file: %v", err)
	}

	// cmd := exec.Command("launchctl", "load", plistPath)
	// return cmd.Run()
	return nil
}

func (s *Service) uninstallMacOS() error {
	plistPath := "/Library/LaunchDaemons/com.gevernus.crossplatformagent.plist"

	cmd := exec.Command("launchctl", "unload", plistPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unload service: %v", err)
	}

	if err := os.Remove(plistPath); err != nil {
		return fmt.Errorf("failed to remove plist file: %v", err)
	}

	return nil
}

const macOSPlistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.gevernus.crossplatformagent</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExecPath}}</string>
        <string>-mode</string>
        <string>service</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.LogFile}}</string>
    <key>StandardErrorPath</key>
    <string>{{.ErrFile}}</string>
</dict>
</plist>`

func (s *Service) installLinux() error {
	logDir := "/var/log/crossplatformagent"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "crossplatformagent.log")
	errFile := filepath.Join(logDir, "crossplatformagent.err")
	if err := os.WriteFile(logFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	if err := os.WriteFile(errFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create error log file: %v", err)
	}

	serviceFilePath := "/etc/systemd/system/crossplatformagent.service"

	tmpl, err := template.New("service").Parse(linuxServiceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %v", err)
	}

	file, err := os.Create(serviceFilePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}
	defer file.Close()

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	err = tmpl.Execute(file, struct {
		ExecPath   string
		ConfigPath string
		LogFile    string
		ErrFile    string
	}{
		ExecPath:   execPath,
		ConfigPath: s.cfg.ConfigPath,
		LogFile:    logFile,
		ErrFile:    errFile,
	})
	if err != nil {
		return fmt.Errorf("failed to write service file: %v", err)
	}

	cmd := exec.Command("systemctl", "daemon-reload")
	return cmd.Run()
}

func (s *Service) uninstallLinux() error {
	serviceFilePath := "/etc/systemd/system/crossplatformagent.service"

	cmd := exec.Command("systemctl", "disable", "crossplatformagent.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable service: %v", err)
	}

	if err := os.Remove(serviceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %v", err)
	}

	cmd = exec.Command("systemctl", "daemon-reload")
	return cmd.Run()
}

const linuxServiceTemplate = `[Unit]
Description=Cross Platform Agent Service
After=network.target

[Service]
ExecStart={{.ExecPath}} -config {{.ConfigPath}}
Restart=always
User=root
StandardOutput=append:{{.LogFile}}
StandardError=append:{{.ErrFile}}

[Install]
WantedBy=multi-user.target`
