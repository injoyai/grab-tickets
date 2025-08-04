package main

import (
	"context"
	"github.com/injoyai/grab-tickets/internal/proxy"
)

func main() {

	s := proxy.Default(
		proxy.WithPort(8888),
		//proxy.WithProxy("http://127.0.0.1:1081"),
		proxy.WithProxyPac("http://127.0.0.1:1081", proxy.Domains),
	)

	//log.SetOutput(io.Discard)
	//s.Verbose = true

	s.OnResponse(
	//proxy.HostLike("(.*)\\.damai\\.cn"),
	).PrintHost() //.Print()

	s.Run(context.Background())

}
