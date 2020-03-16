package cmd

import (
	"github.com/pixiake/kubeocean/install"
	"github.com/pixiake/kubeocean/scale"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"sync"
)

func NewCmdScaleCluster() *cobra.Command {
	var (
		clusterCfgFile string
	)
	var clusterCmd = &cobra.Command{
		Use:   "scale",
		Short: "Scale cluster",
		Run: func(cmd *cobra.Command, args []string) {
			scaleCluster(clusterCfgFile)
		},
	}

	clusterCmd.Flags().StringVarP(&clusterCfgFile, "cluster-info", "", "", "")
	return clusterCmd
}

func scaleCluster(clusterCfgFile string) {
	cfg, err := cluster.ResolveClusterInfoFile(clusterCfgFile)
	if err != nil {
		log.Fatal(err)
	}
	allNodes, _, masterNodes, _, _ := cfg.GroupHosts()

	clusterStatusInfo, joinMasterCmd, joinWorkerCmd := getClusterStatusInfo(masterNodes)

	log.Info("Install Files Download")
	install.GenerateBootStrapScript(cfg)
	install.InstallFilesDownload(cfg)
	install.GenerateKubeletService()
	install.BootStrapOS(allNodes)
	newNodes := NewNodes(clusterStatusInfo, allNodes)
	install.OverrideHostname(&newNodes)
	install.InstallDocker(&newNodes)
	install.GetKubeBinary(cfg, &newNodes)
	install.SetKubeletService(&newNodes)

	result := make(chan string)
	ccons := make(chan struct{}, ssh.DefaultCon)
	masterNum := len(masterNodes.Hosts) - 1
	wg := &sync.WaitGroup{}
	go ssh.CheckResults(result, masterNum, wg, ccons)

	for _, node := range newNodes.Hosts {
		ccons <- struct{}{}
		wg.Add(1)
		go func(joinMasterCmd, joinWorkerCmd string, node *cluster.ClusterNodeCfg, rs chan string) {
			if node.IsMaster {
				scale.JoinMaster(node, joinMasterCmd)
				if node.IsWorker {
					install.RemoveMasterTaint(node)
				}
			} else {
				if node.IsWorker {
					scale.JoinWorker(node, joinWorkerCmd)
				}
			}
			rs <- "ok"
		}(joinMasterCmd, joinWorkerCmd, &node, result)
	}
	wg.Wait()
}

func getClusterStatusInfo(masters *cluster.MasterNodes) (string, string, string) {
	var clusterStatusInfo string
	var err error
	getInfoCmd := "/usr/local/bin/kubectl get node -o wide"
	for _, master := range masters.Hosts {
		clusterStatusInfo, err = master.CmdExecOut(getInfoCmd)
		if err != nil {
			continue
		} else {
			joinMasterCmd, joinWorkerCmd := scale.GetJoinCmd(&master)
			return clusterStatusInfo, joinMasterCmd, joinWorkerCmd
		}
	}
	return "", "", ""
}

func NewNodes(clusterStatusInfo string, nodes *cluster.AllNodes) cluster.AllNodes {
	newNodes := cluster.AllNodes{}
	for _, node := range nodes.Hosts {
		if strings.Contains(clusterStatusInfo, node.Node.HostName) == false && strings.Contains(clusterStatusInfo, node.Node.InternalAddress) == false && (node.IsMaster || node.IsWorker) {
			newNodes.Hosts = append(newNodes.Hosts, node)
		}
	}
	return newNodes
}
