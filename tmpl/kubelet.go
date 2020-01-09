package tmpl

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var (
	kubeletServiceTempl = template.Must(template.New("kubeletService").Parse(
		dedent.Dedent(`[Unit]
Description=Kubernetes Kubelet Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=docker.service
Wants=docker.socket

[Service]
User=root
EnvironmentFile=/etc/kubernetes/kubelet.env
ExecStart=/usr/local/bin/kubelet \
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

	kubeletContainerTempl = template.Must(template.New("kubeletContainer").Parse(
		dedent.Dedent(`#!/bin/bash
/usr/bin/docker run \
  --net=host \
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
  -v /etc/kubernetes:/etc/kubernetes:ro \
  -v /etc/os-release:/etc/os-release:ro \
  {{.KubeRepo}}/hyperkube:{{.KubeVersion}} \
  ./hyperkube kubelet \
  "$@"
    `)))
)
