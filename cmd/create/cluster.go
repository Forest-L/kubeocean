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
	masters := cluster.MasterNodes{}
	log.Info("BootStrap")
	install.InitOS(cfg, &nodes)
	log.Info("Override Hostname")
	install.OverrideHostname(&nodes)
	log.Info("Install Files Download")
	install.InstallFilesDownload(cfg.KubeVersion)
	log.Info("Install Docker")
	install.DockerInstall(&nodes)
	log.Info("Get Kube Binary")
	install.GetKubeBinary(cfg, &nodes)
	log.Info("Set Kubelet Service")
	install.SetKubeletService(&nodes)
	log.Info("Init Cluster")
	install.InitCluster(cfg, &masters)
	install.RemoveMastersTaint(&masters)

}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	allNodes, _, masterNodes, workerNodes, k8sNodes := cfg.GroupHosts()
	log.Info("BootStrap")
	install.InitOS(cfg, allNodes)
	log.Info("Override Hostname")
	install.OverrideHostname(allNodes)
	log.Info("Install Files Download")
	install.InstallFilesDownload(cfg.KubeVersion)
	log.Info("Install Docker")
	install.DockerInstall(allNodes)
	log.Info("Get Kube Binary")
	install.GetKubeBinary(cfg, allNodes)
	log.Info("Set Kubelet Service")
	install.SetKubeletService(allNodes)
	log.Info("Init Cluster")
	install.InitCluster(cfg, masterNodes)

	if len(k8sNodes.Hosts) > 1 {
		joinMasterCmd, joinWorkerCmd := scale.GetJoinCmd(&masterNodes.Hosts[0])
		for index, master := range masterNodes.Hosts {
			if index != 0 {
				scale.JoinMaster(&master, joinMasterCmd)
			}
		}
		install.RemoveMastersTaint(masterNodes)
		for _, worker := range workerNodes.Hosts {
			if worker.IsMaster != true {
				scale.JoinWorker(&worker, joinWorkerCmd)
			}
		}
	}
}
