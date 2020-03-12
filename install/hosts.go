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

func DockerInstall(cfg *cluster.ClusterCfg) {
	for _, host := range cfg.Hosts {
		if checkDocker(&host) == false {
			installDocker(&host)
		}
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

func InitOS(cfg *cluster.ClusterCfg) {
	log.Info("BootStrap")
	tmpl.GenerateBootStrapScript()
	src := "/tmp/kubeocean/bootStrapScript.sh"
	dst := "/tmp/kubeocean"

	if cfg.Hosts == nil {
		if err := exec.Command(src).Run(); err != nil {
			log.Errorf("Bootstrap is Failed: %v", err)
		}
	} else {
		for _, host := range cfg.Hosts {
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "mkdir -p /tmp/kubeocean")
			ssh.PushFile(host.Address, src, dst, host.User, host.Port, host.Password, true)
			if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, src); err != nil {
				log.Fatalf("Bootstrap is Failed (%s):\n", host.Address)
			}
		}
	}
}

func GetKubeBinary(cfg *cluster.ClusterCfg) {

	var kubeVersion string
	if cfg.KubeVersion == "" {
		kubeVersion = cluster.DefaultKubeVersion
	} else {
		kubeVersion = cfg.KubeVersion
	}
	kubeadmFile := fmt.Sprintf("kubeadm-%s", kubeVersion)
	kubeletFile := fmt.Sprintf("kubelet-%s", kubeVersion)
	kubectlFile := fmt.Sprintf("kubectl-%s", kubeVersion)
	kubeCniFile := fmt.Sprintf("cni-plugins-linux-%s-%s.tgz", cluster.DefaultArch, "v0.8.1")
	getKubeadmCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/kubeadm", kubeadmFile)
	getKubeletCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/kubelet", kubeletFile)
	getKubectlCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/kubectl", kubectlFile)
	getKubeCniCmd := fmt.Sprintf("tar -zxf /tmp/kubeocean/%s -C /opt/cni/bin", kubeCniFile)
	sudoGetKubeadmCmd := fmt.Sprintf("sudo %s", getKubeadmCmd)
	sudoGetKubeletCmd := fmt.Sprintf("sudo %s", getKubeletCmd)
	sudoGetKubectlCmd := fmt.Sprintf("sudo %s", getKubectlCmd)
	sudoGetKubeCniCmd := fmt.Sprintf("sudo %s", getKubeCniCmd)
	if cfg.Hosts == nil {
		if err := exec.Command("sudo", "-E", "/bin/sh", "-c", getKubeadmCmd).Run(); err != nil {
			log.Errorf("Failed to get kubeadm: %v", err)
		}
		exec.Command("sudo", "-E", "/bin/sh", "-c", "chmod +x /usr/local/bin/kubeadm").Run()

		if err := exec.Command("sudo", getKubeletCmd).Run(); err != nil {
			log.Errorf("Failed to get kubelet: %v", err)
		}
		exec.Command("sudo", "-E", "/bin/sh", "-c", "chmod +x /usr/local/bin/kubelet").Run()
		exec.Command("sudo", "-E", "/bin/sh", "-c", "ln -s /usr/local/bin/kubelet /usr/bin/kubelet").Run()

		if err := exec.Command("sudo", "-E", "/bin/sh", "-c", getKubectlCmd).Run(); err != nil {
			log.Errorf("Failed to get kubectl: %v", err)
		}
		exec.Command("sudo", "-E", "/bin/sh", "-c", "chmod +x /usr/local/bin/kubectl").Run()

		exec.Command("sudo", "-E", "/bin/sh", "-c", "mkdir -p /opt/cni/bin").Run()
		if err := exec.Command("sudo", "-E", "/bin/sh", "-c", getKubeCniCmd).Run(); err != nil {
			log.Errorf("Failed to get kubecni: %v", err)
		}

	} else {
		for _, host := range cfg.Hosts {
			ssh.PushFile(host.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeadmFile), "/tmp/kubeocean", host.User, host.Port, host.Password, true)
			if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, sudoGetKubeadmCmd); err != nil {
				log.Fatalf("Failed to get kubeadm (%s):\n", host.Address)
			}
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "sudo chmod +x /usr/local/bin/kubeadm")

			ssh.PushFile(host.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeletFile), "/tmp/kubeocean", host.User, host.Port, host.Password, true)
			if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, sudoGetKubeletCmd); err != nil {
				log.Fatalf("Failed to get kubelet (%s):\n", host.Address)
			}
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "sudo chmod +x /usr/local/bin/kubelet")
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "sudo ln -s /usr/local/bin/kubelet /usr/bin/kubelet")

			ssh.PushFile(host.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubectlFile), "/tmp/kubeocean", host.User, host.Port, host.Password, true)
			if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, sudoGetKubectlCmd); err != nil {
				log.Fatalf("Failed to get kubectl (%s):\n", host.Address)
			}
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "sudo chmod +x /usr/local/bin/kubectl")

			ssh.PushFile(host.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeCniFile), "/tmp/kubeocean", host.User, host.Port, host.Password, true)
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "sudo mkdir -p /opt/cni/bin")
			if err := ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, sudoGetKubeCniCmd); err != nil {
				log.Fatalf("Failed to get kubecni (%s):\n", host.Address)
			}
		}
	}
}

func SetKubeletService(cfg *cluster.ClusterCfg) {
	var repo, version string
	if cfg.KubeVersion == "" {
		version = cluster.DefaultKubeVersion
	} else {
		version = cfg.KubeVersion
	}
	if cfg.KubeImageRepo == "" {
		version = cluster.DefaultKubeImageRepo
	} else {
		version = cfg.KubeImageRepo
	}
	tmpl.GenerateKubeletFiles(repo, version)
	if cfg.Hosts != nil {
		for _, host := range cfg.Hosts {
			ssh.CmdExec(host.Address, host.User, host.Port, host.Password, false, "mkdir -p /etc/systemd/system/kubelet.service.d")
			ssh.PushFile(host.Address, "/etc/systemd/system/kubelet.service", "/etc/systemd/system", host.User, host.Port, host.Password, true)
			ssh.PushFile(host.Address, "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf", "/etc/systemd/system/kubelet.service.d", host.User, host.Port, host.Password, true)
		}
	}
}
