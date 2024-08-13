package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

var currentOS = runtime.GOOS
var execCommand = exec.Command

type Impl struct{}

func (u *Impl) OpenDirectory(path string) error {
	var cmd *exec.Cmd
	switch currentOS {
	case "windows":
		cmd = execCommand("explorer", path)
	case "darwin":
		cmd = execCommand("open", path)
	case "linux":
		cmd = execCommand("xdg-open", path)
	default:
		return fmt.Errorf("unsupported operating system: %s", currentOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open directory: %w, output: %s", err, output)
	}

	return nil
}

func (u *Impl) OpenFile(filePath string) error {
	var cmd *exec.Cmd
	switch currentOS {
	case "darwin":
		cmd = execCommand("open", filePath)
	case "linux":
		cmd = execCommand("xdg-open", filePath)
	case "windows":
		cmd = execCommand("rundll32", "url.dll,FileProtocolHandler", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	return nil
}
