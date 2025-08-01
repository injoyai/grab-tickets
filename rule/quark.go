package rule

import (
	"github.com/injoyai/grab-tickets/internal/proxy"
	"github.com/injoyai/logs"
	"net/url"
)

func WithQuark(s *proxy.Proxy) {
	s.OnResponse(
		proxy.HostIs("drive-m.quark.cn"),
		proxy.PathIs("/1/clouddrive/capacity/growth/info"),
	).PrintRequest().OnQuery(func(q url.Values) {
		logs.Debug("kps:", q.Get("kps"))
		logs.Debug("sign:", q.Get("sign"))
		logs.Debug("vcode:", q.Get("vcode"))
	})
}
