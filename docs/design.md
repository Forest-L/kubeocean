###功能
* 命令行部署k8s集群
* 使用yaml配置集群部署信息
* 支持多台机器批量执行任务
* 支持自定义证书有效期

###结构体

``` go

type Host struct {
	Name      string    `yaml:"name"`
	Ip        string    `yaml:"ip"`
	Port      int       `yaml:"port"`
	Username  string    `yaml:"username"`
	Passwd    string    `yaml:"passwd"`
	CmdFile   string    `yaml:"cmdfile"`
	Cmds      string    `yaml:"cmds"`
	CmdList   []string  `yaml:"cmdlist"`
	Key       string    `yaml:"key"`
	LinuxMode bool      `yaml:"linuxmode"`
	Result    SSHResult `yaml:"result"`
	Roles     []string  `yaml:"roles"`
}

type Cluster struct {
	Hosts []Host `yaml:"hosts"`
}

```