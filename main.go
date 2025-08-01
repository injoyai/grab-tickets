package main

import (
	"context"
	"github.com/elazarl/goproxy"
	"github.com/injoyai/grab-tickets/internal/proxy"
	"github.com/injoyai/logs"
	"net/http"
	"net/url"
)

func main() {

	//1. 设置任务

	//2. 开始监听

	//3. 设置手机代理

	//4. 等待手机操作

	//5. 到达指定时间,开始执行任务

	//6. 报告任务结果

	s := proxy.Default(
		proxy.WithPort(8888),
		proxy.WithProxy("http://127.0.0.1:1081"),
	)

	//s.Verbose = true

	s.OnResponse(
		proxy.HostIs("www.baidu.com"),
	).Do(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response { return resp },
	)

	s.OnRequest(proxy.HostIs("www.baidu.com")).DoNothing()
	s.OnResponse(proxy.HostIs("www.baidu.com")).DoNothing()

	s.OnResponse(
		proxy.HostIs("www.trae.ai"),
		//proxy.PathIs("/cloudide/api/v3/trae/CheckLogin"),
	).ReplaceBody("Login Successful", "Login Failed...")

	s.OnResponse(
		proxy.HostIs("www.baidu.com"),
	).ReplaceBody("全球领先", "全球不领先")

	s.OnRequest(
		proxy.HostIs("www.baidu.com"),
	).ResponseHtmlFile("./prank.html")

	s.OnResponse(
		proxy.HostIs("drive-m.quark.cn"),
		proxy.PathIs("/1/clouddrive/capacity/growth/info"),
	).PrintRequest().OnQuery(func(q url.Values) {
		logs.Debug("kps:", q.Get("kps"))
		logs.Debug("sign:", q.Get("sign"))
		logs.Debug("vcode:", q.Get("vcode"))
	})

	s.Run(context.Background())

}
