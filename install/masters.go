package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func InitCluster(masters *cluster.MasterNodes) {
	if masters == nil {
		if out, err := exec.Command("/usr/local/bin/kubeadm", "init").CombinedOutput(); err != nil {
			log.Fatalf("failed to init cluster:\n %v", out)
		}
	}
	for index, master := range masters.Hosts {
		fmt.Println(index)
		fmt.Println(master)
	}
}
