package cmd

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/ssh-bak"
	"github.com/spf13/cobra"
	"strings"
)

func NewCmdExec() *cobra.Command {

	exec := ssh_bak.ExecInfo{}
	var execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Batch SSH commands",
		Long:  "A simple parallel SSH tool that allows you to execute command combinations to cluster by SSH.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(strings.Join(args, ","))
			ssh_bak.ExecuteCmd(exec)
		},
	}

	execCmd.Flags().StringVarP(&exec.Hosts, "cluster", "", "", "host address list")
	execCmd.Flags().StringVarP(&exec.Ips, "ips", "", "", "ip address list")
	execCmd.Flags().StringVarP(&exec.Cmds, "cmds", "m", "", "cmds")
	execCmd.Flags().StringVarP(&exec.Username, "username", "u", "root", "username")
	execCmd.Flags().StringVarP(&exec.Password, "password", "p", "", "password")
	execCmd.Flags().StringVarP(&exec.Key, "key", "k", "", "ssh-bak private key")
	execCmd.Flags().IntVarP(&exec.Port, "port", "", 22, "ssh-bak port")
	execCmd.Flags().StringVarP(&exec.Ciphers, "ciphers", "", "", "ciphers")
	execCmd.Flags().StringVarP(&exec.CmdFile, "cmdfile", "", "", "cmdfile path")
	execCmd.Flags().StringVarP(&exec.HostFile, "hostfile", "", "", "hostfile path")
	execCmd.Flags().StringVarP(&exec.IpFile, "ipfile", "", "", "ipfile path")
	execCmd.Flags().StringVarP(&exec.CfgFile, "config", "", "", "config File Path")
	execCmd.Flags().BoolVarP(&exec.JsonMode, "jsonMode", "j", false, "print output in json format")
	execCmd.Flags().BoolVarP(&exec.OutTxt, "outTxt", "", false, "write result into txt")
	execCmd.Flags().StringVarP(&exec.FileLocate, "fileLocate", "", "", "write file locate")
	execCmd.Flags().BoolVarP(&exec.LinuxMode, "linuxMode", "", true, "In linux mode,multi command combine with && ,such as date&&cd /opt&&ls")
	execCmd.Flags().IntVarP(&exec.TimeLimit, "timeLimit", "t", 600, "max timeout")
	execCmd.Flags().IntVarP(&exec.NumLimit, "numLimit", "n", 20, "max execute number")

	return execCmd
}
