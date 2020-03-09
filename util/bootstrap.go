package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func SystemInit() error {

	sysctllist := []string{"net.ipv4.ip_local_reserved_ports=30000-32767", "net.bridge.bridge-nf-call-iptables=1", "net.bridge.bridge-nf-call-arptables=1", "net.bridge.bridge-nf-call-ip6tables"}
	for _, conf := range sysctllist {
		if err := exec.Command("/bin/sh", "-c", "sysctl", conf).Run(); err != nil {
			return fmt.Errorf("failed to sysctl: %v", err)
		}
	}
	return nil
}

func SwapOff() error {
	if err := exec.Command("swapoff", "-a").Run(); err != nil {
		return fmt.Errorf("failed to swapoff: %v", err)
	}

	disableswap := fmt.Sprint(`sed -i /^[^#]*swap*/s/^/\#/g /etc/fstab`)
	if err := exec.Command("/bin/sh", "-c", disableswap).Run(); err != nil {
		return fmt.Errorf("failed to swapoff: %v", err)
	}
	return nil
}

func Modprobe() error {
	var modulespath string = "/etc/modules-load.d"
	if IsExist(modulespath) {
		log.Printf("%s is exist!", modulespath)
	} else {
		err := os.MkdirAll(modulespath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	if err := exec.Command("/bin/sh", "-c", "modprobe br_netfilter").Run(); err != nil {
		return fmt.Errorf("failed to swapoff: %v", err)
	}

	if err := exec.Command("/bin/sh", "-c", "echo 'br_netfilter' > /etc/modules-load.d/kubeocean_br_netfilter.conf").Run(); err != nil {
		return fmt.Errorf("failed to swapoff: %v", err)
	}

	ipvsModList := []string{"ip_vs", "ip_vs_rr", "ip_vs_wrr", "ip_vs_sh"}
	for _, conf := range ipvsModList {
		if err := exec.Command("/bin/sh", "-c", "modprobe", conf).Run(); err != nil {
			return fmt.Errorf("failed to sysctl: %v", err)
		}
	}
	ipvsMod := strings.Join(ipvsModList, "\n")
	if err := exec.Command("/bin/sh", "-c", "echo", ipvsMod, "> /etc/modules-load.d/kubeocean_br_netfilter.conf").Run(); err != nil {
		return fmt.Errorf("failed to swapoff: %v", err)
	}
	return nil
}
