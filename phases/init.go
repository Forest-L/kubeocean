package phases

import (
	"fmt"
	"github.com/pixiake/kubeocean/bootstrap"
	//"github.com/pixiake/kubeocean/statics"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func CreateCluster(configpath string) error {
	etcdHosts := ssh.Cluster{}
	masterHosts := ssh.Cluster{}
	workerHosts := ssh.Cluster{}
	clusterinfo, err := ssh.GetYamlFile(configpath)
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range clusterinfo {
		fmt.Printf("%v", node)
		for _, role := range node.Roles {
			log.Info("has role: " + role)
			switch role {
			case "master":
				masterHosts.Hosts = append(masterHosts.Hosts, node)
				fmt.Println("I am master")
			case "node":
				workerHosts.Hosts = append(workerHosts.Hosts, node)
				fmt.Println("I am node")
			case "etcd":
				etcdHosts.Hosts = append(etcdHosts.Hosts, node)
				fmt.Println("I am etcd")
			default:
				return fmt.Errorf("Failed to recognize host [%s] role %s", node.Ip, role)
			}
		}
	}

	bootstrap.SwapOff()
	bootstrap.SystemInit()
	bootstrap.Modprobe()
	CreateKubeletService("kubelet", cluster.DefaultKubeVersion, cluster.DefaultKubeImageRepo)

	//err1 := static.RestoreAssets("/usr/bin", "kubeadm")
	//if err1 != nil {
	//	log.Fatal(err)
	//}

	if err := exec.Command("/bin/sh", "-c", "chmod +x /usr/bin/kubeadm").Run(); err != nil {
		log.Fatal("failed to init cluster: %v", err)
	}

	if err := exec.Command("/bin/sh", "-c", "docker pull gcr.io/google-containers/hyperkube:v1.17.0").Run(); err != nil {
		log.Fatal("failed to init cluster: %v", err)
	}

	if err := exec.Command("/usr/bin/kubeadm", "init").Run(); err != nil {
		log.Fatal("failed to init cluster: %v", err)
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
