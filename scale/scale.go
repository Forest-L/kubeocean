package scale

import (
	"fmt"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func JoinMasterCmd(masters *cluster.MasterNodes) string {
	var joinMasterCmd string
	master := masters.Hosts[0]
	uploadCertsCmd := "sudo /usr/local/bin/kubeadm init phase upload-certs --upload-certs"
	out, err := ssh.CmdExecOut(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, uploadCertsCmd)
	if err != nil {
		log.Fatalf("Failed to upload-certs (%s):\n", master.Node.Address)
		os.Exit(1)
	}
	outList := strings.Split(out, "Using certificate key:\n")
	certificateKeyStr := strings.Split(outList[1], "\n")
	CertificateKey := certificateKeyStr[0]
	tokenCreateCmd := fmt.Sprintf("sudo /usr/local/bin/kubeadm token create --print-join-command --certificate-key %s", CertificateKey)
	out, err1 := ssh.CmdExecOut(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, tokenCreateCmd)
	if err1 != nil {
		log.Fatalf("Failed to create token (%s):\n", master.Node.Address)
		os.Exit(1)
	}
	joinStrList := strings.Split(out, "kubeadm join")
	joinStr := strings.Split(joinStrList[1], "\n")
	joinMasterCmd = fmt.Sprintf("sudo /usr/local/bin/kubeadm join %s", joinStr[0])
	return joinMasterCmd
}

func JoinWorkerCmd(masters *cluster.MasterNodes) string {
	var joinWorkerCmd string
	master := masters.Hosts[0]

	tokenCreateCmd := "sudo /usr/local/bin/kubeadm token create --print-join-command"
	out, err1 := ssh.CmdExecOut(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, tokenCreateCmd)
	if err1 != nil {
		log.Fatalf("Failed to create token (%s):\n", master.Node.Address)
		os.Exit(1)
	}
	joinStrList := strings.Split(out, "kubeadm join")
	joinStr := strings.Split(joinStrList[1], "\n")
	joinWorkerCmd = fmt.Sprintf("sudo /usr/local/bin/kubeadm join %s", joinStr[0])
	return joinWorkerCmd
}
