package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenDirectory(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func OpenFile(filePath string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", filePath).Start()
	case "linux":
		return exec.Command("xdg-open", filePath).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
