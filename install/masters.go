package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func InitCluster(cfg *cluster.ClusterCfg, master *cluster.ClusterNodeCfg) {
	if master.Node.InternalAddress == "" {
		exec.Command("sh", "-c", "mkdir -p /etc/kubernetes").Run()
		exec.Command("sh", "-c", "cp -f /tmp/kubeocean/kubeadm-config.yaml /etc/kubernetes/kubeadm-config.yaml").Run()
		if out, err := exec.Command("sh", "-c", "/usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml").CombinedOutput(); err != nil {
			log.Fatalf("Failed to init cluster:\n %v", string(out))
		}

		GetKubeConfig(master)

		tmpl.GenerateNetworkPluginFiles(cfg)
		getNetworkPluginFileCmd := "cp -f /tmp/kubeocean/calico.yaml /etc/kubernetes/calico.yaml"
		deployNetworkPluginCmd := fmt.Sprintf("/usr/local/bin/kubectl apply -f /etc/kubernetes/%s.yaml", "calico")
		if err := exec.Command("sh", "-c", getNetworkPluginFileCmd).Run(); err != nil {
			log.Fatalf("Failed to generate network plugin file")
		}
		if err := exec.Command("sh", "-c", deployNetworkPluginCmd).Run(); err != nil {
			log.Fatalf("Failed to deploy network plugin")
		}

	} else {
		master.CmdExec("mkdir -p /etc/kubernetes")
		ssh.PushFile(master.Node.Address, "/tmp/kubeocean/kubeadm-config.yaml", "/tmp/kubeocean", master.Node.User, master.Node.Port, master.Node.Password, true)
		master.CmdExec("cp -f /tmp/kubeocean/kubeadm-config.yaml /etc/kubernetes/kubeadm-config.yaml")
		initClusterCmd := "/usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml"
		if out, err := master.CmdExecOut(initClusterCmd); err != nil {
			fmt.Println(out)
			log.Fatalf("Failed to init cluster (%s):\n%v", master.Node.Address, err)
		}

		GetKubeConfig(master)
		configFile := fmt.Sprintf("/home/%s/.kube/config", master.Node.User)
		if master.Node.User == "root" {
			configFile = "/root/.kube/config"
		}
		ssh.PullFile(master.Node.Address, "/tmp/kubeocean", configFile, master.Node.User, master.Node.Port, master.Node.Password, true)
		deployNetworkPluginCmd := fmt.Sprintf("/usr/local/bin/kubectl apply -f /etc/kubernetes/%s.yaml", cfg.Network.Plugin)
		if cfg.Network.Plugin == "calico" {
			ssh.PushFile(master.Node.Address, "/tmp/kubeocean/calico.yaml", "/tmp/kubeocean", master.Node.User, master.Node.Port, master.Node.Password, true)
		}
		if cfg.Network.Plugin == "flannel" {
			ssh.PushFile(master.Node.Address, "/tmp/kubeocean/flannelyaml", "/tmp/kubeocean", master.Node.User, master.Node.Port, master.Node.Password, true)
		}
		master.CmdExec("cp -f /tmp/kubeocean/calico.yaml /etc/kubernetes/calico.yaml")
		if err := master.CmdExec(deployNetworkPluginCmd); err != nil {
			log.Fatalf("Failed to deploy calico (%s):\n", master.Node.Address)
		}
		if master.IsWorker {
			RemoveMasterTaint(master)
		}
	}
}

func RemoveMasterTaint(master *cluster.ClusterNodeCfg) {
	var removeMasterTaint string
	if master.Node.InternalAddress == "" {
		removeMasterTaint = fmt.Sprintf("/usr/local/bin/kubectl taint nodes %s node-role.kubernetes.io/master=:NoSchedule-", cluster.DefaultHostName)
		exec.Command("sh", "-c", removeMasterTaint).Run()
	} else {
		removeMasterTaint = fmt.Sprintf("/usr/local/bin/kubectl taint nodes %s node-role.kubernetes.io/master=:NoSchedule-", master.Node.HostName)
		master.CmdExec(removeMasterTaint)
		addWorkerLabel := fmt.Sprintf("/usr/local/bin/kubectl label node %s node-role.kubernetes.io/worker=", master.Node.HostName)
		master.CmdExec(addWorkerLabel)
	}
}

func GetKubeConfig(master *cluster.ClusterNodeCfg) {
	if master.Node.InternalAddress == "" {
		createConfigDirCmd := "mkdir -p /root/.kube"
		getKubeConfigCmd := "cp -f /etc/kubernetes/admin.conf /root/.kube/config"
		if err := exec.Command("sh", "-c", createConfigDirCmd).Run(); err != nil {
			log.Fatalf("Failed to generate kubeconfig")
		}
		if err := exec.Command("sh", "-c", getKubeConfigCmd).Run(); err != nil {
			log.Fatalf("Failed to generate kubeconfig")
		}
	} else {
		createConfigDirCmd := "mkdir -p /root/.kube && mkdir -p $HOME/.kube"
		getKubeConfigCmd := "cp -f /etc/kubernetes/admin.conf /root/.kube/config"
		getKubeConfigCmdUsr := "cp -f /etc/kubernetes/admin.conf $HOME/.kube/config"
		chownKubeConfig := "chown $(id -u):$(id -g) $HOME/.kube/config"
		if err := master.CmdExec(createConfigDirCmd); err != nil {
			log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
		}

		if err := master.CmdExec(getKubeConfigCmd); err != nil {
			log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
		}
		if err := master.CmdExec(getKubeConfigCmdUsr); err != nil {
			log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
		}
		if err := master.CmdExec(chownKubeConfig); err != nil {
			log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
		}
	}
}
