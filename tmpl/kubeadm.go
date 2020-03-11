package tmpl

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var KubeadmCfgTempl = template.Must(template.New("kubeadmCfg").Parse(
	dedent.Dedent(`---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
etcd:
  local:
    imageRepository: "quay.azk8s.cn/coreos"
    imageTag: "v3.3.10"
    dataDir: "/var/lib/etcd"
dns:
  type: CoreDNS
  imageRepository: coredns
  imageTag: 1.6.0
imageRepository: {{ .ImageRepo }}/google-containers
kubernetesVersion: {{ .Version }}
certificatesDir: /etc/kubernetes/ssl
clusterName: {{ .ClusterName }}
controlPlaneEndpoint: {{ .ControlPlaneEndpoint }}
networking:
  dnsDomain: {{ .ClusterName }}
  podSubnet: {{ .PodSubnet }}
  serviceSubnet: {{ .ServiceSubnet }}
apiServer:
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
  certSANs:
    {{- range .CertSANs }}
    - {{ . }}
    {{- end }}
    `)))
