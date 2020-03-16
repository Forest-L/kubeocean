package ssh

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

func PrintExecResult(res Result) {
	if res.Err != nil {
		fmt.Printf("[%s]  %s  %s\n", res.Ip, res.Cmd, "return: 1")
		fmt.Printf("[Error]  %s\n", res.Err)
	} else {
		fmt.Printf("[%s]  %s  %s\n", res.Ip, res.Cmd, "return: 0")
	}
	//fmt.Println("----------------------------------------------------------")
}

//print pull result
func PrintPullResult(ip, src, dst string, err error) {
	fmt.Println("ip=", ip)
	fmt.Println("command=", "scp "+" root@"+ip+":"+dst+" "+src)
	if err != nil {
		fmt.Printf("return=1\n")
		fmt.Println(err)
	} else {
		fmt.Printf("return=0\n")
		fmt.Printf("Pull from %s to %s ok.\n", dst, src)
	}
	fmt.Println("----------------------------------------------------------")
}

func PrintResults(crs chan Result, ls int, wt *sync.WaitGroup, ccons chan struct{}) {
	for i := 0; i < ls; i++ {
		select {
		case rs := <-crs:
			PrintExecResult(rs)
		case <-time.After(time.Second * Timeout):
			fmt.Printf("getSSHClient error,SSH-Read-TimeOut,Timeout=%ds", Timeout)
		}
		wt.Done()
		<-ccons
	}
}

func CheckResults(crs chan string, ls int, wt *sync.WaitGroup, ccons chan struct{}) {
	for i := 0; i < ls; i++ {
		select {
		case rs := <-crs:
			fmt.Println(rs)
		case <-time.After(time.Second * Timeout):
			fmt.Printf("getSSHClient error,SSH-Read-TimeOut,Timeout=%ds", Timeout)
		}
		wt.Done()
		<-ccons
	}
}

//check path is exit

func FileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return !fi.IsDir()
	}
}

func PathExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return fi.IsDir()
	}
}

func MakePath(path string) error {
	if FileExists(path) {
		return errors.New(path + " is a normal file ,not a dir")
	}

	if !PathExists(path) {
		return os.MkdirAll(path, os.ModePerm)
	} else {
		return nil
	}
}
