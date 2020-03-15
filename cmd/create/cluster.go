package create

import (
	"github.com/pixiake/kubeocean/install"
	"github.com/pixiake/kubeocean/scale"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdCreateCluster() *cobra.Command {
	var (
		clusterCfgFile string
		kubeadmCfgFile string
	)
	var clusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "Create cluster",
		Run: func(cmd *cobra.Command, args []string) {
			createCluster(clusterCfgFile)
		},
	}

	clusterCmd.Flags().StringVarP(&clusterCfgFile, "cluster-info", "", "", "")
	clusterCmd.Flags().StringVarP(&kubeadmCfgFile, "kubeadm-config", "", "", "")
	return clusterCmd
}

func createCluster(clusterCfgFile string) {
	if clusterCfgFile != "" {
		clusterInfo, err := cluster.ResolveClusterInfoFile(clusterCfgFile)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(clusterInfo)
		log.Info("Welcome to KubeOcean")
		createMultiNodes(clusterInfo)
	} else {
		log.Info("Welcome to KubeOcean")
		clusterInfo := cluster.ClusterCfg{}
		createAllinone(&clusterInfo)
	}
}

func createAllinone(cfg *cluster.ClusterCfg) {
	nodes := cluster.AllNodes{}
	log.Info("Install Files Download")
	install.InstallFilesDownload(cfg)
	install.GenerateBootStrapScript(cfg)
	install.BootStrapOS(&nodes.Hosts[0])
	install.OverrideHostname(&nodes.Hosts[0])
	install.InstallDocker(&nodes.Hosts[0])
	install.GetKubeBinary(cfg, &nodes.Hosts[0])
	install.SetKubeletService(&nodes.Hosts[0])
	log.Info("Init Cluster")
	install.InitCluster(cfg, &nodes.Hosts[0])
	install.RemoveMasterTaint(&nodes.Hosts[0])

}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	allNodes, _, masterNodes, workerNodes, k8sNodes := cfg.GroupHosts()
	log.Info("Install Files Download")
	install.InstallFilesDownload(cfg)
	install.GenerateBootStrapScript(cfg)
	for _, node := range allNodes.Hosts {
		install.BootStrapOS(&node)
		install.OverrideHostname(&node)
		install.InstallDocker(&node)
		install.GetKubeBinary(cfg, &node)
		install.SetKubeletService(&node)
	}

	log.Info("Init Cluster")
	install.InitCluster(cfg, &masterNodes.Hosts[0])

	if len(k8sNodes.Hosts) > 1 {
		joinMasterCmd, joinWorkerCmd := scale.GetJoinCmd(&masterNodes.Hosts[0])
		for index, master := range masterNodes.Hosts {
			if index != 0 {
				scale.JoinMaster(&master, joinMasterCmd)
				install.RemoveMasterTaint(&master)
			}
		}
		for _, worker := range workerNodes.Hosts {
			if worker.IsMaster != true {
				scale.JoinWorker(&worker, joinWorkerCmd)
			}
		}
	}
}
