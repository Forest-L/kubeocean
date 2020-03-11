package scale

import (
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
)

func JoinWorkers(workers *cluster.WorkerNodes, masters *cluster.MasterNodes) {
	joinWorkerCmd := JoinWorkerCmd(masters)
	for _, worker := range workers.Hosts {
		if worker.IsMaster != true {
			if err := ssh.CmdExec(worker.Node.Address, worker.Node.User, worker.Node.Port, worker.Node.Password, false, joinWorkerCmd); err != nil {
				log.Fatalf("Failed to add master (%s):\n", worker.Node.Address)
			}
		}
	}
}
