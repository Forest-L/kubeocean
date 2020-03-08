package cmd

import (
	"github.com/pixiake/kubeocean/cmd/create"
	"github.com/spf13/cobra"
)

func NewKubeoceanCommand() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "ko",
		Short: "Batch SSH commands",
		Long:  "A simple parallel SSH tool that allows you to execute command combinations to cluster by SSH.",
	}
	rootCmd.AddCommand(create.NewCmdCreate())
	rootCmd.AddCommand(NewCmdExec())
	rootCmd.AddCommand(NewCmdVersion())
	return rootCmd
}
