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
		// Fallback for other Unix-like systems
		cmd = exec.Command("sh", "-c", fmt.Sprintf("open '%s' || xdg-open '%s' || x-www-browser '%s'", path, path, path))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open directory: %w, output: %s", err, output)
	}

	return nil
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
