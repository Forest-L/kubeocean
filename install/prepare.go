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
	kubeCniUrl := fmt.Sprintf("https://containernetworking.pek3b.qingstor.com/plugins/releases/download/%s/cni-plugins-linux-%s-%s.tgz", cluster.DefaultCniVersion, cluster.DefaultArch, cluster.DefaultCniVersion)
	HelmUrl := fmt.Sprintf("https://kubernetes-helm.pek3b.qingstor.com/linux-amd64/%s/helm", cluster.DefaultHelmVersion)

	kubeadm := fmt.Sprintf("/tmp/kubeocean/kubeadm-%s", kubeVersion)
	kubelet := fmt.Sprintf("/tmp/kubeocean/kubelet-%s", kubeVersion)
	kubectl := fmt.Sprintf("/tmp/kubeocean/kubectl-%s", kubeVersion)
	kubeCni := fmt.Sprintf("/tmp/kubeocean/cni-plugins-linux-%s-%s.tgz", cluster.DefaultArch, cluster.DefaultCniVersion)
	helm := fmt.Sprintf("/tmp/kubeocean/helm-%s", cluster.DefaultHelmVersion)

	getKubeadmCmd := fmt.Sprintf("curl -o %s  %s", kubeadm, kubeadmUrl)
	getKubeletCmd := fmt.Sprintf("curl -o %s  %s", kubelet, kubeletUrl)
	getKubectlCmd := fmt.Sprintf("curl -o %s  %s", kubectl, kubectlUrl)
	getKubeCniCmd := fmt.Sprintf("curl -o %s  %s", kubeCni, kubeCniUrl)
	getHelmCmd := fmt.Sprintf("curl -o %s  %s", helm, HelmUrl)

	log.Info("Kubeadm being download ...")
	if util.IsExist(kubeadm) == false {
		if out, err := exec.Command("/bin/sh", "-c", getKubeadmCmd).CombinedOutput(); err != nil {
			log.Errorf("Failed to get kubeadm: %v", err)
			fmt.Println(string(out))
		}
	}

	log.Info("Kubelet being download ...")
	if util.IsExist(kubelet) == false {
		if out, err := exec.Command("/bin/sh", "-c", getKubeletCmd).CombinedOutput(); err != nil {
			log.Errorf("Failed to get kubelet: %v", err)
			fmt.Println(string(out))
		}
	}

	log.Info("Kubectl being download ...")
	if util.IsExist(kubectl) == false {
		if out, err := exec.Command("/bin/sh", "-c", getKubectlCmd).CombinedOutput(); err != nil {
			log.Errorf("Failed to get kubectl: %v", err)
			fmt.Println(string(out))
		}
	}

	log.Info("KubeCni being download ...")
	if util.IsExist(kubeCni) == false {
		if out, err := exec.Command("/bin/sh", "-c", getKubeCniCmd).CombinedOutput(); err != nil {
			log.Errorf("Failed to get kubecni: %v", err)
			fmt.Println(string(out))
		}
	}

	log.Info("Helm being download ...")
	if util.IsExist(kubeCni) == false {
		if out, err := exec.Command("/bin/sh", "-c", getHelmCmd).CombinedOutput(); err != nil {
			log.Errorf("Failed to get helm: %v", err)
			fmt.Println(string(out))
		}
	}
}

func GenerateBootStrapScript(cfg *cluster.ClusterCfg) {
	tmpl.GenerateBootStrapScript(cfg.GenerateHosts())
}

func GenerateKubeletService() {
	tmpl.GenerateKubeletFiles()
}
