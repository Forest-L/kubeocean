package install

import (
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
		exec.Command("sudo", "cp -f /tmp/kubeocean/kubeadm-config.yaml /etc/kubernetes").CombinedOutput()
		if out, err := exec.Command("sudo", " /usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml").CombinedOutput(); err != nil {
			log.Fatalf("Failed to init cluster:\n %v", out)
		}
	} else {
		for index, master := range masters.Hosts {
			if index == 0 {
				ssh.PushFile(master.Node.Address, "/tmp/kubeocean/kubeadm-config.yaml", "/etc/kubernetes", master.Node.User, master.Node.Port, master.Node.Password, true)
				initClusterCmd := "sudo /usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml"
				if err := ssh.CmdExec(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, initClusterCmd); err != nil {
					log.Fatalf("Failed to init cluster (%s):\n", master.Node.Address)
				}

				getKubeConfigCmd := "sudo mkdir -p /root/.kube && sudo cp -f /etc/kubernetes/admin.conf /root/.kube/config"
				if err := ssh.CmdExec(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, getKubeConfigCmd); err != nil {
					log.Fatalf("Failed to generate kubeconfig (%s):\n", master.Node.Address)
				}

				deployCalicoCmd := "sudo kubectl apply -f https://docs.projectcalico.org/v3.8/manifests/calico.yaml"
				if err := ssh.CmdExec(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, deployCalicoCmd); err != nil {
					log.Fatalf("Failed to deploy calico (%s):\n", master.Node.Address)
				}
			}
		}
	}
}
