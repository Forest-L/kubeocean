package create

import (
	"fmt"
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
			createCluster(clusterCfgFile, kubeadmCfgFile)
		},
	}

	clusterCmd.Flags().StringVarP(&clusterCfgFile, "cluster-info", "", "", "")
	clusterCmd.Flags().StringVarP(&kubeadmCfgFile, "kubeadm-config", "", "", "")
	return clusterCmd
}

func createCluster(clusterCfgFile string, kubeadmCfgFile string) {

	if clusterCfgFile != "" {
		//dir, _ := os.Executable()
		//exPath := filepath.Dir(dir)
		//configFile := fmt.Sprintf("%s/%s", exPath, "cluster-info.yaml")
		clusterInfo, err := cluster.ResolveClusterInfoFile(clusterCfgFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(clusterInfo)
		createMultiNodes(clusterInfo)
	} else {
		log.Info("Init a allinone cluster")
		createAllinone()
	}

}

func createAllinone() {
	fmt.Printf("test")
}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	for host, _ := range cfg.Hosts {
		fmt.Printf("%v\n", host)
	}
}
