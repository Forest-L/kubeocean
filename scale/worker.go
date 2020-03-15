package scale

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func JoinWorker(worker *cluster.ClusterNodeCfg, joinWorkerCmd string) {
	if err := worker.CmdExec(joinWorkerCmd); err != nil {
		log.Fatalf("Failed to add worker (%s):\n", worker.Node.Address)
	}

	GetKubeconfig(worker)
	addWorkerLabel := fmt.Sprintf("kubectl label node %s node-role.kubernetes.io/worker=", worker.Node.HostName)
	exec.Command("sh", "-c", addWorkerLabel).Run()
}

func GetKubeconfig(worker *cluster.ClusterNodeCfg) {
	configSrc := "/tmp/kubeocean/config"
	var congigDst string
	if worker.Node.User == "root" {
		congigDst = "/root/.kube"
	} else {
		congigDst = fmt.Sprintf("/home/%s/.kube", worker.Node.User)
	}

	createConfigDirCmd := "mkdir -p /root/.kube && mkdir -p $HOME/.kube"
	chownKubeDir := "chown $(id -u):$(id -g) $HOME/.kube"
	chownKubeConfig := "chown $(id -u):$(id -g) $HOME/.kube/config"

	if err := worker.CmdExec(createConfigDirCmd); err != nil {
		log.Fatalf("Failed to generate kubeconfig (%s):\n", worker.Node.Address)
	}
	worker.CmdExec(chownKubeDir)
	ssh.PushFile(worker.Node.Address, configSrc, congigDst, worker.Node.User, worker.Node.Port, worker.Node.Password, true)
	if err := worker.CmdExec(chownKubeConfig); err != nil {
		log.Fatalf("Failed to generate kubeconfig (%s):\n", worker.Node.Address)
	}
	if worker.Node.User != "root" {
		worker.CmdExec(fmt.Sprintf("cp -f /home/%s/.kube/config /root/.kube", worker.Node.User))
	}
}
