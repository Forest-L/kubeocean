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
	log.Info("BootStrap")
	install.InitOS(nil)
	install.OverrideHostname(nil)
	install.InstallFilesDownload(cfg.KubeVersion)
	install.DockerInstall(nil)
	install.GetKubeBinary(cfg, nil)
	install.SetKubeletService(nil)
	install.InjectHosts(cfg, nil)
	install.InitCluster(cfg, nil)

}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	allNodes, _, masterNodes, workerNodes := cfg.GroupHosts()
	install.InitOS(allNodes)
	install.OverrideHostname(allNodes)
	install.InstallFilesDownload(cfg.KubeVersion)
	install.DockerInstall(allNodes)
	install.GetKubeBinary(cfg, allNodes)
	install.SetKubeletService(allNodes)
	install.InjectHosts(cfg, allNodes)
	//install.SetUpEtcd(etcdNodes)
	install.InitCluster(cfg, masterNodes)
	install.RemoveMasterTaint(masterNodes)
	if len(masterNodes.Hosts) > 1 {
		scale.JoinMasters(masterNodes)
	}
	scale.JoinWorkers(workerNodes, masterNodes)
}
