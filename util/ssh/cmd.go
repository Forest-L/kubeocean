package ssh

import (
	"fmt"
	"sync"
)

func CmdExec(host string, user string, port string, password string, force bool, privilegeCmd string, cmd string) error {
	puser := NewUser(user, port, password, force)
	if privilegeCmd != "" {
		cmd = fmt.Sprintf("%s\"%s\"", privilegeCmd, cmd)
	}
	err := SingleRun(host, cmd, puser, force)
	if err != nil {
		return err
	}
	return nil
}

func SingleRun(host string, cmd string, user *CommonUser, force bool) error {
	server := NewCmdServer(host, user.port, user.user, user.psw, "cmd", cmd, force)
	r := server.SRunCmd()
	//PrintExecResult(r)
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func CmdExecOut(host string, user string, port string, password string, force bool, privilegeCmd string, cmd string) (string, error) {
	puser := NewUser(user, port, password, force)
	if privilegeCmd != "" {
		cmd = fmt.Sprintf("%s\"%s\"", privilegeCmd, cmd)
	}
	out, err := SingleRunOut(host, cmd, puser, force)
	if err != nil {
		return "", err
	} else {
		return out, nil
	}
}

func SingleRunOut(host string, cmd string, user *CommonUser, force bool) (string, error) {
	server := NewCmdServer(host, user.port, user.user, user.psw, "cmd", cmd, force)
	r := server.SRunCmd()
	//PrintExecResult(r)
	if r.Err != nil {
		return "", r.Err
	} else {
		return r.Result, nil
	}
}

func GetIps(h []Host) []string {
	ips := make([]string, 0)
	for _, v := range h {
		ips = append(ips, v.Ip)
	}
	return ips
}

//func ServersRun(cmd string, cu *CommonUser, wt *sync.WaitGroup, crs chan machine.Result, ipFile string, ccons chan struct{}) {
func ServersRun(cmd string, hosts []Host) {

	result := make(chan Result)
	ccons := make(chan struct{}, DefaultCon)
	hostNum := len(hosts)
	wg := &sync.WaitGroup{}

	go PrintResults(result, hostNum, wg, ccons)

	for _, h := range hosts {
		ccons <- struct{}{}
		//if h.PrivilegeCmd != "" {
		//	cmd = fmt.Sprintf("%s\"%s\"", h.PrivilegeCmd, cmd)
		//}
		server := NewCmdServer(h.Ip, h.Port, h.User, h.Psw, "cmd", cmd, true)
		wg.Add(1)
		go server.PRunCmd(result)
	}
}

//push file or dir to remote servers
func ServersPush(src, dst string, hosts []Host) {
	result := make(chan Result)
	ccons := make(chan struct{}, DefaultCon)
	hostNum := len(hosts)
	wg := &sync.WaitGroup{}

	go PrintResults(result, hostNum, wg, ccons)

	for _, h := range hosts {
		ccons <- struct{}{}
		server := NewScpServer(h.Ip, h.Port, h.User, h.Psw, "scp", src, dst, true)
		wg.Add(1)
		go server.PRunScp(result)
	}
}
