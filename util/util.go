package util

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"os"
)

const (
	VERSION = "KubeOcean Version v0.0.1\nKubernetes Version v1.17.0"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func CmdExec(isLocal bool, cmd string, host cluster.NodeCfg) {
	if isLocal {

	}
}
