package ssh

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type ExecInfo struct {
	Hosts      string
	Ips        string
	Cmds       string
	Username   string
	Password   string
	Key        string
	Port       int
	Ciphers    string
	CmdFile    string
	HostFile   string
	IpFile     string
	CfgFile    string
	JsonMode   bool
	OutTxt     bool
	FileLocate string
	LinuxMode  bool
	TimeLimit  int
	NumLimit   int
}

func ExecuteCmd(exec ExecInfo) {
	var cmdList []string
	var hostList []string
	//var cipherList []string
	var err error

	sshHosts := Cluster{}.Hosts
	newSshHosts := Cluster{}.Hosts
	var hostInfo Host

	if exec.IpFile != "" {
		hostList, err = GetIpListFromFile(exec.IpFile)
		if err != nil {
			log.Println("load iplist error: ", err)
			return
		}
	}

	if exec.HostFile != "" {
		hostList, err = Getfile(exec.HostFile)
		if err != nil {
			log.Println("load hostfile error: ", err)
			return
		}
	}

	if exec.Ips != "" {
		hostList, err = GetIpList(exec.Ips)
		if err != nil {
			log.Println("load iplist error: ", err)
			return
		}
	}

	if exec.Hosts != "" {
		hostList = SplitString(exec.Hosts)
	}

	if exec.CmdFile != "" {
		cmdList, err = Getfile(exec.CmdFile)
		if err != nil {
			log.Println("load cmdfile error: ", err)
		}
	}

	if exec.Cmds != "" {
		cmdList = SplitString(exec.Cmds)
	}

	//if exec.Ciphers != "" {
	//	cipherList = SplitString(exec.Ciphers)
	//}

	if exec.CfgFile == "" {
		for _, host := range hostList {
			hostInfo.Ip = host
			hostInfo.User = exec.Username
			hostInfo.Passwd = exec.Password
			hostInfo.Port = exec.Port
			hostInfo.CmdList = cmdList
			hostInfo.Key = exec.Key
			hostInfo.LinuxMode = exec.LinuxMode
			newSshHosts = append(newSshHosts, hostInfo)
		}
	} else {
		log.Info("Get config info")
		sshHosts, err = GetYamlFile(exec.CfgFile)
		if err != nil {
			log.Println("load cfgFile error: ", err)
		}
		for _, host := range sshHosts {
			//log.Info(host)
			hostInfo.Ip = host.Ip
			hostInfo.User = host.User
			hostInfo.Passwd = host.Passwd
			hostInfo.Key = host.Key
			hostInfo.Port = host.Port
			hostInfo.Roles = host.Roles
			hostInfo.CmdList = cmdList
			hostInfo.LinuxMode = exec.LinuxMode
			newSshHosts = append(newSshHosts, hostInfo)
		}

		//for i := 0; i < len(sshHosts); i++ {
		//	if sshHosts[i].Cmds != "" {
		//		sshHosts[i].CmdList = SplitString(sshHosts[i].Cmds)
		//	} else {
		//		cmdList, err = Getfile(sshHosts[1].CmdFile)
		//		if err != nil {
		//			log.Println("load cmdFile error: ", err)
		//			return
		//		}
		//		sshHosts[i].CmdList = cmdList
		//	}
		//}
	}

	//chLimit := make(chan bool, exec.NumLimit)
	//chs := make([]chan SSHResult, len(sshHosts))
	startTime := time.Now()
	log.Println("Welcome to KubeOcean")
	wg := sync.WaitGroup{}
	//limitFunc := func(chLimit chan bool, host Host) {
	//	wg.Add(1)
	//	Dossh(host.Username, host.Passwd, host.Ip, host.Key, host.CmdList, host.Port, exec.TimeLimit, cipherList, host.LinuxMode)
	//	log.Info(host.Roles)
	//	<-chLimit
	//	//log.Info(a)
	//}
	limitFunc := func(host Host) {
		//Dossh(host.User, host.Passwd, host.Ip, host.Key, host.CmdList, host.Port, exec.TimeLimit, cipherList, host.LinuxMode)
		log.Info(host.Roles)
		//a :=  <-chLimit
		//log.Info(a)
	}
	for _, host := range newSshHosts {
		//chLimit<- true
		//go limitFunc(chLimit, host)
		limitFunc(host)
	}

	//sshResults := []SSHResult{}
	//for _, ch := range chs {
	//	res := <-ch
	//	if res.Result != "" {
	//		sshResults = append(sshResults, res)
	//	}
	//}
	wg.Wait()
	endTime := time.Now()
	log.Printf("KubeOcean finished. Process time %s. Number of active ip is %d", endTime.Sub(startTime), len(sshHosts))

	//if exec.OutTxt {
	//	for _, sshResult := range sshResults {
	//		err = WriteIntoTxt(sshResult, exec.FileLocate)
	//		if err != nil {
	//			log.Println("write into txt error: ", err)
	//			return
	//		}
	//	}
	//	return
	//}
	//if exec.JsonMode {
	//	jsonResult, err := json.Marshal(sshResults)
	//	if err != nil {
	//		log.Println("json Marshal error: ", err)
	//	}
	//	fmt.Println(string(jsonResult))
	//	return
	//}
	//for _, sshResults := range sshResults {
	//	fmt.Println("host: ", sshResults.Host)
	//	fmt.Println("========= Result =========")
	//	fmt.Println(sshResults.Result)
	//}
}
