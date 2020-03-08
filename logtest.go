package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	// 设置日志等级
	log.SetLevel(log.DebugLevel)
	// 设置日志输出到什么地方去
	// 将日志输出到标准输出，就是直接在控制台打印出来。
	log.SetOutput(os.Stdout)
	// 设置为true则显示日志在代码什么位置打印的
	//log.SetReportCaller(true)

	// 设置日志以json格式输出， 如果不设置默认以text格式输出
	log.SetFormatter(&log.TextFormatter{})

	// 打印日志
	log.Debug("调试信息")
	log.Info("提示信息")
	log.Warn("警告信息")
	log.Error("错误信息")
	//log.Panic("致命错误")
	//
	// 为日志加上字段信息，log.Fields其实就是map[string]interface{}类型的别名
	log.WithFields(log.Fields{
		"user_id":    1001,
		"ip":         "123.12.12.11",
		"request_id": "kana012uasdb8a918gad712",
	}).Info("用户登陆失败.")

	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	// Calls os.Exit(1) after logging
	log.Fatal("Bye.")
	// Calls panic() after logging
	log.Panic("I'm bailing.")
}
