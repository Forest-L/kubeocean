package ssh

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	NO_PASSWORD = "GET PASSWORD ERROR\n"
)

const (
	NO_EXIST = "0"
	IS_FILE  = "1"
	IS_DIR   = "2"
)

type CommonUser struct {
	user  string
	port  string
	psw   string
	force bool
}

func NewUser(user, port, psw string, force bool) *CommonUser {
	return &CommonUser{
		user:  user,
		port:  port,
		psw:   psw,
		force: force,
	}
}

type Server struct {
	Ip         string
	Port       string
	User       string
	Psw        string
	Action     string
	Cmd        string
	FileName   string
	RemotePath string
	Force      bool
}

type ScpConfig struct {
	Src string
	Dst string
}

type Result struct {
	Ip     string
	Cmd    string
	Result string
	Err    error
}

func NewCmdServer(ip, port, user, psw, action, cmd string, force bool) *Server {
	server := &Server{
		Ip:     ip,
		Port:   port,
		User:   user,
		Action: action,
		Cmd:    cmd,
		Psw:    psw,
		Force:  force,
	}
	return server
}

func NewScpServer(ip, port, user, psw, action, file, rpath string, force bool) *Server {
	rfile := path.Join(rpath, path.Base(file))
	cmd := createShell(rfile)
	server := &Server{
		Ip:         ip,
		Port:       port,
		User:       user,
		Psw:        psw,
		Action:     action,
		FileName:   file,
		RemotePath: rpath,
		Cmd:        cmd,
		Force:      force,
	}
	return server
}
func NewPullServer(ip, port, user, psw, action, file, rpath string, force bool) *Server {
	cmd := createShell(rpath)
	server := &Server{
		Ip:         ip,
		Port:       port,
		User:       user,
		Psw:        psw,
		Action:     action,
		FileName:   file,
		RemotePath: rpath,
		Cmd:        cmd,
		Force:      force,
	}
	return server
}

// set Server.Cmd
func (s *Server) SetCmd(cmd string) {
	s.Cmd = cmd
}

//create shell script for running on remote server
func createShell(file string) string {
	s1 := "bash << EOF \n"
	s2 := "if [[ -f " + file + " ]];then \n"
	s3 := "echo '1'\n"
	s4 := "elif [[ -d " + file + " ]];then \n"
	s5 := `echo "2"
else 
echo "0"
fi
EOF`
	cmd := s1 + s2 + s3 + s4 + s5
	return cmd
}

// implement ssh auth method [password keyboard-interactive] and [password]
func (server *Server) getSshClient() (client *ssh.Client, err error) {
	authMethods := []ssh.AuthMethod{}
	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {

		if len(questions) == 0 {
			return []string{}, nil
		}
		/*
			for i, question := range questions {
				log.Debug("SSH Question %d: %s", i+1, question)
			}
		*/

		answers = make([]string, len(questions))
		for i := range questions {
			yes, _ := regexp.MatchString("*yes*", questions[i])
			if yes {
				fmt.Println("yes")
				answers[i] = "yes"

			} else {
				fmt.Println("passwd")
				answers[i] = server.Psw
			}
		}
		return answers, nil
	}
	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	authMethods = append(authMethods, ssh.Password(server.Psw))

	sshConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	ip_port := server.Ip + ":" + server.Port
	client, err = ssh.Dial("tcp", ip_port, sshConfig)
	return
}

func (server *Server) SRunCmd() Result {
	rs := Result{
		Ip:  server.Ip,
		Cmd: server.Cmd,
	}

	if server.Psw == NO_PASSWORD {
		rs.Err = errors.New(NO_PASSWORD)
		return rs
	}

	client, err := server.getSshClient()
	if err != nil {
		rs.Err = err
		return rs
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		rs.Err = err
		return rs
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Fatalf("%v", err)
	}

	in, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	out, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	var output []byte

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if (strings.HasPrefix(line, "[sudo] password for ") || strings.HasPrefix(line, "Password")) && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(server.Psw + "\n"))
				if err != nil {
					break
				}
			}
		}
	}(in, out, &output)

	cmd := server.Cmd
	_, err1 := session.Output(cmd)
	if err1 != nil {
		rs.Err = err1
		return rs
	}
	rs.Result = string(output)
	return rs
}

