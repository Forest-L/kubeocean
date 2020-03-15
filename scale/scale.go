package scale

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

func GetJoinCmd(master *cluster.ClusterNodeCfg) (string, string) {
	var joinMasterCmd, joinWorkerCmd string

	// Get Join Master Command
	uploadCertsCmd := "/usr/local/bin/kubeadm init phase upload-certs --upload-certs"
	out, err := master.CmdExecOut(uploadCertsCmd)
	if err != nil {
		log.Fatalf("Failed to upload-certs (%s):\n", master.Node.Address)
		os.Exit(1)
	}
	reg := regexp.MustCompile("[0-9|a-z]{64}")
	CertificateKey := reg.FindAllString(out, -1)[0]

	// Get Join Worker Command
	tokenCreateWorkerCmd := "/usr/local/bin/kubeadm token create --print-join-command"
	outWorkerCmd, errWorker := master.CmdExecOut(tokenCreateWorkerCmd)
	if errWorker != nil {
		log.Fatalf("Failed to create token (%s):\n", master.Node.Address)
		os.Exit(1)
	}
	joinWorkerStrList := strings.Split(outWorkerCmd, "kubeadm join")
	joinWorkerStr := strings.Split(joinWorkerStrList[1], "\n")
	joinWorkerCmd = fmt.Sprintf("/usr/local/bin/kubeadm join %s", joinWorkerStr[0])
	joinMasterCmd = fmt.Sprintf("/usr/local/bin/kubeadm join %s --control-plane --certificate-key %s", joinWorkerStr[0], CertificateKey)
	return joinMasterCmd, joinWorkerCmd
}
