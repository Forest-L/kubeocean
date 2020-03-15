package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func InstallFilesDownload(cfg *cluster.ClusterCfg) {
	var kubeVersion string
	if cfg.KubeVersion == "" {
		kubeVersion = cluster.DefaultKubeVersion
	} else {
		kubeVersion = cfg.KubeVersion
	}

	kubeadmUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubeadm", kubeVersion, cluster.DefaultArch)
	kubeletUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubelet", kubeVersion, cluster.DefaultArch)
	kubectlUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubectl", kubeVersion, cluster.DefaultArch)
	kubeCniUrl := fmt.Sprintf("https://containernetworking.pek3b.qingstor.com/plugins/releases/download/%s/cni-plugins-linux-%s-%s.tgz", "v0.8.1", cluster.DefaultArch, "v0.8.1")

	kubeadm := fmt.Sprintf("/tmp/kubeocean/kubeadm-%s", kubeVersion)
	kubelet := fmt.Sprintf("/tmp/kubeocean/kubelet-%s", kubeVersion)
	kubectl := fmt.Sprintf("/tmp/kubeocean/kubectl-%s", kubeVersion)
	kubeCni := fmt.Sprintf("/tmp/kubeocean/cni-plugins-linux-%s-%s.tgz", cluster.DefaultArch, "v0.8.1")

	getKubeadmCmd := fmt.Sprintf("curl -o %s  %s", kubeadm, kubeadmUrl)
	getKubeletCmd := fmt.Sprintf("curl -o %s  %s", kubelet, kubeletUrl)
	getKubectlCmd := fmt.Sprintf("curl -o %s  %s", kubectl, kubectlUrl)
	getKubeCniCmd := fmt.Sprintf("curl -o %s  %s", kubeCni, kubeCniUrl)

	log.Info("Kubeadm being download ...")
	if util.IsExist(kubeadm) == false {
		if err := exec.Command("/bin/sh", "-c", getKubeadmCmd).Run(); err != nil {
			log.Errorf("Failed to get kubeadm: %v", err)
		}
	}

	log.Info("Kubelet being download ...")
	if util.IsExist(kubelet) == false {
		if err := exec.Command("/bin/sh", "-c", getKubeletCmd).Run(); err != nil {
			log.Errorf("Failed to get kubelet: %v", err)
		}
	}

	log.Info("Kubectl being download ...")
	if util.IsExist(kubectl) == false {
		if err := exec.Command("/bin/sh", "-c", getKubectlCmd).Run(); err != nil {
			log.Errorf("Failed to get kubectl: %v", err)
		}
	}

	log.Info("KubeCni being download ...")
	if util.IsExist(kubeCni) == false {
		if err := exec.Command("/bin/sh", "-c", getKubeCniCmd).Run(); err != nil {
			log.Errorf("Failed to get kubecni: %v", err)
		}
	}
}

func GenerateBootStrapScript(cfg *cluster.ClusterCfg) {
	tmpl.GenerateBootStrapScript(cfg.GenerateHosts())
}

func GenerateKubeletService() {
	tmpl.GenerateKubeletFiles()
}
