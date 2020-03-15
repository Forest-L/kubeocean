package cmd

import (
	"github.com/pixiake/kubeocean/install"
	"github.com/pixiake/kubeocean/scale"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
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
	install.InstallFilesDownload(cfg)
	install.GenerateBootStrapScript(cfg)
	install.GenerateKubeletService()

	for _, node := range NewNodes(clusterStatusInfo, allNodes) {
		install.BootStrapOS(&node)
		install.OverrideHostname(&node)
		install.InstallDocker(&node)
		install.GetKubeBinary(cfg, &node)
		install.SetKubeletService(&node)
		if node.IsMaster {
			scale.JoinMaster(&node, joinMasterCmd)
			if node.IsWorker {
				install.RemoveMasterTaint(&node)
			}
		} else {
			if node.IsWorker {
				scale.JoinWorker(&node, joinWorkerCmd)
			}
		}
	}
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

func NewNodes(clusterStatusInfo string, nodes *cluster.AllNodes) []cluster.ClusterNodeCfg {
	newNodes := []cluster.ClusterNodeCfg{}
	for _, node := range nodes.Hosts {
		if strings.Contains(clusterStatusInfo, node.Node.HostName) == false && strings.Contains(clusterStatusInfo, node.Node.InternalAddress) == false && (node.IsMaster || node.IsWorker) {
			newNodes = append(newNodes, node)
		}
	}
	return newNodes
}
