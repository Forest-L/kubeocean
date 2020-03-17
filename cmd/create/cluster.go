package create

import (
	"github.com/pixiake/kubeocean/install"
	"github.com/pixiake/kubeocean/scale"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sync"
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
	allinone := cluster.ClusterNodeCfg{}
	nodeList := []cluster.ClusterNodeCfg{allinone}
	nodes := cluster.AllNodes{nodeList}
	log.Info("Install Files Download")
	install.GenerateBootStrapScript(cfg)
	install.InstallFilesDownload(cfg)
	install.GenerateKubeletService()
	install.BootStrapOS(&nodes)
	install.OverrideHostname(&nodes)
	install.InstallDocker(&nodes)
	install.GetKubeBinary(cfg, &nodes)
	install.SetKubeletService(&nodes)
	log.Info("Init Cluster")
	install.InitCluster(cfg, &nodes.Hosts[0])
	install.RemoveMasterTaint(&nodes.Hosts[0])

}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	allNodes, _, masterNodes, workerNodes, k8sNodes := cfg.GroupHosts()
	log.Info("Install Files Download")
	install.GenerateBootStrapScript(cfg)
	install.InstallFilesDownload(cfg)
	install.GenerateKubeletService()
	install.InstallDocker(allNodes)
	install.OverrideHostname(allNodes)
	install.BootStrapOS(allNodes)
	install.GetKubeBinary(cfg, allNodes)
	install.SetKubeletService(allNodes)

	log.Info("Init Cluster")
	install.InitCluster(cfg, &masterNodes.Hosts[0])

	if len(k8sNodes.Hosts) > 1 {
		log.Infof("Join Masters")
		joinMasterCmd, joinWorkerCmd := scale.GetJoinCmd(&masterNodes.Hosts[0])

		result := make(chan string)
		ccons := make(chan struct{}, ssh.DefaultCon)
		masterNum := len(masterNodes.Hosts) - 1
		wg := &sync.WaitGroup{}
		go ssh.CheckResults(result, masterNum, wg, ccons)

		for index, master := range masterNodes.Hosts {
			if index != 0 {
				host := master
				ccons <- struct{}{}
				wg.Add(1)
				go func(joinMasterCmd string, host *cluster.ClusterNodeCfg, rs chan string) {
					scale.JoinMaster(host, joinMasterCmd)
					if master.IsWorker {
						install.RemoveMasterTaint(host)
					}
					rs <- "ok"
				}(joinMasterCmd, &host, result)
			}
		}
		wg.Wait()

		workerNum := len(workerNodes.Hosts)
		go ssh.CheckResults(result, workerNum, wg, ccons)
		for _, worker := range workerNodes.Hosts {
			host := worker
			ccons <- struct{}{}
			wg.Add(1)
			go func(joinWorkerCmd string, host *cluster.ClusterNodeCfg, rs chan string) {
				if host.IsMaster != true {
					scale.JoinWorker(host, joinWorkerCmd)
				}
				rs <- "ok"
			}(joinWorkerCmd, &host, result)
		}
		wg.Wait()
	}
}
