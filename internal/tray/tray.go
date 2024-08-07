package tray

import (
	"crossplatform-agent/assets"
	"crossplatform-agent/internal/config"
	"crossplatform-agent/pkg/utils"
	"fmt"

	"github.com/getlantern/systray"
	log "github.com/sirupsen/logrus"
)

type TrayManager struct {
	cfg    *config.Config
	onExit func()
}

// Stop implements service.TrayManager.
func (tm *TrayManager) Stop() {
	panic("unimplemented")
}

func NewTrayManager(cfg *config.Config, onExit func()) *TrayManager {
	return &TrayManager{
		cfg:    cfg,
		onExit: onExit,
	}
}

func (tm *TrayManager) Run() {
	systray.Run(tm.onReady, tm.onExit)
}

func (tm *TrayManager) onReady() {
	systray.SetIcon(assets.MustAsset("tray_icon_24x24.png"))
	systray.SetTitle("Agent Service")
	systray.SetTooltip("Agent Service")

	mStatus := systray.AddMenuItem("Status: Active", "Agent Status")
	mStatus.Disable()

	mAPIURL := systray.AddMenuItem(fmt.Sprintf("API: %s", tm.cfg.APIURL), "API URL")
	mAPIURL.Disable()

	mLogs := systray.AddMenuItem("Open Logs", "Open log directory")

	mQuit := systray.AddMenuItem("Quit", "Quit the agent service")

	go func() {
		for {
			select {
			case <-mLogs.ClickedCh:
				tm.openLogs()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (tm *TrayManager) openLogs() {
	log.Info("Opening logs directory")
	err := utils.OpenDirectory(tm.cfg.LogPath)
	if err != nil {
		log.Error("Failed to open logs directory:", err)
	}
}
