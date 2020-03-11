package cluster

import (
	"fmt"
	"github.com/pixiake/kubeocean/util"
	"log"
	"os"
	"strings"
)

const (
	DefaultSSHPort        = "22"
	DefaultDockerSockPath = "/var/run/docker.sock"
	DefaultLBPort         = "6443"
	DefaultLBDomain       = "lb.kubesphere.local"
	DefaultNetworkPlugin  = "calico"
	DefaultPodsCIDR       = "10.233.64.0/18"
	DefaultServiceCIDR    = "10.233.0.0/18"
	DefaultKubeVersion    = "v1.17.3"
	DefaultKubeImageRepo  = "gcr.azk8s.cn"
	DefaultClusterName    = "cluster.local"
	ETCDRole              = "etcd"
	MasterRole            = "master"
	WorkerRole            = "worker"
)

type ClusterCfg struct {
	Hosts []NodeCfg `yaml:"hosts" json:"hosts,omitempty"`
	//SSHKeyPath     string    `yaml:"ssh_key_path" json:"sshKeyPath,omitempty" norman:"nocreate,noupdate"`
	LBKubeApiserver LBKubeApiserverCfg `yaml:"lb_kubeapiserver" json:"lb_kubeapiserver,omitempty"`
	KubeVersion     string             `yaml:"kube_version" json:"kube_version,omitempty"`
	KubeImageRepo   string             `yaml:"kube_image_repo" json:"kube_image_repo,omitempty"`
	Network         NetworkConfig      `yaml:"network" json:"network,omitempty"`
}

type EtcdNodes struct {
	Hosts []ClusterNodeCfg
}

type MasterNodes struct {
	Hosts []ClusterNodeCfg
}

type WorkerNodes struct {
	Hosts []ClusterNodeCfg
}
type ClusterNodeCfg struct {
	Node     NodeCfg
	IsEtcd   bool
	IsMaster bool
	IsWorker bool
}
type NodeCfg struct {
	HostName        string   `yaml:"hostName,omitempty" json:"hostName,omitempty"`
	Address         string   `yaml:"address" json:"address,omitempty"`
	Port            string   `yaml:"port" json:"port,omitempty"`
	InternalAddress string   `yaml:"internal_address" json:"internalAddress,omitempty"`
	Role            []string `yaml:"role" json:"role,omitempty" norman:"type=array[enum],options=etcd|worker|worker"`
	//HostnameOverride string   `yaml:"hostname_override" json:"hostnameOverride,omitempty"`
	User     string `yaml:"user" json:"user,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
	//SSHAgentAuth     bool              `yaml:"ssh_agent_auth,omitempty" json:"sshAgentAuth,omitempty"`
	//SSHKey           string            `yaml:"ssh_key" json:"sshKey,omitempty" norman:"type=password"`
	//SSHKeyPath       string            `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	//SSHCert          string            `yaml:"ssh_cert" json:"sshCert,omitempty"`
	//SSHCertPath      string            `yaml:"ssh_cert_path" json:"sshCertPath,omitempty"`
	Labels map[string]string `yaml:"labels" json:"labels,omitempty"`
	Taints []Taint           `yaml:"taints" json:"taints,omitempty"`
}

type Taint struct {
	Key    string      `json:"key,omitempty" yaml:"key"`
	Value  string      `json:"value,omitempty" yaml:"value"`
	Effect TaintEffect `json:"effect,omitempty" yaml:"effect"`
}

type TaintEffect string

const (
	TaintEffectNoSchedule       TaintEffect = "NoSchedule"
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"
	TaintEffectNoExecute        TaintEffect = "NoExecute"
)

type NodeInfo struct {
	HostName string
}

type NetworkConfig struct {
	Plugin          string `yaml:"plugin" json:"plugin,omitempty"`
	KubePodsCIDR    string `yaml:"kube_pods_cidr" json:"kube_pods_cidr,omitempty"`
	KubeServiceCIDR string `yaml:"kube_service_cidr" json:"kube_service_cidr,omitempty"`
}

type LBKubeApiserverCfg struct {
	Domain  string `yaml:"domain" json:"domain,omitempty"`
	Address string `yaml:"address" json:"address,omitempty"`
	Port    string `yaml:"port" json:"port,omitempty"`
}

type KubeadmCfg struct {
	ClusterName          string
	ControlPlaneEndpoint string
	PodSubnet            string
	ServiceSubnet        string
	ImageRepo            string
	Version              string
	CertSANs             []string
}

