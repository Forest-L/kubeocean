package install

import (
	"fmt"
	"github.com/pixiake/kubeocean/util"
	"github.com/pixiake/kubeocean/util/cluster"
	"github.com/pixiake/kubeocean/util/ssh-bak"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func SystemInit(isLocal bool, host *cluster.NodeCfg) error {

	sysctllist := []string{"net.ipv4.ip_local_reserved_ports=30000-32767", "net.bridge.bridge-nf-call-iptables=1", "net.bridge.bridge-nf-call-arptables=1", "net.bridge.bridge-nf-call-ip6tables"}
	for _, conf := range sysctllist {
		cmd := fmt.Sprintf("sysctl %s", conf)

		if isLocal {
			if err := exec.Command("/bin/sh", "-c", cmd).Run(); err != nil {
				return fmt.Errorf("failed to sysctl: %v", err)
			}
		} else {
			key := ""
			portInt, err := strconv.Atoi(host.Port)
			if err != nil {
				return err
			}
			cmdList := []string{cmd}
			err_exec := ssh_bak.DosshRun(host.User, host.Password, host.Address, key, cmdList, portInt, nil)
			if err_exec != nil {
				return err_exec
			}
		}

	}
	return nil
}

func SwapOff(isLocal bool, host *cluster.NodeCfg) error {

	cmd := "swapoff -a"
	disableswap := fmt.Sprint(`sed -i /^[^#]*swap*/s/^/\#/g /etc/fstab`)
	if isLocal {
		if err := exec.Command("/bin/sh", "-c", cmd).Run(); err != nil {
			return fmt.Errorf("failed to swapoff: %v", err)
		}
		if err := exec.Command("/bin/sh", "-c", disableswap).Run(); err != nil {
			return fmt.Errorf("failed to swapoff: %v", err)
		}
	} else {
		key := ""
		portInt, err := strconv.Atoi(host.Port)
		if err != nil {
			return err
		}
		cmdList := []string{cmd, disableswap}
		err_exec := ssh_bak.DosshRun(host.User, host.Password, host.Address, key, cmdList, portInt, nil)
		if err_exec != nil {
			return err_exec
		}
	}

	return nil
}

func Modprobe() error {
	var modulespath string = "/etc/modules-load.d"
	if util.IsExist(modulespath) {
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
