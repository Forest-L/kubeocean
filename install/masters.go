package install

import (
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func SetKubeadmService(cfg *cluster.ClusterCfg) {
	tmpl.GenerateKubeadmFiles(cfg)
}

func InitCluster(cfg *cluster.ClusterCfg, masters *cluster.MasterNodes) {
	if masters == nil {
		SetKubeadmService(cfg)
		if out, err := exec.Command("/usr/local/bin/kubeadm", "init --config=/etc/kubernetes/kubeadm-config.yaml").CombinedOutput(); err != nil {
			log.Fatalf("Failed to init cluster:\n %v", out)
		}
	} else {
		for index, master := range masters.Hosts {
			if index == 0 {
				SetKubeadmService(cfg)
				initClusterCmd := "/usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yam"
				if err := ssh.CmdExec(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, initClusterCmd); err != nil {
					log.Fatalf("Failed to init cluster (%s):\n", master.Node.Address)
				}

				getKubeConfigCmd := "mkdir -p /root/.kube && cp -f /etc/kubernetes/admin.conf /root/.kube/config"
				if err := ssh.CmdExec(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, getKubeConfigCmd); err != nil {
					log.Fatalf("Failed to init cluster (%s):\n", master.Node.Address)
				}

				deployCalicoCmd := "kubectl apply -f https://docs.projectcalico.org/v3.8/manifests/calico.yaml"
				if err := exec.Command("/bin/sh", "-c", deployCalicoCmd).Run(); err != nil {
					log.Fatalf("Failed to deploy calico: %v", err)
				}
			}
		}
	}
}
