package create

import "github.com/spf13/cobra"

func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Kubeocean create config and cluster",
	}
	cmd.AddCommand(NewCmdCreateCfg())

	return cmd
}
