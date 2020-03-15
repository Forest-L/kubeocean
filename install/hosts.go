package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func InstallDocker(node *cluster.ClusterNodeCfg) {
	installDockerCmd := "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh"
	if node.Node.Address == "" && CheckDocker(node) == false {
		log.Infof("Docker being installed ...")
		if output, err := exec.Command("/bin/sh", "-c", "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh").CombinedOutput(); err != nil {
			log.Fatal("Install Docker Failed:\n")
			fmt.Println(output)
		}
	} else {
		if CheckDocker(node) == false {
			log.Infof("Docker being installed... [%s]", node.Node.Address)
			err := ssh.CmdExec(node.Node.Address, node.Node.User, node.Node.Port, node.Node.Password, false, "", installDockerCmd)
			if err != nil {
				log.Fatalf("Install Docker Failed [%s]:\n", node.Node.Address)
				os.Exit(1)
			}
		}
	}
}

func CheckDocker(host *cluster.ClusterNodeCfg) bool {
	dockerCheckCmd := "which docker"
	if host.Node.InternalAddress == "" {
		if err := exec.Command("which", "docker").Run(); err != nil {
			return false
		} else {
			log.Info("Docker already exists.")
			return true
		}
	} else {
		if err := host.CmdExec(dockerCheckCmd); err != nil {
			return false
		} else {
			log.Infof("Docker already exists. (%s)", host.Node.Address)
			return true
		}
	}
}

func BootStrapOS(node *cluster.ClusterNodeCfg) {
	src := "/tmp/kubeocean/bootStrapScript.sh"
	dst := "/tmp/kubeocean"

	if node.Node.Address == "" {
		log.Info("BootStrapOS")
		if err := exec.Command(src).Run(); err != nil {
			log.Errorf("Bootstrap is Failed: %v", err)
		}
	} else {
		log.Infof("BootStrapOS [%s]", node.Node.InternalAddress)
		node.CmdExec("mkdir -p /tmp/kubeocean")
		ssh.PushFile(node.Node.Address, src, dst, node.Node.User, node.Node.Port, node.Node.Password, true)
		if out, err := node.CmdExecOut(src); err != nil {
			fmt.Println(out)
			log.Fatalf("Bootstrap is Failed [%s]:\n", node.Node.Address)
		}
	}
}

func GetKubeBinary(cfg *cluster.ClusterCfg, node *cluster.ClusterNodeCfg) {
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

	if node.Node.Address == "" {
		log.Info("Get Kube Binary Files")
		if err := exec.Command("/bin/sh", "-c", getKubeadmCmd).Run(); err != nil {
			log.Errorf("Failed to get kubeadm: %v", err)
		}
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubeadm").Run()

		if err := exec.Command("/bin/sh", "-c", getKubeletCmd).Run(); err != nil {
			log.Errorf("Failed to get kubelet: %v", err)
		}
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubelet").Run()
		exec.Command("/bin/sh", "-c", "ln -s /usr/local/bin/kubelet /usr/bin/kubelet").Run()

		if err := exec.Command("/bin/sh", "-c", getKubectlCmd).Run(); err != nil {
			log.Errorf("Failed to get kubectl: %v", err)
		}
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubectl").Run()

		exec.Command("/bin/sh", "-c", "mkdir -p /opt/cni/bin").Run()
		if err := exec.Command("/bin/sh", "-c", getKubeCniCmd).Run(); err != nil {
			log.Errorf("Failed to get kubecni: %v", err)
		}

	} else {
		log.Info("Get Kube Binary Files [%s]", node.Node.InternalAddress)
		ssh.PushFile(node.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeadmFile), "/tmp/kubeocean", node.Node.User, node.Node.Port, node.Node.Password, true)
		if err := node.CmdExec(getKubeadmCmd); err != nil {
			log.Fatalf("Failed to get kubeadm [%s]:\n", node.Node.Address)
		}
		node.CmdExec("chmod +x /usr/local/bin/kubeadm")

		ssh.PushFile(node.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeletFile), "/tmp/kubeocean", node.Node.User, node.Node.Port, node.Node.Password, true)
		if err := node.CmdExec(getKubeletCmd); err != nil {
			log.Fatalf("Failed to get kubelet [%s]:\n", node.Node.Address)
		}
		node.CmdExec("chmod +x /usr/local/bin/kubelet")
		node.CmdExec("ln -s /usr/local/bin/kubelet /usr/bin/kubelet")

		ssh.PushFile(node.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubectlFile), "/tmp/kubeocean", node.Node.User, node.Node.Port, node.Node.Password, true)
		if err := node.CmdExec(getKubectlCmd); err != nil {
			log.Fatalf("Failed to get kubectl [%s]:\n", node.Node.Address)
		}
		node.CmdExec("chmod +x /usr/local/bin/kubectl")

		ssh.PushFile(node.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeCniFile), "/tmp/kubeocean", node.Node.User, node.Node.Port, node.Node.Password, true)
		node.CmdExec("mkdir -p /opt/cni/bin")
		if err := node.CmdExec(getKubeCniCmd); err != nil {
			log.Fatalf("Failed to get kubecni [%s]:\n", node.Node.Address)
		}
	}
}

func SetKubeletService(node *cluster.ClusterNodeCfg) {

	if node.Node.Address == "" {
		log.Info("Set Kubelet Service")
		exec.Command("/bin/sh", "-c", "mkdir -p /etc/systemd/system/kubelet.service.d").Run()
		exec.Command("/bin/sh", "-c", "cp -f /tmp/kubeocean/kubelet.service /etc/systemd/system").Run()
		exec.Command("/bin/sh", "-c", "cp -f /tmp/kubeocean/10-kubeadm.conf /etc/systemd/system/kubelet.service.d").Run()
	} else {
		log.Info("Set Kubelet Service [%s]", node.Node.InternalAddress)
		node.CmdExec("mkdir -p /etc/systemd/system/kubelet.service.d")
		ssh.PushFile(node.Node.Address, "/tmp/kubeocean/kubelet.service", "/etc/systemd/system", node.Node.User, node.Node.Port, node.Node.Password, true)
		ssh.PushFile(node.Node.Address, "/tmp/kubeocean/10-kubeadm.conf", "/etc/systemd/system/kubelet.service.d", node.Node.User, node.Node.Port, node.Node.Password, true)
	}
}

func OverrideHostname(node *cluster.ClusterNodeCfg) {
	if node.Node.Address == "" {
		log.Info("Override Hostname")
		err := exec.Command("/bin/sh", "-c", fmt.Sprintf("hostnamectl set-hostname %s", cluster.DefaultHostName)).Run()
		if err != nil {
			log.Fatalf("Failed to Override Hostname: %v", err)
		}
	} else {
		log.Info("Override Hostname [%s]", node.Node.InternalAddress)
		node.CmdExec(fmt.Sprintf("hostnamectl set-hostname %s", node.Node.HostName))
	}
}
