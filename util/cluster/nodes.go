package cluster

type nodeCfg struct {
	HostName         string            `yaml:"hostName,omitempty" json:"hostName,omitempty" norman:"type=reference[node]"`
	Address          string            `yaml:"address" json:"address,omitempty"`
	Port             string            `yaml:"port" json:"port,omitempty"`
	InternalAddress  string            `yaml:"internal_address" json:"internalAddress,omitempty"`
	Role             []string          `yaml:"role" json:"role,omitempty" norman:"type=array[enum],options=etcd|worker|worker"`
	HostnameOverride string            `yaml:"hostname_override" json:"hostnameOverride,omitempty"`
	User             string            `yaml:"user" json:"user,omitempty"`
	SSHAgentAuth     bool              `yaml:"ssh_agent_auth,omitempty" json:"sshAgentAuth,omitempty"`
	SSHKey           string            `yaml:"ssh_key" json:"sshKey,omitempty" norman:"type=password"`
	SSHKeyPath       string            `yaml:"ssh_key_path" json:"sshKeyPath,omitempty"`
	SSHCert          string            `yaml:"ssh_cert" json:"sshCert,omitempty"`
	SSHCertPath      string            `yaml:"ssh_cert_path" json:"sshCertPath,omitempty"`
	Labels           map[string]string `yaml:"labels" json:"labels,omitempty"`
	Taints           []Taint           `yaml:"taints" json:"taints,omitempty"`
}

type Taint struct {
	Key    string      `json:"key,omitempty" yaml:"key"`
	Value  string      `json:"value,omitempty" yaml:"value"`
	Effect TaintEffect `json:"effect,omitempty" yaml:"effect"`
}

type TaintEffect string

const (
	TaintEffectNoSchedule       TaintEffect = "NoSchedule"
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"
	TaintEffectNoExecute        TaintEffect = "NoExecute"
)
