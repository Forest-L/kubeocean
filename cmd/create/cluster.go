package create

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
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
		dir, _ := os.Executable()
		exPath := filepath.Dir(dir)
		configFile := fmt.Sprintf("%s/%s", exPath, "cluster-info.yaml")
		clusterInfo, err := cluster.ResolveClusterInfoFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", clusterInfo)
	} else {
		fmt.Printf("Init a allinone cluster")
	}

}
