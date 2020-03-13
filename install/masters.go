package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func SetKubeadmCfg(cfg *cluster.ClusterCfg) {
	tmpl.GenerateKubeadmFiles(cfg)
}

func InitCluster(cfg *cluster.ClusterCfg, masters *cluster.MasterNodes) {
	SetKubeadmCfg(cfg)
	if masters == nil {
		exec.Command("sh", "-c", "cp -f /tmp/kubeocean/kubeadm-config.yaml /etc/kubernetes").CombinedOutput()
		if out, err := exec.Command("sh", "-c", "/usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml").CombinedOutput(); err != nil {
			log.Fatalf("Failed to init cluster:\n %v", out)
		}

		createConfigDirCmd := "mkdir -p /root/.kube"
		getKubeConfigCmd := "cp -f /etc/kubernetes/admin.conf /root/.kube/config"
		if err := exec.Command("sh", "-c", createConfigDirCmd).Run(); err != nil {
			log.Fatalf("Failed to generate kubeconfig")
		}
		if err := exec.Command("sh", "-c", getKubeConfigCmd).Run(); err != nil {
			log.Fatalf("Failed to generate kubeconfig")
		}

		tmpl.GenerateNetworkPluginFiles(cfg)
		getNetworkPluginFileCmd := "cp -f /tmp/kubeocean/calico.yaml /etc/kubernetes"
		deployNetworkPluginCmd := fmt.Sprintf("/usr/local/bin/kubectl apply -f /etc/kubernetes/%s.yaml", "calico")
		if err := exec.Command("sh", "-c", getNetworkPluginFileCmd).Run(); err != nil {
			log.Fatalf("Failed to generate network plugin file")
		}
		if err := exec.Command("sh", "-c", deployNetworkPluginCmd).Run(); err != nil {
			log.Fatalf("Failed to deploy network plugin")
		}

	} else {
		for index, master := range masters.Hosts {
			if index == 0 {
				ssh.PushFile(master.Node.Address, "/tmp/kubeocean/kubeadm-config.yaml", "/etc/kubernetes", master.Node.User, master.Node.Port, master.Node.Password, true)
				initClusterCmd := "/usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml"
				if err := master.CmdExec(initClusterCmd); err != nil {
					log.Fatalf("Failed to init cluster (%s):\n", master.Node.Address)
				}

				createConfigDirCmd := "mkdir -p /root/.kube"
				getKubeConfigCmd := "cp -f /etc/kubernetes/admin.conf /root/.kube/config"
				if err := master.CmdExec(createConfigDirCmd); err != nil {
					log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
				}

				if err := master.CmdExec(getKubeConfigCmd); err != nil {
					log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
				}
				tmpl.GenerateNetworkPluginFiles(cfg)
				deployNetworkPluginCmd := fmt.Sprintf("/usr/local/bin/kubectl apply -f /etc/kubernetes/%s.yaml", cfg.Network.Plugin)
				if cfg.Network.Plugin == "calico" {
					ssh.PushFile(master.Node.Address, "/tmp/kubeocean/calico.yaml", "/etc/kubernetes", master.Node.User, master.Node.Port, master.Node.Password, true)
				}
				if cfg.Network.Plugin == "flannel" {
					ssh.PushFile(master.Node.Address, "/tmp/kubeocean/flannelyaml", "/etc/kubernetes", master.Node.User, master.Node.Port, master.Node.Password, true)
				}
				if err := master.CmdExec(deployNetworkPluginCmd); err != nil {
					log.Fatalf("Failed to deploy calico (%s):\n", master.Node.Address)
				}
			}
		}
	}
}

func RemoveMasterTaint(masters *cluster.MasterNodes) {
	var removeMasterTaint string
	if masters == nil {
		removeMasterTaint = fmt.Sprintf("/usr/local/bin/kubectl taint nodes %s node-role.kubernetes.io/master=:NoSchedule-", cluster.DefaultHostName)
		exec.Command("sh", "-c", removeMasterTaint).Run()
	} else {
		for _, master := range masters.Hosts {
			if master.IsWorker == true {
				removeMasterTaint = fmt.Sprintf("/usr/local/bin/kubectl taint nodes %s node-role.kubernetes.io/master=:NoSchedule-", masters.Hosts[0].Node.HostName)
				master.CmdExec(removeMasterTaint)
			}
		}
	}
}
