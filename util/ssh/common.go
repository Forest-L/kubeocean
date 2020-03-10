package ssh

import (
	"errors"
	"fmt"
	"os"
)

func PrintExecResult(res Result) {
	fmt.Printf("ip=%s\n", res.Ip)
	fmt.Printf("command=%s\n", res.Cmd)
	if res.Err != nil {
		fmt.Printf("return=1\n")
		fmt.Printf("%s\n", res.Err)
	} else {
		fmt.Printf("return=0\n")
		fmt.Printf("%s\n", res.Result)
	}
	fmt.Println("----------------------------------------------------------")
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
