package tmpl

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var (
	KubeletServiceTempl = template.Must(template.New("kubeletService").Parse(
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

	KubeletContainerTempl = template.Must(template.New("kubeletContainer").Parse(
		dedent.Dedent(`[Service]
Environment="KUBELET_CONTAINER= --net=host \
  --pid=host \
  --privileged \
  --name=kubelet \
  --restart=on-failure:5 \
  --memory=256M \
  --cpu-shares=100 \
  -v /dev:/dev:rw \
  -v /etc/cni:/etc/cni:ro \
  -v /opt/cni:/opt/cni:ro \
  -v /etc/ssl:/etc/ssl:ro \
  -v /etc/resolv.conf:/etc/resolv.conf \
  -v /etc/calico/certs:/etc/calico/certs:ro \
  -v /usr/share/ca-certificates:/usr/share/ca-certificates:ro \
  -v /sys:/sys:ro \
  -v /var/lib/docker:/var/lib/docker:rw \
  -v /var/log:/var/log:rw \
  -v /var/lib/kubelet:/var/lib/kubelet:shared \
  -v /var/lib/calico:/var/lib/calico:shared \
  -v /var/lib/cni:/var/lib/cni:shared \
  -v /var/run:/var/run:rw \
  -v /etc/kubernetes:/etc/kubernetes:shared \
  -v /usr/libexec/kubernetes:/usr/libexec/kubernetes:shared \
  -v /etc/os-release:/etc/os-release:ro \
  {{ .Repo }}/google-containers/hyperkube:{{ .Version }} \
  kubelet "
    `)))

	KubeletTempl = template.Must(template.New("kubeletContainer").Parse(
		dedent.Dedent(`#!/bin/bash
/usr/bin/docker run --rm {{ .Repo }}/google-containers/hyperkube:{{ .Version }} kubelet "$@"
    `)))
)
