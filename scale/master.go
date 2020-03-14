package scale

import (
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
)

func JoinMaster(master *cluster.ClusterNodeCfg, joinMasterCmd string) {
	if err := master.CmdExec(joinMasterCmd); err != nil {
		log.Fatalf("Failed to add master (%s):\n", master.Node.Address)
	}
}
