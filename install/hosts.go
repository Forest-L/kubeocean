package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"sync"
)

func BootStrapOS(nodes *cluster.AllNodes) {
	src := "/tmp/kubeocean/bootStrapScript.sh"
	dst := "/tmp/kubeocean"
	if nodes.Hosts[0].Node.InternalAddress == "" {
		log.Info("BootStrapOS")
		if err := exec.Command(src).Run(); err != nil {
			log.Errorf("Bootstrap is Failed: %v", err)
		}
	} else {
		log.Info("BootStrapOS")
		nodes.GoExec("mkdir -p /tmp/kubeocean -m 777")
		nodes.GoPush(src, dst)
		nodes.GoExec(src)
	}
}

func OverrideHostname(nodes *cluster.AllNodes) {
	if nodes.Hosts[0].Node.InternalAddress == "" {
		log.Info("Override Hostname")
		err := exec.Command("/bin/sh", "-c", fmt.Sprintf("hostnamectl set-hostname %s", cluster.DefaultHostName)).Run()
		if err != nil {
			log.Fatalf("Failed to Override Hostname: %v", err)
		}
	} else {
		log.Infof("Override Hostname")
		result := make(chan string)
		ccons := make(chan struct{}, ssh.DefaultCon)
		hostNum := len(nodes.Hosts)
		wg := &sync.WaitGroup{}
		go ssh.CheckResults(result, hostNum, wg, ccons)
		for _, node := range nodes.Hosts {
			ccons <- struct{}{}
			wg.Add(1)
			cmd := fmt.Sprintf("hostnamectl set-hostname %s", node.Node.HostName)
			go func(rs chan string, cmd string) {
				node.CmdExec(cmd)
				rs <- "ok"
			}(result, cmd)
		}
		wg.Wait()
	}
}

func InstallDocker(nodes *cluster.AllNodes) {
	dockerCheckCmd := "which docker"
	installDockerCmd := "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh"
	if nodes.Hosts[0].Node.InternalAddress == "" {
		if err := exec.Command("which", "docker").Run(); err != nil {
			log.Infof("Docker being installed ...")
			if err := exec.Command("/bin/sh", "-c", "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh").Run(); err != nil {
				log.Fatal("Install Docker Failed:\n")
			}
		}
	} else {
		log.Infof("Docker being installed")
		result := make(chan string)
		ccons := make(chan struct{}, ssh.DefaultCon)
		hostNum := len(nodes.Hosts)
		wg := &sync.WaitGroup{}
		go ssh.CheckResults(result, hostNum, wg, ccons)

		for _, node := range nodes.Hosts {
			ccons <- struct{}{}
			wg.Add(1)
			go func(host *cluster.ClusterNodeCfg, rs chan string) {
				if err := host.CmdExec(dockerCheckCmd); err != nil {
					ssh.CmdExec(host.Node.Address, host.Node.User, host.Node.Port, host.Node.Password, true, "", installDockerCmd)
				} else {
					log.Infof("Docker already exists. [%s]", host.Node.InternalAddress)
				}
				rs <- "ok"
			}(&node, result)
		}
		wg.Wait()
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

	if nodes.Hosts[0].Node.InternalAddress == "" {
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
		log.Info("Get Kube Binary Files")
		nodes.GoPush(fmt.Sprintf("/tmp/kubeocean/%s", kubeadmFile), "/tmp/kubeocean")
		nodes.GoExec(getKubeadmCmd)
		nodes.GoExec("chmod +x /usr/local/bin/kubeadm")

		nodes.GoPush(fmt.Sprintf("/tmp/kubeocean/%s", kubeletFile), "/tmp/kubeocean")
		nodes.GoExec(getKubeletCmd)
		nodes.GoExec("chmod +x /usr/local/bin/kubelet")
		nodes.GoExec("ln -s /usr/local/bin/kubelet /usr/bin/kubelet")

		nodes.GoPush(fmt.Sprintf("/tmp/kubeocean/%s", kubectlFile), "/tmp/kubeocean")
		nodes.GoExec(getKubectlCmd)
		nodes.GoExec("chmod +x /usr/local/bin/kubectl")

		nodes.GoPush(fmt.Sprintf("/tmp/kubeocean/%s", kubeCniFile), "/tmp/kubeocean")

		nodes.GoExec("mkdir -p /opt/cni/bin")
		nodes.GoExec(getKubeCniCmd)
	}
}

func SetKubeletService(nodes *cluster.AllNodes) {

	if nodes.Hosts[0].Node.InternalAddress == "" {
		log.Info("Set Kubelet Service")
		exec.Command("/bin/sh", "-c", "mkdir -p /etc/systemd/system/kubelet.service.d").Run()
		exec.Command("/bin/sh", "-c", "cp -f /tmp/kubeocean/kubelet.service /etc/systemd/system").Run()
		exec.Command("/bin/sh", "-c", "cp -f /tmp/kubeocean/10-kubeadm.conf /etc/systemd/system/kubelet.service.d").Run()
	} else {
		log.Info("Set Kubelet Service")
		nodes.GoExec("mkdir -p /etc/systemd/system/kubelet.service.d")
		nodes.GoPush("/tmp/kubeocean/kubelet.service", "/tmp/kubeocean")
		nodes.GoPush("/tmp/kubeocean/10-kubeadm.conf", "/tmp/kubeocean")
		nodes.GoExec(fmt.Sprintf("cp -f /tmp/kubeocean/kubelet.service /etc/systemd/system"))
		nodes.GoExec(fmt.Sprintf("cp -f /tmp/kubeocean/10-kubeadm.conf /etc/systemd/system/kubelet.service.d"))
	}
}