func (server *Server) RunCmd() (result string, err error) {
	if server.Psw == NO_PASSWORD {
		return NO_PASSWORD, nil
	}
	client, err := server.getSshClient()
	if err != nil {
		return "getSSHClient error", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "newSession error", err
	}
	defer session.Close()

	cmd := server.Cmd
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(bs), err
	}
	return string(bs), nil
}

func (server *Server) checkRemoteFile() (result string) {
	re, _ := server.RunCmd()
	return re
}

func (server *Server) RunScpDir() (err error) {
	re := strings.TrimSpace(server.checkRemoteFile())
	log.Debug("server.checkRemoteFile()=%s\n", re)

	//远程机器存在同名文件
	if re == IS_FILE && server.Force == false {
		errString := "<ERROR>\nRemote Server's " + server.RemotePath + " has the same file " + server.FileName + "\nYou can use `-f` option force to cover the remote file.\n</ERROR>\n"
		return errors.New(errString)
	}

	rfile := server.RemotePath
	cmd := createShell(rfile)
	server.SetCmd(cmd)
	re = strings.TrimSpace(server.checkRemoteFile())
	log.Debug("server.checkRemoteFile()=%s\n", re)

	//远程目录不存在
	if re != IS_DIR {
		errString := "[" + server.Ip + ":" + server.RemotePath + "] does not exist or not a dir\n"
		return errors.New(errString)
	}

	client, err := server.getSshClient()
	if err != nil {
		return err
	}
	defer client.Close()

	filename := server.FileName
	fi, err := os.Stat(filename)
	if err != nil {
		log.Debug("open source file %s error\n", filename)
		return err
	}
	scp := NewScp(client)
	if fi.IsDir() {
		err = scp.PushDir(filename, server.RemotePath)
		return err
	}
	err = scp.PushFile(filename, server.RemotePath)
	return err
}

//pull file from remote to local server
func (server *Server) PullScp() (err error) {

	//判断远程源文件情况
	re := strings.TrimSpace(server.checkRemoteFile())
	log.Debug("server.checkRemoteFile()=%s\n", re)

	//不存在报错
	if re == NO_EXIST {
		errString := "Remote Server's " + server.RemotePath + " doesn't exist.\n"
		return errors.New(errString)
	}

	//不支持拉取目录
	if re == IS_DIR {
		errString := "Remote Server's " + server.RemotePath + " is a directory ,not support.\n"
		return errors.New(errString)
	}

	//仅仅支持普通文件
	if re != IS_FILE {
		errString := "Get info from Remote Server's " + server.RemotePath + " error.\n"
		return errors.New(errString)
	}

	//本地目录
	dst := server.FileName
	//远程文件
	src := server.RemotePath

	log.Debug("src=%s", src)
	log.Debug("dst=%s", dst)

	//本地路径不存在，自动创建
	err = MakePath(dst)
	if err != nil {
		return err
	}

	//检查本地是否有同名文件
	fileName := filepath.Base(src)
	localFile := filepath.Join(dst, fileName)

	flag := FileExists(localFile)
	log.Debug("flag=%v", flag)
	log.Debug("localFile=%s", localFile)

	//-f 可以强制覆盖
	if flag && !server.Force {
		return errors.New(localFile + " is exist, use -f to cover the old file")
	}

	//执行pull
	client, err := server.getSshClient()
	if err != nil {
		return err
	}
	defer client.Close()

	scp := NewScp(client)
	err = scp.PullFile(dst, src)
	return err
}
