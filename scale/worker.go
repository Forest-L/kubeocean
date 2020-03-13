package scale

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func JoinWorkers(workers *cluster.WorkerNodes, masters *cluster.MasterNodes) {
	joinWorkerCmd := JoinWorkerCmd(masters)
	for _, worker := range workers.Hosts {
		if worker.IsMaster != true {
			if err := worker.CmdExec(joinWorkerCmd); err != nil {
				log.Fatalf("Failed to add master (%s):\n", worker.Node.Address)
			}
			addWorkerLabel := fmt.Sprintf("kubectl label node %s node-role.kubernetes.io/worker=", worker.Node.HostName)
			exec.Command("sh", "-c", addWorkerLabel).Run()
		}
	}
}
