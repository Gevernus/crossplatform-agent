package service

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func (s *Service) Install() error {
	log.Info("Starting service installation")
	var err error
	switch runtime.GOOS {
	case "windows":
		err = s.installWindows()
	case "darwin":
		err = s.installMacOS()
	case "linux":
		err = s.installLinux()
	default:
		err = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	if err != nil {
		log.Error("Service installation failed:", err)
		return err
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
		return exec.Command("launchctl", "load", "/Library/LaunchDaemons/com.yourcompany.crossplatformagent.plist").Run()
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
		return exec.Command("launchctl", "unload", "/Library/LaunchDaemons/com.yourcompany.crossplatformagent.plist").Run()
	case "linux":
		return exec.Command("systemctl", "stop", "crossplatformagent.service").Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (s *Service) installWindows() error {
	cmd := exec.Command("sc", "create", "CrossPlatformAgentService",
		"binPath=", os.Args[0],
		"start=", "auto",
		"DisplayName=", "Cross Platform Agent Service")
	return cmd.Run()
}

func (s *Service) uninstallWindows() error {
	return exec.Command("sc", "delete", "CrossPlatformAgentService").Run()
}

func (s *Service) installMacOS() error {
	plistPath := "/Library/LaunchDaemons/com.yourcompany.crossplatformagent.plist"

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
	}{
		ExecPath:   execPath,
		ConfigPath: s.cfg.ConfigPath,
	})
	if err != nil {
		return fmt.Errorf("failed to write plist file: %v", err)
	}

	cmd := exec.Command("launchctl", "load", plistPath)
	return cmd.Run()
}

func (s *Service) uninstallMacOS() error {
	plistPath := "/Library/LaunchDaemons/com.yourcompany.crossplatformagent.plist"

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
    <string>com.yourcompany.crossplatformagent</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ExecPath}}</string>
        <string>-config</string>
        <string>{{.ConfigPath}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/crossplatformagent.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/crossplatformagent.err</string>
</dict>
</plist>`

func (s *Service) installLinux() error {
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
	}{
		ExecPath:   execPath,
		ConfigPath: s.cfg.ConfigPath,
	})
	if err != nil {
		return fmt.Errorf("failed to write service file: %v", err)
	}

	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	cmd = exec.Command("systemctl", "enable", "crossplatformagent.service")
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

[Install]
WantedBy=multi-user.target`
