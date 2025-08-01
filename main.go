package main

import (
	"context"
	"github.com/injoyai/grab-tickets/internal/proxy"
)

/*

	//1. 设置任务

	//2. 开始监听

	//3. 设置手机代理

	//4. 等待手机操作

	//5. 到达指定时间,开始执行任务

	//6. 报告任务结果

*/

func main() {

	s := proxy.Default(
		proxy.WithPort(8888),
		//proxy.WithProxy("http://127.0.0.1:1081"),
	)

	//log.SetOutput(io.Discard)
	s.Verbose = true

	s.OnResponse(
	//proxy.HostLike("(.*)\\.damai\\.cn"),
	).Print(true)

	s.Run(context.Background())

}
