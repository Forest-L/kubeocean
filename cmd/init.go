package cmd

import (
	"github.com/pixiake/kubeocean/phases"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdInit() *cobra.Command {
	var configpath string
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create a kubernetes cluster",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := phases.CreateCluster(configpath)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	initCmd.Flags().StringVarP(&configpath, "configpath", "", "", "Cluster configuration file path")
	return initCmd
}
