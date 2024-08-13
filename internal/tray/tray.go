package tray

import (
	"crossplatform-agent/assets"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/pkg/utils"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type ServiceInterface interface {
	IsActive() (bool, error)
	StartService() error
	StopService() error
}

type UtilsInterface interface {
	OpenDirectory(path string) error
	OpenFile(filePath string) error
}

type TrayManager struct {
	cfg     *config.Config
	svc     ServiceInterface
	systray SystrayWrapper
	utils   UtilsInterface
}

func NewTrayManager(cfg *config.Config, svc ServiceInterface) *TrayManager {
	return &TrayManager{
		cfg:     cfg,
		svc:     svc,
		systray: &DefaultSystrayWrapper{},
		utils:   &utils.Impl{},
	}
}

func (tm *TrayManager) Run() error {
	tm.systray.Run(tm.onReady, tm.OnExit)
	return nil
}

func (tm *TrayManager) OnExit() {
	// tm.StopService()
}

func (tm *TrayManager) onReady() {
	tm.systray.SetTooltip("Agent Service")
	err := setTrayIcon(tm)
	if err != nil {
		log.Errorf("Failed to set tray icon: %v", err)
	}

	active, err := tm.svc.IsActive()
	if err != nil {
		log.Fatalf("Failed to check if service is active: %v", err)
	}

	mStatus := tm.systray.AddMenuItem("Status: Active", "Agent Status")
	mStatus.Disable()

	mAPIURL := tm.systray.AddMenuItem(fmt.Sprintf("API: %s", tm.cfg.APIURL), "API URL")
	mAPIURL.Disable()

	mToggleService := tm.systray.AddMenuItemCheckbox("Turn Off Service", "Toggle service status", true)
	if active {
		mToggleService.Check()
	} else {
		mToggleService.SetTitle("Turn On Service")
		mStatus.SetTitle("Status: Inactive")
		mToggleService.Uncheck()
	}
	mCommandHistory := tm.systray.AddMenuItem("Executed Commands", "Show executed commands")

	mLogs := tm.systray.AddMenuItem("Open Logs", "Open log directory")
	mQuit := tm.systray.AddMenuItem("Quit", "Quit the agent service")

	go func() {
		for {
			select {
			case <-mLogs.ClickedCh:
				tm.openLogs()
			case <-mToggleService.ClickedCh:
				if mToggleService.Checked() {
					if err := tm.svc.StopService(); err != nil {
						log.Error("Failed to stop service:", err)
					} else {
						mToggleService.SetTitle("Turn On Service")
						mStatus.SetTitle("Status: Inactive")
						mToggleService.Uncheck()
					}
				} else {
					if err := tm.svc.StartService(); err != nil {
						log.Error("Failed to start service:", err)
					} else {
						mToggleService.SetTitle("Turn Off Service")
						mStatus.SetTitle("Status: Active")
						mToggleService.Check()
					}
				}
			case <-mCommandHistory.ClickedCh:
				tm.showCommandHistory()
			case <-mQuit.ClickedCh:
				tm.systray.Quit()
				return
			}
		}
	}()
}

func setTrayIcon(tm *TrayManager) error {
	var iconBytes []byte
	var err error

	switch runtime.GOOS {
	case "windows":
		iconBytes, err = assets.Asset("assets/tray_icon.ico")
	default:
		iconBytes, err = assets.Asset("assets/tray_icon_24x24.png")
	}

	if err != nil {
		return fmt.Errorf("failed to load tray icon: %w", err)
	}

	tm.systray.SetIcon(iconBytes)
	return nil
}

func (tm *TrayManager) showCommandHistory() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to determine executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)

	logFilePath := filepath.Join(execDir, "Command_history.log")
	log.Infof("Opening log file: %s", logFilePath)
	err = tm.utils.OpenFile(logFilePath)
	if err != nil {
		log.Error("Failed to open log file:", err)
	}
}

type Command struct {
	Command   string
	Timestamp string
}

func (tm *TrayManager) openLogs() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to determine executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)
	log.Info("Opening logs directory: ", execDir)
	err = tm.utils.OpenDirectory(execDir)
	if err != nil {
		log.Error("Failed to open logs directory:", err)
	}
}
