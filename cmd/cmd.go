package cmd

import (
	"github.com/spf13/cobra"
)

func NewKubeoceanCommand() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "kubeocean",
		Short: "Batch SSH commands",
		Long:  "A simple parallel SSH tool that allows you to execute command combinations to hosts by SSH.",
	}

	rootCmd.AddCommand(NewCmdExec())
	rootCmd.AddCommand(NewCmdVersion())
	return rootCmd
}
