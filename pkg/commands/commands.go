package commands

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OSGetter defines a function type for getting the current OS
type OSGetter func() string

// DefaultOSGetter is the default implementation using runtime.GOOS
var DefaultOSGetter OSGetter = func() string {
	return runtime.GOOS
}

func Shutdown() error {
	var cmd *exec.Cmd

	switch DefaultOSGetter() {
	case "windows":
		cmd = DefaultCommandRunner("shutdown", "/s", "/t", "0")
	case "darwin":
		cmd = DefaultCommandRunner("shutdown", "-h", "now")
	case "linux":
		cmd = DefaultCommandRunner("shutdown", "-h", "now")
	default:
		return fmt.Errorf("unsupported operating system: %s", DefaultOSGetter())
	}

	return cmd.Run()
}

func ShutdownNetwork() error {
	var cmd *exec.Cmd

	switch DefaultOSGetter() {
	case "windows":
		cmd = DefaultCommandRunner("netsh", "interface", "set", "interface", "name=*", "admin=disabled")
	case "darwin":
		cmd = DefaultCommandRunner("networksetup", "-setairportpower", "en0", "off")
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = DefaultCommandRunner("networksetup", "-setnetworkserviceenabled", "Ethernet", "off")
	case "linux":
		cmd = DefaultCommandRunner("nmcli", "networking", "off")
	default:
		return fmt.Errorf("unsupported operating system: %s", DefaultOSGetter())
	}

	return cmd.Run()
}
