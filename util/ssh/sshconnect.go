package ssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

//ssh connect
func connect(user, password, host, key string, port int, cipherList []string) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
		log.Infof("passwd is %v", password)
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	//if len(cipherList) == 0 {
	//	config = ssh.Config{
	//		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	//	}
	//} else {
	//	config = ssh.Config{
	//		Ciphers: cipherList,
	//	}
	//}
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	log.Infof("%v", clientConfig.User)
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	log.Info(addr)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		log.Errorf("ssh dial failed: %v", err)
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}

func SftpConnect(user, password, host, key string, port int, cipherList []string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(cipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: cipherList,
		}
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(client); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func Dossh(username, password, host, key string, cmdlist []string, port, timeout int, cipherList []string, linuxMode bool) {
	//chSSH := make(chan SSHResult)
	//log.Info(linuxMode)
	if linuxMode {
		log.Info("ssh exec")
		//go dossh_run(username, password, host, key, cmdlist, port, cipherList)
		dossh_run(username, password, host, key, cmdlist, port, cipherList)
	} else {
		//go dossh_session(username, password, host, key, cmdlist, port, cipherList, chSSH)
		log.Warn("no support session")
	}
	//var res SSHResult

	//select {
	//case <-time.After(time.Duration(timeout) * time.Second):
	//	res.Host = host
	//	res.Success = false
	//	res.Result = ("SSH run timeout：" + strconv.Itoa(timeout) + " second.")
	//	ch <- res
	//case res = <-chSSH:
	//	ch <- res
	//}
	//return
}

func dossh_session(username, password, host, key string, cmdlist []string, port int, cipherList []string, ch chan SSHResult) {
	session, err := connect(username, password, host, key, port, cipherList)
	var sshResult SSHResult
	sshResult.Host = host

	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>", err.Error())
		ch <- sshResult
		return
	}
	defer session.Close()

	cmdlist = append(cmdlist, "exit")

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	err = session.Shell()
	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>", err.Error())
		ch <- sshResult
		return
	}
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
	}
	session.Wait()
	if errbt.String() != "" {
		sshResult.Success = false
		sshResult.Result = errbt.String()
		ch <- sshResult
	} else {
		sshResult.Success = true
		sshResult.Result = outbt.String()
		ch <- sshResult
	}

	return
}

func dossh_run(username, password, host, key string, cmdlist []string, port int, cipherList []string) {
	session, err := connect(username, password, host, key, port, cipherList)
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()

	newCmdList := []string{}
	for _, cmd := range cmdlist {
		if strings.Contains(cmd, "sudo") || strings.Index(cmd, "sudo") == 0 {
			cmd = strings.Replace(cmd, "sudo", "sudo -S", 1)
			cmd = fmt.Sprint("echo ", password, " | ", cmd)
		}
		newCmdList = append(newCmdList, cmd)
	}

	//newCmdList = append(newCmdList, "exit")
	newcmd := strings.Join(newCmdList, "&&")
	//fmt.Println(newcmd)
	//var outbt, errbt bytes.Buffer
	//session.Stdout = &outbt
	//
	//session.Stderr = &errbt
	log.Info(newcmd)
	b, err := session.CombinedOutput(newcmd)
	if err != nil {
		log.Fatal("任务执行失败：", err)
		return
	}
	fmt.Printf("执行结果：\n%s", string(b))
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)

	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)

	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)

	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + "  copy file to remote server finished!")
}
