package scale

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func JoinWorker(worker *cluster.ClusterNodeCfg, joinWorkerCmd string) {
	if err := worker.CmdExec(joinWorkerCmd); err != nil {
		log.Fatalf("Failed to add worker (%s):\n", worker.Node.Address)
	}
	addWorkerLabel := fmt.Sprintf("kubectl label node %s node-role.kubernetes.io/worker=", worker.Node.HostName)
	exec.Command("sh", "-c", addWorkerLabel).Run()
}
