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
		rs := make(chan string)
		cn := make(chan struct{}, ssh.DefaultCon)
		hostNum := len(nodes.Hosts)
		wg := &sync.WaitGroup{}
		defer close(rs)
		defer close(cn)
		go ssh.CheckResults(rs, hostNum, wg, cn)
		for _, node := range nodes.Hosts {
			host := node
			cn <- struct{}{}
			wg.Add(1)
			go func(rs chan string, host *cluster.ClusterNodeCfg) {
				cmd := fmt.Sprintf("hostnamectl set-hostname %s", host.Node.HostName)
				host.CmdExec(cmd)
				rs <- host.Node.InternalAddress
			}(rs, &host)
		}
		wg.Wait()
	}
}

func InstallDocker(nodes *cluster.AllNodes) {

	if nodes.Hosts[0].Node.InternalAddress == "" {
		if err := exec.Command("which", "docker").Run(); err != nil {
			log.Infof("Docker being installed ...")
			if err := exec.Command("/bin/sh", "-c", "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh").Run(); err != nil {
				log.Fatal("Install Docker Failed:\n")
			}
			exec.Command("sh", "-c", "systemctl enable docker").Run()
		}
	} else {
		log.Infof("Docker being installed")
		rs := make(chan string)
		cn := make(chan struct{}, ssh.DefaultCon)
		hostNum := len(nodes.Hosts)
		defer close(rs)
		defer close(cn)
		wg := &sync.WaitGroup{}
		go ssh.CheckResults(rs, hostNum, wg, cn)

		for _, node := range nodes.Hosts {
			host := node
			cn <- struct{}{}
			wg.Add(1)
			go func(rs chan string, host *cluster.ClusterNodeCfg) {
				dockerCheckCmd := "which docker"
				installDockerCmd := "curl https://raw.githubusercontent.com/pixiake/kubeocean/master/scripts/docker-istall.sh | sh"
				if err := host.CmdExec(dockerCheckCmd); err != nil {
					ssh.CmdExec(host.Node.Address, host.Node.User, host.Node.Port, host.Node.Password, true, "", installDockerCmd)
					host.CmdExec("systemctl enable docker")
				} else {
					log.Infof("Docker already exists. [%s]", host.Node.InternalAddress)
				}
				rs <- host.Node.InternalAddress
			}(rs, &host)
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
	kubeCniFile := fmt.Sprintf("cni-plugins-linux-%s-%s.tgz", cluster.DefaultArch, cluster.DefaultCniVersion)
	helmFile := fmt.Sprintf("helm-%s", cluster.DefaultHelmVersion)
	getKubeadmCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/kubeadm", kubeadmFile)
	getKubeletCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/kubelet", kubeletFile)
	getKubectlCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/kubectl", kubectlFile)
	getKubeCniCmd := fmt.Sprintf("tar -zxf /tmp/kubeocean/%s -C /opt/cni/bin", kubeCniFile)
	getHelmCmd := fmt.Sprintf("cp -f /tmp/kubeocean/%s /usr/local/bin/helm", helmFile)

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

		log.Info("Get Helm Binary Files")
		if err := exec.Command("/bin/sh", "-c", getHelmCmd).Run(); err != nil {
			log.Errorf("Failed to get helm: %v", err)
		}
		exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/helm").Run()

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

		nodes.GoPush(fmt.Sprintf("/tmp/kubeocean/%s", helmFile), "/tmp/kubeocean")
		nodes.GoExec(getHelmCmd)
		nodes.GoExec("chmod +x /usr/local/bin/helm")

	}
}

func SetKubeletService(nodes *cluster.AllNodes) {

	if nodes.Hosts[0].Node.InternalAddress == "" {
		log.Info("Set Kubelet Service")
		exec.Command("/bin/sh", "-c", "mkdir -p /etc/systemd/system/kubelet.service.d").Run()
		exec.Command("/bin/sh", "-c", "cp -f /tmp/kubeocean/kubelet.service /etc/systemd/system").Run()
		exec.Command("/bin/sh", "-c", "cp -f /tmp/kubeocean/10-kubeadm.conf /etc/systemd/system/kubelet.service.d").Run()
		exec.Command("/bin/sh", "-c", "systemctl enable kubelet").Run()
	} else {
		log.Info("Set Kubelet Service")
		nodes.GoExec("mkdir -p /etc/systemd/system/kubelet.service.d")
		nodes.GoPush("/tmp/kubeocean/kubelet.service", "/tmp/kubeocean")
		nodes.GoPush("/tmp/kubeocean/10-kubeadm.conf", "/tmp/kubeocean")
		nodes.GoExec(fmt.Sprintf("cp -f /tmp/kubeocean/kubelet.service /etc/systemd/system"))
		nodes.GoExec(fmt.Sprintf("cp -f /tmp/kubeocean/10-kubeadm.conf /etc/systemd/system/kubelet.service.d"))
		nodes.GoExec("systemctl enable kubelet")
	}
}
