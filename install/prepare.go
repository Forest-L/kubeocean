package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func InstallFilesDownload(kubeVersion string) {
	if kubeVersion == "" {
		kubeVersion = cluster.DefaultKubeVersion
	}
	kubeadmUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubeadm", kubeVersion, cluster.DefaultArch)
	kubeletUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubelet", kubeVersion, cluster.DefaultArch)
	kubectlUrl := fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubectl", kubeVersion, cluster.DefaultArch)
	kubeCniUrl := fmt.Sprintf("https://containernetworking.pek3b.qingstor.com/plugins/releases/download/%s/cni-plugins-linux-%s-%s.tgz", "v0.8.1", cluster.DefaultArch, "v0.8.1")

	getKubeadmCmd := fmt.Sprintf("curl -o /tmp/kubeocean/kubeadm-%s  %s", kubeVersion, kubeadmUrl)
	getKubeletCmd := fmt.Sprintf("curl -o /tmp/kubeocean/kubelet-%s  %s", kubeVersion, kubeletUrl)
	getKubectlCmd := fmt.Sprintf("curl -o /tmp/kubeocean/kubectl-%s  %s", kubeVersion, kubectlUrl)
	getKubeCniCmd := fmt.Sprintf("curl -o /tmp/kubeocean/cni-plugins-linux-amd64-v0.8.1.tgz  %s", kubeCniUrl)

	if err := exec.Command("sudo", getKubeadmCmd).Run(); err != nil {
		log.Errorf("Failed to get kubeadm: %v", err)
	}

	if err := exec.Command("sudo", getKubeletCmd).Run(); err != nil {
		log.Errorf("Failed to get kubelet: %v", err)
	}

	if err := exec.Command("sudo", getKubectlCmd).Run(); err != nil {
		log.Errorf("Failed to get kubectl: %v", err)
	}

	if err := exec.Command("sudo", getKubeCniCmd).Run(); err != nil {
		log.Errorf("Failed to get kubecni: %v", err)
	}
}
