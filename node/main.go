package main

import (
	"github.com/injoyai/logs"
)

func main() {

	//1. 获取自身IP地址
	ip, err := GetNetIP()
	logs.PanicErr(err)
	_ = ip

	//2. 与主服务器建立连接
	//redial.TCP()

	//3. 等待主服务器发送任务

	//4. 执行任务

}
