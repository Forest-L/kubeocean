package ssh

import "errors"

func CmdExec(host string, user string, port string, password string, force bool, encFlag bool, cmd string) {
	puser := NewUser(user, port, password, force, encFlag)
	SingleRun(host, cmd, puser, force)
}

func SingleRun(host string, cmd string, user *CommonUser, force bool) {
	server := NewCmdServer(host, user.port, user.user, user.psw, "cmd", cmd, force)
	r := server.SRunCmd()
	PrintExecResult(r)
}
