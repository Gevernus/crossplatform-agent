package commands

import (
	"fmt"
	"os/exec"
	"runtime"
)

var currentOS = runtime.GOOS
var execCommand = exec.Command

func Shutdown() error {
	var cmd *exec.Cmd

	switch currentOS {
	case "windows":
		cmd = execCommand("shutdown", "/s", "/t", "0")
	case "darwin", "linux":
		cmd = execCommand("shutdown", "-h", "now")
	default:
		return fmt.Errorf("unsupported operating system: %s", currentOS)
	}

	return cmd.Run()
}

func ShutdownNetwork() error {
	var cmd *exec.Cmd

	switch currentOS {
	case "windows":
		cmd = execCommand("netsh", "interface", "set", "interface", "name=*", "admin=disabled")
	case "darwin":
		cmd = execCommand("networksetup", "-setairportpower", "en0", "off")
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = execCommand("networksetup", "-setnetworkserviceenabled", "Ethernet", "off")
	case "linux":
		cmd = execCommand("nmcli", "networking", "off")
	default:
		return fmt.Errorf("unsupported operating system: %s", currentOS)
	}

	return cmd.Run()
}
