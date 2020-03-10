package tmpl

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var KubeadmCfgTempl = template.Must(template.New("kubeadmCfg").Parse(
	dedent.Dedent(`[Unit]
Description=Kubernetes Kubelet Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=docker.service
Wants=docker.socket

[Service]
User=root
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
EnvironmentFile=-/etc/default/kubelet
EnvironmentFile=/var/lib/kubelet/kubeadm-flags.env
ExecStart=/usr/bin/docker run $KUBELET_CONTAINER \
        $KUBELET_KUBECONFIG_ARGS \
        $KUBELET_CONFIG_ARGS  \
        $KUBELET_KUBEADM_ARGS  \
        $KUBELET_EXTRA_ARGS  \
		$KUBE_LOGTOSTDERR \
		$KUBE_LOG_LEVEL \
		$KUBELET_API_SERVER \
		$KUBELET_ADDRESS \
		$KUBELET_PORT \
		$KUBELET_HOSTNAME \
		$KUBE_ALLOW_PRIV \
		$KUBELET_ARGS \
		$DOCKER_SOCKET \
		$KUBELET_NETWORK_PLUGIN \
		$KUBELET_VOLUME_PLUGIN \
		$KUBELET_CLOUDPROVIDER
Restart=always
RestartSec=10s
ExecStartPre=-/usr/bin/docker rm -f kubelet
ExecStartPre=-/bin/mkdir -p /var/lib/kubelet/volume-plugins
ExecReload=/usr/bin/docker restart kubelet


[Install]
WantedBy=multi-user.target
    `)))
