package tray

import "github.com/getlantern/systray"

type SystrayWrapper interface {
	Run(onReady func(), onExit func())
	SetTooltip(tooltip string)
	SetIcon(iconBytes []byte)
	AddMenuItem(title string, tooltip string) *systray.MenuItem
	AddMenuItemCheckbox(title string, tooltip string, checked bool) *systray.MenuItem
	Quit()
}

type DefaultSystrayWrapper struct{}

func (d *DefaultSystrayWrapper) Run(onReady func(), onExit func()) {
	systray.Run(onReady, onExit)
}

func (d *DefaultSystrayWrapper) SetTooltip(tooltip string) {
	systray.SetTooltip(tooltip)
}

func (d *DefaultSystrayWrapper) SetIcon(iconBytes []byte) {
	systray.SetIcon(iconBytes)
}

func (d *DefaultSystrayWrapper) AddMenuItem(title string, tooltip string) *systray.MenuItem {
	return systray.AddMenuItem(title, tooltip)
}

func (d *DefaultSystrayWrapper) AddMenuItemCheckbox(title string, tooltip string, checked bool) *systray.MenuItem {
	return systray.AddMenuItemCheckbox(title, tooltip, checked)
}

func (d *DefaultSystrayWrapper) Quit() {
	systray.Quit()
}
