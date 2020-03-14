package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func InjectHosts(cfg *cluster.ClusterCfg, nodes *cluster.AllNodes) {
	hosts := cfg.GenerateHosts()
	injectHostsCmd := fmt.Sprintf("echo \"%s\"  >> /etc/hosts", hosts)
	removeDuplicatesCmd := "awk ' !x[$0]++{print > \"/etc/hosts\"}' /etc/hosts"
	if nodes.Hosts == nil {
		if err := exec.Command("/bin/sh", "-c", injectHostsCmd).Run(); err != nil {
			log.Fatal("Failed to Inject Hosts:\n%v", err)
		}
		if err1 := exec.Command("/bin/sh", "-c", "awk", "' !x[$0]++{print > \"/etc/hosts\"}'", "/etc/hosts").Run(); err1 != nil {
			log.Fatalf("Failed to Inject Hosts:\n%v", err1)
		}
	} else {
		for _, host := range nodes.Hosts {
			if err := host.CmdExec(injectHostsCmd); err != nil {
				log.Fatal("Failed to Inject Hosts:\n%v", err)
			}
			host.CmdExec(removeDuplicatesCmd)
		}
	}
}

func DockerInstall(nodes *cluster.AllNodes) {
	installDockerCmd := "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh"
	if nodes.Hosts == nil && CheckDocker(nil) == false {
		log.Infof("Docker being installed ...")
		if output, err := exec.Command("/bin/sh", "-c", "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh").CombinedOutput(); err != nil {
			log.Fatal("Install Docker Failed:\n")
			fmt.Println(output)
		}
	} else {
		for _, host := range nodes.Hosts {
			if CheckDocker(&host) == false {
				log.Infof("Docker being installed... (%s)", host.Node.Address)
				out, err := ssh.CmdExecOut(host.Node.Address, host.Node.User, host.Node.Port, host.Node.Password, false, "", installDockerCmd)
				if err != nil {
					log.Fatalf("Install Docker Failed (%s):\n", host.Node.Address)
					fmt.Println(out)
					os.Exit(1)
				}
			}
		}
	}
}

func CheckDocker(host *cluster.ClusterNodeCfg) bool {
	dockerCheckCmd := "which docker"
	if host == nil {
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

func InitOS(nodes *cluster.AllNodes) {
	tmpl.GenerateBootStrapScript()
	src := "/tmp/kubeocean/bootStrapScript.sh"
	dst := "/tmp/kubeocean"

	if nodes.Hosts == nil {
		if err := exec.Command(src).Run(); err != nil {
			log.Errorf("Bootstrap is Failed: %v", err)
		}
	} else {
		for _, host := range nodes.Hosts {
			host.CmdExec("mkdir -p /tmp/kubeocean")
			ssh.PushFile(host.Node.Address, src, dst, host.Node.User, host.Node.Port, host.Node.Password, true)
			if err := host.CmdExec(src); err != nil {
				log.Fatalf("Bootstrap is Failed (%s):\n", host.Node.Address)
			}
		}
	}
}

func GetKubeBinary(cfg *cluster.ClusterCfg, nodes *cluster.AllNodes) {

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

	if nodes.Hosts == nil {
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
		for _, host := range nodes.Hosts {
			ssh.PushFile(host.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeadmFile), "/tmp/kubeocean", host.Node.User, host.Node.Port, host.Node.Password, true)
			if err := host.CmdExec(getKubeadmCmd); err != nil {
				log.Fatalf("Failed to get kubeadm (%s):\n", host.Node.Address)
			}
			host.CmdExec("chmod +x /usr/local/bin/kubeadm")

			ssh.PushFile(host.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeletFile), "/tmp/kubeocean", host.Node.User, host.Node.Port, host.Node.Password, true)
			if err := host.CmdExec(getKubeletCmd); err != nil {
				log.Fatalf("Failed to get kubelet (%s):\n", host.Node.Address)
			}
			host.CmdExec("chmod +x /usr/local/bin/kubelet")
			host.CmdExec("ln -s /usr/local/bin/kubelet /usr/bin/kubelet")

			ssh.PushFile(host.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubectlFile), "/tmp/kubeocean", host.Node.User, host.Node.Port, host.Node.Password, true)
			if err := host.CmdExec(getKubectlCmd); err != nil {
				log.Fatalf("Failed to get kubectl (%s):\n", host.Node.Address)
			}
			host.CmdExec("chmod +x /usr/local/bin/kubectl")

			ssh.PushFile(host.Node.Address, fmt.Sprintf("/tmp/kubeocean/%s", kubeCniFile), "/tmp/kubeocean", host.Node.User, host.Node.Port, host.Node.Password, true)
			host.CmdExec("mkdir -p /opt/cni/bin")
			if err := host.CmdExec(getKubeCniCmd); err != nil {
				log.Fatalf("Failed to get kubecni (%s):\n", host.Node.Address)
			}
		}
	}
}

func SetKubeletService(nodes *cluster.AllNodes) {
	tmpl.GenerateKubeletFiles()
	if nodes.Hosts != nil {
		for _, host := range nodes.Hosts {
			host.CmdExec("mkdir -p /etc/systemd/system/kubelet.service.d")
			ssh.PushFile(host.Node.Address, "/etc/systemd/system/kubelet.service", "/etc/systemd/system", host.Node.User, host.Node.Port, host.Node.Password, true)
			ssh.PushFile(host.Node.Address, "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf", "/etc/systemd/system/kubelet.service.d", host.Node.User, host.Node.Port, host.Node.Password, true)
		}
	}
}

func OverrideHostname(nodes *cluster.AllNodes) {
	if nodes.Hosts == nil {
		err := exec.Command("/bin/sh", "-c", "hostnamectl", "set-hostname", cluster.DefaultHostName)
		if err != nil {
			log.Fatalf("Failed to Override Hostname: %v", err)
		}
	} else {
		for _, node := range nodes.Hosts {
			node.CmdExec(fmt.Sprintf("hostnamectl set-hostname %s", node.Node.HostName))
		}
	}
}
