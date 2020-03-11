package scale

import (
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
)

func JoinMasters(master *cluster.MasterNodes) {
	joinMasterCmd := JoinMasterCmd(master)
	for index, master := range master.Hosts {
		if index != 0 {
			if err := ssh.CmdExec(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, joinMasterCmd); err != nil {
				log.Fatalf("Failed to add master (%s):\n", master.Node.Address)
			}
		}
	}
}
