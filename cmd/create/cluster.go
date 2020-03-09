package create

import (
	"fmt"
	"github.com/pixiake/kubeocean/tmpl"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
)

func NewCmdCreateCluster() *cobra.Command {
	var (
		clusterCfgFile string
		kubeadmCfgFile string
	)
	var clusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "Create cluster",
		Run: func(cmd *cobra.Command, args []string) {
			createCluster(clusterCfgFile, kubeadmCfgFile)
		},
	}

	clusterCmd.Flags().StringVarP(&clusterCfgFile, "cluster-info", "", "", "")
	clusterCmd.Flags().StringVarP(&kubeadmCfgFile, "kubeadm-config", "", "", "")
	return clusterCmd
}

func createCluster(clusterCfgFile string, kubeadmCfgFile string) {

	if clusterCfgFile != "" {
		//dir, _ := os.Executable()
		//exPath := filepath.Dir(dir)
		//configFile := fmt.Sprintf("%s/%s", exPath, "cluster-info.yaml")
		clusterInfo, err := cluster.ResolveClusterInfoFile(clusterCfgFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(clusterInfo)
		createMultiNodes(clusterInfo)
	} else {
		log.Info("Init a allinone cluster")
		createAllinone()
	}

}

func createAllinone() {
	clusterCfg := cluster.ClusterCfg{}
	clusterCfg.KubeImageRepo = cluster.DefaultKubeImageRepo
	clusterCfg.KubeVersion = cluster.DefaultKubeVersion
	tmpl.GenerateBootStrapScript()
	if err := exec.Command("/bin/bash", "/tmp/kubeocean/bootStrapScript.sh").Run(); err != nil {
		log.Errorf("bootstrap is failed: %v", err)
	}
	hyperKubeImagePull := fmt.Sprintf("docker pull %s/google-containers/hyperkube:%s", cluster.DefaultKubeImageRepo, cluster.DefaultKubeVersion)
	if err := exec.Command("/bin/bash", hyperKubeImagePull).Run(); err != nil {
		log.Errorf("hyperkube image pull failed: %v", err)
	}

	kubeadmDownLoad := fmt.Sprintf("curl -o /usr/local/bin/kubeadm https://kubernetes-release.pek3b.qingstor.com/release/%s//bin/linux/amd64/kubeadm", cluster.DefaultKubeVersion)
	if err := exec.Command("/bin/bash", kubeadmDownLoad).Run(); err != nil {
		log.Fatal("failed to init cluster: %v", err)
	}

	if err := exec.Command("/bin/sh", "-c", "chmod +x /usr/local/bin/kubeadm").Run(); err != nil {
		log.Fatal("failed to init cluster: %v", err)
	}
	tmpl.GenerateKubeletFiles(clusterCfg)
	if err := exec.Command("/usr/local/bin/kubeadm", "init").Run(); err != nil {
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

	cmdkubectl := "docker cp kubelet:/usr/local/bin/kubectl /usr/local/bin/kubectl"
	if err := exec.Command("/bin/sh", "-c", cmdkubectl).Run(); err != nil {
		log.Fatal("failed to create config: %v", err)
	}

}

func createMultiNodes(cfg *cluster.ClusterCfg) {
	hosts := cfg.Hosts
	etcdNodes := cluster.EtcdNodes{}
	masterNodes := cluster.MasterNodes{}
	workerNodes := cluster.WorkerNodes{}
	for _, host := range hosts {
		for _, role := range host.Role {
			if role == "etcd" {
				etcdNodes.Hosts = append(etcdNodes.Hosts, host)
			}
			if role == "master" {
				masterNodes.Hosts = append(masterNodes.Hosts, host)
			}
			if role == "worker" {
				workerNodes.Hosts = append(workerNodes.Hosts, host)
			}
		}
	}

}
