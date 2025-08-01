package proxy

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

type ReqHandler = func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
type RespHandler = func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response

type Condition = goproxy.RespCondition
