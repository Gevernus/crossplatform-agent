package tray

import (
	"crossplatform-agent/assets"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/pkg/utils"
	"fmt"
	"path/filepath"

	"github.com/getlantern/systray"
	log "github.com/sirupsen/logrus"
)

type TrayManager struct {
	cfg          *config.Config
	startService func() error
	stopService  func() error
	onExit       func()
}

func NewTrayManager(cfg *config.Config) *TrayManager {
	return &TrayManager{
		cfg:          cfg,
		startService: nil,
		stopService:  nil,
		onExit:       nil,
	}
}

func (tm *TrayManager) SetStartServiceCallback(callback func() error) {
	tm.startService = callback
}

func (tm *TrayManager) SetStopServiceCallback(callback func() error) {
	tm.stopService = callback
}

func (tm *TrayManager) SetOnExitCallback(callback func()) {
	tm.onExit = callback
}

func (tm *TrayManager) Run() error {
	systray.Run(tm.onReady, tm.onExit)
	return nil
}

func (tm *TrayManager) onReady() {
	systray.SetIcon(assets.MustAsset("tray_icon_24x24.png"))
	systray.SetTooltip("Agent Service")
	mStatus := systray.AddMenuItem("Status: Active", "Agent Status")
	mStatus.Disable()

	mAPIURL := systray.AddMenuItem(fmt.Sprintf("API: %s", tm.cfg.APIURL), "API URL")
	mAPIURL.Disable()

	mToggleService := systray.AddMenuItemCheckbox("Turn Off Service", "Toggle service status", true)
	mCommandHistory := systray.AddMenuItem("Executed Commands", "Show executed commands")

	mLogs := systray.AddMenuItem("Open Logs", "Open log directory")
	mQuit := systray.AddMenuItem("Quit", "Quit the agent service")

	go func() {

		for {
			select {
			case <-mLogs.ClickedCh:
				tm.openLogs()
			case <-mToggleService.ClickedCh:
				if mToggleService.Checked() {
					if err := tm.stopService(); err != nil {
						log.Error("Failed to stop service:", err)
					} else {
						mToggleService.SetTitle("Turn On Service")
						mStatus.SetTitle("Status: Inactive")
						mToggleService.Uncheck()
					}
				} else {
					if err := tm.startService(); err != nil {
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
				systray.Quit()
				return
			}
		}
	}()
}

func (tm *TrayManager) showCommandHistory() {
	logFilePath := filepath.Join(tm.cfg.LogDir, "command_history.log")
	log.Infof("Opening log file: %s", logFilePath)
	err := utils.OpenFile(logFilePath)
	if err != nil {
		log.Error("Failed to open log file:", err)
	}
}

type Command struct {
	Command   string
	Timestamp string
}

func (tm *TrayManager) openLogs() {
	log.Info("Opening logs directory")
	err := utils.OpenDirectory(tm.cfg.LogDir)
	if err != nil {
		log.Error("Failed to open logs directory:", err)
	}
}
