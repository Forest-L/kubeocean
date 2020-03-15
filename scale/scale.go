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
	tokenCreateMasterCmd := fmt.Sprintf("/usr/local/bin/kubeadm token create --print-join-command --certificate-key %s", CertificateKey)
	outMasterCmd, errMaster := master.CmdExecOut(tokenCreateMasterCmd)
	if errMaster != nil {
		log.Fatalf("Failed to create token (%s):\n", master.Node.Address)
		os.Exit(1)
	}
	joinMasterStrList := strings.Split(outMasterCmd, "kubeadm join")
	joinMasterStr := strings.Split(joinMasterStrList[1], CertificateKey)
	joinMasterCmd = fmt.Sprintf("/usr/local/bin/kubeadm join %s %s", joinMasterStr[0], CertificateKey)

	joinWorkerStrList := strings.Split(joinMasterCmd, "--control-plane")
	joinWorkerCmd = joinWorkerStrList[0]

	return joinMasterCmd, joinWorkerCmd
}
