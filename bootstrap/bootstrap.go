package bootstrap

import (
	"fmt"
	"os/exec"
)

func SystemInit() error {

	if err := exec.Command("sysctl", "net.bridge.bridge-nf-call-iptables").Run(); err != nil {
		return fmt.Errorf("failed to sysctl: %v", err)
	}

	if err := exec.Command("swapoff", "-a").Run(); err != nil {
		return fmt.Errorf("failed to swapoff: %v", err)
	}

	return nil
}
