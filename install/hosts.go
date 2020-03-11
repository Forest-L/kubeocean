package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func InjectHosts(cfg *cluster.ClusterCfg) {
	hosts := cfg.GenerateHosts()
	injectHostsCmd := fmt.Sprintf("echo \"%s\"  >> /etc/hosts", hosts)
	removeDuplicatesCmd := "awk ' !x[$0]++{print > \"/etc/hosts\"}' /etc/hosts"
	if cfg.Hosts == nil {
		if err := exec.Command("sudo", injectHostsCmd).Run(); err != nil {
			log.Fatal("Failed to Inject Hosts:\n")
		}
		exec.Command("sudo", removeDuplicatesCmd).Run()
	} else {
		for _, host := range cfg.Hosts {
			if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, fmt.Sprintf("sudo %s", injectHostsCmd)); err != nil {
				log.Fatal("Failed to Inject Hosts:\n")
			}
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, fmt.Sprintf("sudo %s", removeDuplicatesCmd))
		}
	}
}

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

func GetKubeBinary(host *cluster.NodeCfg, version string) {

	var kubeVersion string
	if version == "" {
		kubeVersion = cluster.DefaultKubeVersion
	} else {
		kubeVersion = version
	}
	kubeadmUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/amd64/kubeadm", kubeVersion)
	kubeletUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/amd64/kubelet", kubeVersion)
	kubectlUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/amd64/kubectl", kubeVersion)
	kubeCniUrl := fmt.Sprintf("https://containernetworking.pek3b.qingstor.com/plugins/releases/download/%s/cni-plugins-linux-amd64-%s.tgz", "v0.8.1", "v0.8.1")
	getKubeadmCmd := fmt.Sprintf("curl -o /tmp/kubeocean/kubeadm  %s", kubeadmUrl)
	getKubeletCmd := fmt.Sprintf("curl -o /tmp/kubeocean/kubelet  %s", kubeletUrl)
	getKubectlCmd := fmt.Sprintf("curl -o /tmp/kubeocean/kubectl  %s", kubectlUrl)
	getKubeCniCmd := fmt.Sprintf("curl -o /tmp/kubeocean/cni-plugins-linux-amd64-v0.8.1.tgz  %s", kubeCniUrl)
	if host == nil {
		if err := exec.Command("/bin/sh", "-c", getKubeadmCmd).Run(); err != nil {
			log.Errorf("Failed to get kubeadm: %v", err)
		}
		exec.Command("/bin/sh", "-c", "cp /tmp/kubeocean/kubeadm /usr/local/bin/kubeadm").Run()
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubeadm").Run()

		if err := exec.Command("/bin/sh", "-c", getKubeletCmd).Run(); err != nil {
			log.Errorf("Failed to get kubelet: %v", err)
		}
		exec.Command("/bin/sh", "-c", "cp /tmp/kubeocean/kubelet /usr/local/bin/kubelet").Run()
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubelet").Run()
		exec.Command("/bin/sh", "-c", "ln -s /usr/local/bin/kubelet /usr/bin/kubelet").Run()

		if err := exec.Command("/bin/sh", "-c", getKubectlCmd).Run(); err != nil {
			log.Errorf("Failed to get kubectl: %v", err)
		}
		exec.Command("/bin/sh", "-c", "cp /tmp/kubeocean/kubectl /usr/local/bin/kubectl").Run()
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubectl").Run()

		if err := exec.Command("/bin/sh", "-c", getKubeCniCmd).Run(); err != nil {
			log.Errorf("Failed to get kubecni: %v", err)
		}
		exec.Command("/bin/sh", "-c", "cp /tmp/kubeocean/kubectl /usr/local/bin/kubectl").Run()
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubectl").Run()

		if err := exec.Command("/bin/sh", "-c", getKubeCniCmd).Run(); err != nil {
			log.Errorf("Failed to get kubecni: %v", err)
		}
		exec.Command("/bin/sh", "-c", "mkdir -p /opt/cni/bin").Run()
		exec.Command("/bin/sh", "-c", "tar -zxf /tmp/kubeocean/cni-plugins-linux-amd64-v0.8.1.tgz -C /opt/cni/bin").Run()
	} else {
		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, getKubeadmCmd); err != nil {
			log.Fatalf("Failed to get kubeadm (%s):\n", host.Address)
		}
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "cp /tmp/kubeocean/kubeadm /usr/local/bin/kubeadm")
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "chmod +x /usr/local/bin/kubeadm")

		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, getKubeletCmd); err != nil {
			log.Fatalf("Failed to get kubelet (%s):\n", host.Address)
		}
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "cp /tmp/kubeocean/kubelet /usr/local/bin/kubelet")
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "chmod +x /usr/local/bin/kubelet")
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "ln -s /usr/local/bin/kubelet /usr/bin/kubelet")

		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, getKubectlCmd); err != nil {
			log.Fatalf("Failed to get kubectl (%s):\n", host.Address)
		}
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "cp /tmp/kubeocean/kubectl /usr/local/bin/kubectl")
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "chmod +x /usr/local/bin/kubectl")

		if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, getKubeCniCmd); err != nil {
			log.Fatalf("Failed to get kubecni (%s):\n", host.Address)
		}
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "mkdir -p /opt/cni/bin")
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "tar -zxf /tmp/kubeocean/cni-plugins-linux-amd64-v0.8.1.tgz -C /opt/cni/bin")
	}
}

func SetKubeletService(host *cluster.NodeCfg, repo string, version string) {
	tmpl.GenerateKubeletFiles(repo, version)

	if host != nil {
		ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "mkdir -p /etc/systemd/system/kubelet.service.d")
		ssh.PushFile(host.Address, "/etc/systemd/system/kubelet.service", "/etc/systemd/system", host.User, host.Port, host.Password, true)
		ssh.PushFile(host.Address, "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf", "/etc/systemd/system/kubelet.service.d", host.User, host.Port, host.Password, true)
	}
}
