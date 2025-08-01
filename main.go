package main

import (
	"context"
	"github.com/elazarl/goproxy"
	"github.com/injoyai/grab-tickets/internal/proxy"
	"net/http"
)

func main() {

	//1. 设置任务

	//2. 开始监听

	//3. 设置手机代理

	//4. 等待手机操作

	//5. 到达指定时间,开始执行任务

	//6. 报告任务结果

	s := proxy.Default(proxy.WithProxy("http://127.0.0.1:1081"))

	s.OnResponseHostReplace([]string{"www.baidu.com"}, "全球领先", "全球不领先")

	s.Run(context.Background())

}

func RespReplaceByHost(s *proxy.Proxy, host string, old, new string) {
	s.OnResponse(proxy.RespHostIs(host)).
		DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			resp, _ = proxy.RespReplaceBody(resp, old, new)
			return resp
		})
}
