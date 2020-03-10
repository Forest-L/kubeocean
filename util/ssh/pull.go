package ssh

func PullFile(host string, src string, dst string, user string, port string, password string, force bool) {
	puser := NewUser(user, port, password, force)
	SinglePull(host, src, dst, puser, force)
}

func SinglePull(host string, src string, dst string, cu *CommonUser, force bool) {
	server := NewPullServer(host, cu.port, cu.user, cu.psw, "scp", src, dst, force)
	err := server.PullScp()
	PrintPullResult(host, src, dst, err)
}
