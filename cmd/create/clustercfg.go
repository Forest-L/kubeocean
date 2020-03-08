package create

import (
	"bufio"
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func NewCmdCreateCfg() *cobra.Command {
	var culsterCfgCmd = &cobra.Command{
		Use:   "config",
		Short: "Create cluster info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("test")
			clusterConfig()
		},
	}
	return culsterCfgCmd
}

func getConfig(reader *bufio.Reader, text, def string) (string, error) {
	for {
		if def == "" {
			fmt.Printf("[+] %s [%s]: ", text, "none")
		} else {
			fmt.Printf("[+] %s [%s]: ", text, def)
		}
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)

		if input != "" {
			return input, nil
		}
		return def, nil
	}
}

func writeConfig(cluster *cluster.ClusterCfg, configFile string, print bool) error {
	yamlConfig, err := yaml.Marshal(*cluster)
	if err != nil {
		return err
	}
	log.Debugf("Deploying cluster configuration file: %s", configFile)

	configString := fmt.Sprintf("%s", string(yamlConfig))
	if print {
		fmt.Printf("Configuration File: \n%s", configString)
		return nil
	}
	return ioutil.WriteFile(configFile, []byte(configString), 0640)
}

func clusterConfig() error {
	clusterCfg := cluster.ClusterCfg{}
	reader := bufio.NewReader(os.Stdin)

	// Get number of hosts
	numberOfHostsString, err := getConfig(reader, "Number of Hosts", "1")
	if err != nil {
		return err
	}
	numberOfHostsInt, err := strconv.Atoi(numberOfHostsString)
	if err != nil {
		return err
	}

	sshKeyPath, err := getConfig(reader, "Cluster Level SSH Private Key Path", "~/.ssh/id_rsa")
	if err != nil {
		return err
	}
	clusterCfg.SSHKeyPath = sshKeyPath
	// Get Hosts config
	clusterCfg.Nodes = make([]cluster.NodeCfg, 0)
	for i := 0; i < numberOfHostsInt; i++ {
		hostCfg, err := getHostConfig(reader, i, clusterCfg.SSHKeyPath)
		if err != nil {
			return err
		}
		clusterCfg.Nodes = append(clusterCfg.Nodes, *hostCfg)
	}

	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	println(exPath)
	return writeConfig(&clusterCfg, exPath+"cluster-info.yaml", true)
}

func getHostConfig(reader *bufio.Reader, index int, clusterSSHKeyPath string) (*cluster.NodeCfg, error) {
	host := cluster.NodeCfg{}

	address, err := getConfig(reader, fmt.Sprintf("SSH Address of host (%d)", index+1), "")
	if err != nil {
		return nil, err
	}
	host.Address = address

	port, err := getConfig(reader, fmt.Sprintf("SSH Port of host (%d)", index+1), cluster.DefaultSSHPort)
	if err != nil {
		return nil, err
	}
	host.Port = port

	sshKeyPath, err := getConfig(reader, fmt.Sprintf("SSH Private Key Path of host (%s)", address), "")
	if err != nil {
		return nil, err
	}
	if len(sshKeyPath) == 0 {
		fmt.Printf("[-] You have entered empty SSH key path, trying fetch from SSH key parameter\n")
		sshKey, err := getConfig(reader, fmt.Sprintf("SSH Private Key of host (%s)", address), "")
		if err != nil {
			return nil, err
		}
		if len(sshKey) == 0 {
			fmt.Printf("[-] You have entered empty SSH key, defaulting to cluster level SSH key: %s\n", clusterSSHKeyPath)
			host.SSHKeyPath = clusterSSHKeyPath
		} else {
			host.SSHKey = sshKey
		}
	} else {
		host.SSHKeyPath = sshKeyPath
	}

	sshUser, err := getConfig(reader, fmt.Sprintf("SSH User of host (%s)", address), "root")
	if err != nil {
		return nil, err
	}
	host.User = sshUser

	hostRole, err := getConfig(reader, fmt.Sprintf("What's host (%s) role?(0: etcd, 1: master, 2: worker)", address), "012")
	if err != nil {
		return nil, err
	}

	if strings.Contains(hostRole, "0") {
		host.Role = append(host.Role, cluster.ETCDRole)
	}
	if strings.Contains(hostRole, "1") {
		host.Role = append(host.Role, cluster.MasterRole)
	}
	if strings.Contains(hostRole, "2") {
		host.Role = append(host.Role, cluster.WorkerRole)
	}

	hostnameOverride, err := getConfig(reader, fmt.Sprintf("Override Hostname of host (%s)", address), "")
	if err != nil {
		return nil, err
	}
	host.HostnameOverride = hostnameOverride

	internalAddress, err := getConfig(reader, fmt.Sprintf("Internal IP of host (%s)", address), "")
	if err != nil {
		return nil, err
	}
	host.InternalAddress = internalAddress

	return &host, nil
}
