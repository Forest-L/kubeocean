package cmd

import (
	"fmt"
	"github.com/pixiake/kubeocean/util"
	"github.com/spf13/cobra"
)

func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "get version information",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(util.VERSION)
		},
	}

	return versionCmd
}
