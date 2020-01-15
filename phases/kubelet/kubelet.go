package kubelet

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

type kubeletContainer struct {
	KubeRepo    string
	KubeVersion string
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
func createDirectory(directory []string) {
	dirs := directory
	for _, v := range dirs {
		if isExist(v) {
			fmt.Printf("%s is exist!", v)
		} else {
			err := os.MkdirAll(v, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func generateFile(name string, perm os.FileMode, tmpl *template.Template, container kubeletContainer) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_SYNC, perm)
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(file, container)
	file.Close()
}

func CreateKubeletService(service string, version string, repo string) {
	if len(service) == 0 {
		fmt.Println("Please provide the service name !")
	} else if service == "kubelet" {
		dir := []string{"/etc/systemd/system/kubelet.service.d"}
		createDirectory(dir)
		kubeletContainertest := kubeletContainer{repo, version}
		generateFile("/usr/bin/kubelet", 0755, tmpl.GetTmpl("kubelet"), kubeletContainertest)
		generateFile("/etc/systemd/system/kubelet.service.d/kubelet-contain.conf", 0644, tmpl.GetTmpl("kubeletContainer"), kubeletContainertest)
		generateFile("/etc/systemd/system/kubelet.service", 0644, tmpl.GetTmpl("kubeletService"), kubeletContainertest)

	} else {
		fmt.Println("nonsupport this service!")
	}
}

func commandCreate() {
	var version string
	var registry string

	var cmdStartKubelet = &cobra.Command{
		Use:   "start [service name]",
		Short: "start the service",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("the service is running: %s", args[0])
			CreateKubeletService(args[0], version, registry)
		},
	}

	cmdStartKubelet.Flags().StringVarP(&registry, "registry", "r", "gcr.io", "kubernetes containers repo url")
	cmdStartKubelet.Flags().StringVarP(&version, "version", "v", "v1.17.0", "kubernetes version")

	var rootCmd = &cobra.Command{Use: "kubeocean"}
	rootCmd.AddCommand(cmdStartKubelet)
	rootCmd.Execute()

}
