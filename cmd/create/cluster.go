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
	install.InitOS(&nodes)
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
	log.Info("Inject Hosts")
	install.InjectHosts(cfg, &nodes)
	log.Info("Init Cluster")
	install.InitCluster(cfg, &masters)
	install.RemoveMasterTaint(&masters)

}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	allNodes, _, masterNodes, workerNodes := cfg.GroupHosts()
	log.Info("BootStrap")
	install.InitOS(allNodes)
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
	log.Info("Inject Hosts")
	install.InjectHosts(cfg, allNodes)
	log.Info("Init Cluster")
	install.InitCluster(cfg, masterNodes)
	install.RemoveMasterTaint(masterNodes)
	if len(masterNodes.Hosts) > 1 {
		scale.JoinMasters(masterNodes)
	}
	log.Info("Join Workers")
	scale.JoinWorkers(workerNodes, masterNodes)
}
