package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	"log"
	"os/exec"
)

func DockerInstall(host *cluster.NodeCfg) {
	if checkDocker(host) == false {
		installDocker(host)
	}
}

func checkDocker(host *cluster.NodeCfg) bool {
	dockerCheckCmd := "which docker"
	if host == nil {
		if err := exec.Command("which", "docker").Run(); err != nil {
			return false
		} else {
			return true
		}
	} else {
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, dockerCheckCmd); err != nil {
			return false
		} else {
			return true
		}
	}
}

func installDocker(host *cluster.NodeCfg) {
	installDockerCmd := "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh"
	if host == nil {
		if output, err := exec.Command("curl", "https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh").CombinedOutput(); err != nil {
			log.Fatal("Install Docker Failed:\n")
			fmt.Println(output)
		}
	} else {
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, installDockerCmd); err != nil {
			log.Fatalf("Install Docker Failed (%s):\n", host.Address)
		}
	}
}
