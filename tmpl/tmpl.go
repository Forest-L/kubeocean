package tmpl

import (
	"fmt"
	"github.com/pixiake/kubeocean/util"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os"
)

type KubeContainer struct {
	Repo    string
	Version string
}

type File struct {
	Name string
	pem  os.FileMode
}

func GenerateBootStrapScript() {
	tmpPath := "/tmp/kubeocean"
	if util.IsExist(tmpPath) == false {
		err := os.MkdirAll(tmpPath, os.ModePerm)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
	bootStrapScript := fmt.Sprintf("%s/bootStrapScript.sh")
	file, err := os.OpenFile(bootStrapScript, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0755)
	defer file.Close()
	if err != nil {
		log.Errorf("%v", err)
	}
	BootStrapTmpl.Execute(file, nil)

}

func GenerateKubeletFiles(cfg cluster.ClusterCfg) {
	kubeContainerInfo := KubeContainer{}
	kubeContainerInfo.Repo = cfg.KubeImageRepo
	kubeContainerInfo.Version = cfg.KubeVersion

	kubelet := File{Name: "/usr/bin/kubelet", pem: 0755}
	kubeletContainer := File{Name: "/etc/systemd/system/kubelet.service.d/kubelet-contain.conf", pem: 0644}
	kubeletService := File{Name: "/etc/systemd/system/kubelet.service", pem: 0644}
	kubeletFiles := []File{kubelet, kubeletContainer, kubeletService}
	for _, f := range kubeletFiles {
		file, err := os.OpenFile(f.Name, os.O_CREATE|os.O_WRONLY|os.O_SYNC, f.pem)
		defer file.Close()
		if err != nil {
			log.Errorf("%v", err)
		}
		BootStrapTmpl.Execute(file, kubeContainerInfo)
	}
}
