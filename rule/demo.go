package rule

import "github.com/injoyai/grab-tickets/internal/proxy"

func WithPrank(s *proxy.Proxy) {
	s.OnRequest().ResponseHtmlFile("./prank.html")
}
