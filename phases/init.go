package phases

import (
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func CreateCluster(configpath string) error {

	CreateKubeletService("kubelet", cluster.DefaultKubeVersion, cluster.DefaultKubeImageRepo)

	//err1 := static.RestoreAssets("/usr/bin", "kubeadm")
	//if err1 != nil {
	//	log.Fatal(err)
	//}

	if err := exec.Command("/bin/sh", "-c", "chmod +x /usr/bin/kubeadm").Run(); err != nil {
		log.Fatal("failed to initsystem cluster: %v", err)
	}

	if err := exec.Command("/bin/sh", "-c", "docker pull gcr.io/google-containers/hyperkube:v1.17.0").Run(); err != nil {
		log.Fatal("failed to initsystem cluster: %v", err)
	}

	if err := exec.Command("/usr/bin/kubeadm", "initsystem").Run(); err != nil {
		log.Fatal("failed to initsystem cluster: %v", err)
	}

	cmdconfig := "mkdir -p /root/.kube && cp /etc/kubernetes/admin.conf /root/.kube/config"
	if err := exec.Command("/bin/sh", "-c", cmdconfig).Run(); err != nil {
		log.Fatal("failed to create config: %v", err)
	}
	cmdcalico := "kubectl apply -f https://docs.projectcalico.org/v3.8/manifests/calico.yaml"
	if err := exec.Command("/bin/sh", "-c", cmdcalico).Run(); err != nil {
		log.Fatal("failed to create config: %v", err)
	}

	cmdkubectl := "docker cp kubelet:/usr/local/bin/kubectl /usr/bin/kubectl"
	if err := exec.Command("/bin/sh", "-c", cmdkubectl).Run(); err != nil {
		log.Fatal("failed to create config: %v", err)
	}
	return nil
}
