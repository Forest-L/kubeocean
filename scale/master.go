package scale

import (
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
)

func JoinMasters(master *cluster.MasterNodes) {
	joinMasterCmd := JoinMasterCmd(master)
	for index, master := range master.Hosts {
		if index != 0 {
			if err := master.CmdExec(joinMasterCmd); err != nil {
				log.Fatalf("Failed to add master (%s):\n", master.Node.Address)
			}
		}
	}
}
