package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
)

func AddWorkers(workers *cluster.WorkerNodes) {
	fmt.Println(workers)
}