func (cfg *ClusterCfg) GroupHosts() (*EtcdNodes, *MasterNodes, *WorkerNodes) {
	hosts := cfg.Hosts
	etcdNodes := EtcdNodes{}
	masterNodes := MasterNodes{}
	workerNodes := WorkerNodes{}

	for _, host := range hosts {
		clusterNode := ClusterNodeCfg{Node: host}
		for _, role := range host.Role {
			if role == "etcd" {
				clusterNode.IsEtcd = true
			}
			if role == "master" {
				clusterNode.IsMaster = true
			}
			if role == "worker" {
				clusterNode.IsWorker = true
			}
		}
		if clusterNode.IsEtcd == true {
			etcdNodes.Hosts = append(etcdNodes.Hosts, clusterNode)
		}
		if clusterNode.IsMaster == true {
			masterNodes.Hosts = append(masterNodes.Hosts, clusterNode)
		}
		if clusterNode.IsWorker == true {
			workerNodes.Hosts = append(workerNodes.Hosts, clusterNode)
		}
	}
	return &etcdNodes, &masterNodes, &workerNodes
}

func (cfg *ClusterCfg) GenerateKubeadmCfg() *KubeadmCfg {
	kubeadm := KubeadmCfg{}
	kubeadm.ClusterName = DefaultClusterName
	kubeadm.PodSubnet = cfg.Network.KubePodsCIDR
	kubeadm.ServiceSubnet = cfg.Network.KubeServiceCIDR
	kubeadm.ImageRepo = cfg.KubeImageRepo
	kubeadm.Version = cfg.KubeVersion
	if cfg.LBKubeApiserver.Domain == "" {
		kubeadm.ControlPlaneEndpoint = fmt.Sprintf("%s:%s", DefaultLBDomain, DefaultLBPort)
	} else {
		kubeadm.ControlPlaneEndpoint = fmt.Sprintf("%s:%s", cfg.LBKubeApiserver.Domain, cfg.LBKubeApiserver.Port)
	}
	if cfg.LBKubeApiserver.Address == "" {
		kubeadm.ControlPlaneEndpoint = fmt.Sprintf("%s:%s", DefaultLBDomain, DefaultLBPort)
	}

	kubeadm.CertSANs = cfg.GenerateCertSANs(kubeadm.ClusterName)
	return &kubeadm
}

func (cfg *ClusterCfg) GenerateCertSANs(clusterName string) []string {
	clusterSvc := fmt.Sprintf("kubernetes.default.svc.%s", clusterName)
	defaultCertSANs := []string{"kubernetes", "kubernetes.default", "kubernetes.default.svc", clusterSvc, "localhost", "127.0.0.1"}
	extraCertSANs := []string{}
	if cfg.LBKubeApiserver.Domain == "" {
		extraCertSANs = append(extraCertSANs, DefaultLBDomain)
	} else {
		extraCertSANs = append(extraCertSANs, cfg.LBKubeApiserver.Domain)
	}
	if cfg.LBKubeApiserver.Address != "" {
		extraCertSANs = append(extraCertSANs, cfg.LBKubeApiserver.Address)
	}
	if cfg.Hosts != nil {
		for _, host := range cfg.Hosts {
			if host.HostName != "" {
				extraCertSANs = append(extraCertSANs, host.HostName)
			}
			if host.Address != "" {
				extraCertSANs = append(extraCertSANs, host.Address)
			}
			if host.InternalAddress != "" && host.InternalAddress != host.Address {
				extraCertSANs = append(extraCertSANs, host.InternalAddress)
			}
		}
	}
	if cfg.Network.KubeServiceCIDR == "" {
		extraCertSANs = append(extraCertSANs, util.ParseIp(DefaultServiceCIDR)[0])
	} else {
		extraCertSANs = append(extraCertSANs, util.ParseIp(cfg.Network.KubeServiceCIDR)[0])
	}
	defaultCertSANs = append(defaultCertSANs, extraCertSANs...)
	return defaultCertSANs
}

func (cfg *ClusterCfg) GenerateHosts() string {
	var lbHost string
	hostsList := []string{}

	_, masters, _ := cfg.GroupHosts()
	if cfg.Hosts == nil {
		localIp, err := util.GetLocalIP()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		lbHost = fmt.Sprintf("%s  %s", localIp, DefaultLBDomain)
	} else {
		if cfg.LBKubeApiserver.Address != "" {
			lbHost = fmt.Sprintf("%s  %s", cfg.LBKubeApiserver.Address, cfg.LBKubeApiserver.Domain)
		} else {
			if len(masters.Hosts) == 1 {
				lbHost = fmt.Sprintf("%s  %s", masters.Hosts[0].Node.InternalAddress, DefaultLBDomain)
			}
		}

		for _, host := range cfg.Hosts {
			if host.HostName != "" {
				hostsList = append(hostsList, fmt.Sprintf("%s  %s", host.InternalAddress, host.HostName))
			}
		}
	}

	hostsList = append(hostsList, lbHost)
	hosts := strings.Join(hostsList, "\n")
	return hosts
}
