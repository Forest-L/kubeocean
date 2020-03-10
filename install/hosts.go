package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
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
			log.Info("Docker already exists.")
			return true
		}
	} else {
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, dockerCheckCmd); err != nil {
			return false
		} else {
			log.Infof("Docker already exists. (%s)", host.Address)
			return true
		}
	}
}

func installDocker(host *cluster.NodeCfg) {
	installDockerCmd := "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh"
	if host == nil {
		log.Infof("Docker being installed...")
		if output, err := exec.Command("curl", "https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh").CombinedOutput(); err != nil {
			log.Fatal("Install Docker Failed:\n")
			fmt.Println(output)
		}
	} else {
		log.Infof("Docker being installed... (%s)", host.Address)
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, installDockerCmd); err != nil {
			log.Fatalf("Install Docker Failed (%s):\n", host.Address)
		}
	}
}

func InitOS(host *cluster.NodeCfg) {
	log.Info("BootStrap")
	tmpl.GenerateBootStrapScript()
	src := "/tmp/kubeocean/bootStrapScript.sh"
	dst := "/tmp/kubeocean"
	initOsCmd := fmt.Sprintf("/bin/bash %s", src)
	if host == nil {
		if err := exec.Command("/bin/bash", src).Run(); err != nil {
			log.Errorf("Bootstrap is Failed: %v", err)
		}
	} else {
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "mkdir -p /tmp/kubeocean")
		ssh.PushFile(host.Address, src, dst, host.User, host.Port, host.Password, true)
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, initOsCmd); err != nil {
			log.Fatalf("Bootstrap is Failed (%s):\n", host.Address)
		}
	}
}

func PullHyperKubeImage(host *cluster.NodeCfg, repo string, version string) {
	var hyperKubeImage string
	if repo == "" && version == "" {
		hyperKubeImage = fmt.Sprintf("%s/google-containers/hyperkube:%s", cluster.DefaultKubeImageRepo, cluster.DefaultKubeVersion)
	} else {
		hyperKubeImage = fmt.Sprintf("%s/google-containers/hyperkube:%s", repo, version)
	}

	pullHyperKubeImageCmd := fmt.Sprintf("docker pull %s", hyperKubeImage)
	if host == nil {
		if err := exec.Command("/bin/sh", "-c", pullHyperKubeImageCmd).Run(); err != nil {
			log.Errorf("Failed to pull hyperKube image: %v", err)
		}
	} else {
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, pullHyperKubeImageCmd); err != nil {
			log.Fatalf("Failed to pull hyperKube image (%s):\n", host.Address)
		}
	}
}

func GetKubeadm(host *cluster.NodeCfg, version string) {

	var kubeVersion string
	if version == "" {
		kubeVersion = cluster.DefaultKubeVersion
	} else {
		kubeVersion = version
	}
	kubeadmUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/amd64/kubeadm", kubeVersion)
	getKubeadmCmd := fmt.Sprintf("curl -o /usr/local/bin/kubeadm  %s", kubeadmUrl)
	if host == nil {
		if err := exec.Command("/bin/sh", "-c", getKubeadmCmd).Run(); err != nil {
			log.Errorf("Failed to get kubeadm: %v", err)
		}
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubeadm").Run()

	} else {
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, getKubeadmCmd); err != nil {
			log.Fatalf("Failed to get kubeadm (%s):\n", host.Address)
		}
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "chmod +x /usr/local/bin/kubeadm")
	}
}

func SetKubeletService(host *cluster.NodeCfg, repo string, version string) {
	tmpl.GenerateKubeletFiles(repo, version)

	if host != nil {
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "mkdir -p /etc/systemd/system/kubelet.service.d")
		ssh.PushFile(host.Address, "/usr/local/bin/kubelet", "/usr/local/bin", host.User, host.Port, host.Password, true)
		ssh.PushFile(host.Address, "/etc/systemd/system/kubelet.service", "/etc/systemd/system", host.User, host.Port, host.Password, true)
		ssh.PushFile(host.Address, "/etc/systemd/system/kubelet.service.d/kubelet-contain.conf", "/etc/systemd/system/kubelet.service.d", host.User, host.Port, host.Password, true)
	}
}
