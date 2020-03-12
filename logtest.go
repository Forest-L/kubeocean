package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//var joinMasterCmd string
	cmd := exec.Command("sudo", "/usr/local/bin/kubeadm init phase upload-certs --upload-certs")

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("v%", err)
		os.Exit(1)
	}
	fmt.Println(string(out))

	outList := strings.Split(string(out), "Using certificate key:\r\n")
	certificateKeyStr := strings.Split(outList[1], "\r\n")
	CertificateKey := certificateKeyStr[0]
	fmt.Println(CertificateKey)

	//tokenCreateCmd := fmt.Sprintf("sudo /usr/local/bin/kubeadm token create --print-join-command --certificate-key %s", CertificateKey)
	//out, err1 := ssh.CmdExecOut(master.Node.Address, master.Node.User, master.Node.Port, master.Node.Password, false, tokenCreateCmd)
	//if err1 != nil {
	//	log.Fatalf("Failed to create token (%s):\n", master.Node.Address)
	//	os.Exit(1)
	//}
	//joinStrList := strings.Split(out, "kubeadm join")
	//joinStr := strings.Split(joinStrList[1], "\r\n")
	//joinMasterCmd = fmt.Sprintf("sudo /usr/local/bin/kubeadm join %s", joinStr[0])
}
