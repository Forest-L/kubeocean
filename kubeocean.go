package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/kubelet"
	"os"
	"text/template"
	"github.com/pixiake/kubeocean/tmpl/tmpl"
)



type kubeletContainer struct {
	KubeRepo string
	KubeVersion string
}
func main()  {
	//kubeletContainertest := kubeletContainer{"gcr.io/google_container", "v1.20.3"}
	//generateFile("/root/kubelet", 0755, kubeletContainerTempl, kubeletContainertest)
	//generateFile("kubelet.service", 0755, kubeletServiceTempl, kubeletContainertest)
	//checkPort("10250")
	commandCreate()
}

func generateFile(name string, perm os.FileMode, tmpl *template.Template, container kubeletContainer)  {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_SYNC, perm)
	if err != nil {
		fmt.Println(err)
	}

	tmpl.Execute(file, container)
	file.Close()
}

func startService(service string, version string, repo string) {
	if len(service) == 0 {
		fmt.Println("Please provide the service name !")
	} else if service == "kubelet" {
		kubeletContainertest := kubeletContainer{repo, version}
		generateFile("/usr/local/bin/kubelet", 0755, tmpl.kubeletContainerTempl, kubeletContainertest)
		generateFile("/etc/systemd/system/kubelet.service", 0644, tmpl.kubeletServiceTempl, kubeletContainertest)
		//kubelet.TryStartKubelet()

		fmt.Println(kubeletContainertest.KubeRepo)
		fmt.Println(kubeletContainertest.KubeVersion)

		//generateFile("kubelet", 0755, kubeletContainerTempl, kubeletContainertest)
		//generateFile("kubelet.service", 0644, kubeletServiceTempl, kubeletContainertest)
		kubelet.TryStartKubelet()

	} else {
		fmt.Println("nonsupport this service!")
	}
}


func commandCreate()  {
	var version string
	var registry string

	var cmdStartKubelet = &cobra.Command{
		Use:                        "start [service name]",
		Short:                      "start the service",
		Args:                       cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("the service is running: %s", args[0])
			startService(args[0], version, registry)
		},
	}

	cmdStartKubelet.Flags().StringVarP(&registry, "registry", "r", "gcr.io/google_containers", "kubernetes containers repo url")
	cmdStartKubelet.Flags().StringVarP(&version, "version", "v", "v1.17.0", "kubernetes version")


	var rootCmd = &cobra.Command{Use: "kubeocean"}
	rootCmd.AddCommand(cmdStartKubelet)
	rootCmd.Execute()
}
