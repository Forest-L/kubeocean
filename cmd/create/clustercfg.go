package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewCmdCreateCfg() *cobra.Command {
	var culsterCfgCmd = &cobra.Command{
		Use:   "config",
		Short: "Create cluster info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("test")
			//ssh.ExecuteCmd(exec)
		},
	}
	return culsterCfgCmd
}

func clusterCfg() {

}
