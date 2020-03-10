package ssh

func PushFile(host string, src string, dst string, user string, port string, password string, force bool) {
	puser := NewUser(user, port, password, force)
	SinglePush(host, src, dst, puser, force)
}

func SinglePush(ip, src, dst string, cu *CommonUser, f bool) {
	server := NewScpServer(ip, cu.port, cu.user, cu.psw, "scp", src, dst, f)
	cmd := "push " + server.FileName + " to " + server.Ip + ":" + server.RemotePath

	rs := Result{
		Ip:  server.Ip,
		Cmd: cmd,
	}
	err := server.RunScpDir()
	if err != nil {
		rs.Err = err
	} else {
		rs.Result = cmd + " ok\n"
	}
	PrintExecResult(rs)
}
