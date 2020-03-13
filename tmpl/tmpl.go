package tmpl

import (
	"fmt"
	"github.com/pixiake/kubeocean/util"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os"
	"text/template"
)

type KubeContainer struct {
	Repo    string
	Version string
}

type File struct {
	Name string
	Pem  os.FileMode
	Tmpl *template.Template
}

func GenerateBootStrapScript() {
	tmpPath := "/tmp/kubeocean"
	if util.IsExist(tmpPath) == false {
		err := os.MkdirAll(tmpPath, os.ModePerm)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	BootStrapTmpl, e := template.ParseFiles(fmt.Sprintf("%s/tmpl/bootStrapScript.sh", wd))
	if e != nil {
		log.Fatalf("%v", e)
	}
	bootStrapScript := fmt.Sprintf("%s/bootStrapScript.sh", tmpPath)
	file, err := os.OpenFile(bootStrapScript, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0755)
	defer file.Close()
	if err != nil {
		log.Errorf("%v", err)
	}
	err1 := BootStrapTmpl.Execute(file, "KubeOcean")
	if err1 != nil {
		fmt.Println("test")
		log.Errorf("%v", err1)
	}

}

func createDirectory(directory []string) {
	dirs := directory
	for _, v := range dirs {
		if util.IsExist(v) {
			fmt.Printf("%s is exist!", v)
		} else {
			err := os.MkdirAll(v, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func GenerateKubeletFiles() {
	dir := []string{"/etc/systemd/system/kubelet.service.d"}
	createDirectory(dir)
	KubeletServiceTempl, _ := template.ParseFiles("kubelet.service")
	KubeletEnvTempl, _ := template.ParseFiles("10-kubeadm.conf")
	kubeletService := File{Name: "/etc/systemd/system/kubelet.service", Pem: 0644, Tmpl: KubeletServiceTempl}
	kubeletEnv := File{Name: "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf", Pem: 0644, Tmpl: KubeletEnvTempl}

	kubeletFiles := []File{kubeletService, kubeletEnv}
	for _, f := range kubeletFiles {
		file, err := os.OpenFile(f.Name, os.O_CREATE|os.O_WRONLY|os.O_SYNC, f.Pem)
		defer file.Close()
		if err != nil {
			log.Errorf("%v", err)
		}
		kube := KubeContainer{}
		f.Tmpl.Execute(file, kube)
	}
}

func GenerateKubeadmFiles(cfg *cluster.ClusterCfg) {
	dir := []string{"/etc/kubernetes"}
	createDirectory(dir)
	kubeadmCfg := cfg.GenerateKubeadmCfg()

	KubeadmCfgTempl, _ := template.ParseFiles("kubeadm-config.yaml")
	kubeadmCfgFile := File{Name: "/tmp/kubeocean/kubeadm-config.yaml", Pem: 0644, Tmpl: KubeadmCfgTempl}

	kubeFiles := []File{kubeadmCfgFile}
	for _, f := range kubeFiles {
		file, err := os.OpenFile(f.Name, os.O_CREATE|os.O_WRONLY|os.O_SYNC, f.Pem)
		defer file.Close()
		if err != nil {
			log.Errorf("%v", err)
		}
		f.Tmpl.Execute(file, kubeadmCfg)
	}
}

func GenerateNetworkPluginFiles(cfg *cluster.ClusterCfg) {
	var fileName string
	var tmpl *template.Template
	if cfg.Network.Plugin == "calico" || cfg.Network.Plugin == "" {
		tmpl, _ = template.ParseFiles("calico.yaml")
		fileName = "/tmp/kubeocean/calico.yaml"
	}

	if cfg.Network.Plugin == "flannel" {
		tmpl, _ = template.ParseFiles("flannel.yaml")
		fileName = "/tmp/kubeocean/flannel.yaml"
	}

	networkPluginFile := File{Name: fileName, Pem: 0644, Tmpl: tmpl}

	file, err := os.OpenFile(networkPluginFile.Name, os.O_CREATE|os.O_WRONLY|os.O_SYNC, networkPluginFile.Pem)
	defer file.Close()
	if err != nil {
		log.Errorf("%v", err)
	}
	networkPluginFile.Tmpl.Execute(file, cfg.Network)

}
