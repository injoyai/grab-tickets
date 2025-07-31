package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/injoyai/logs"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
)

type Option func(*Proxy)

func WithCA(certFile, keyFile string) Option {
	return func(p *Proxy) {
		err := p.SetCA(certFile, keyFile)
		logs.PrintErr(err)
	}
}

func WithCABytes(certFile, keyFile []byte) Option {
	return func(p *Proxy) {
		err := p.SetCABytes(certFile, keyFile)
		logs.PrintErr(err)
	}
}

func WithProxy(u string) Option {
	return func(p *Proxy) {
		err := p.SetProxy(u)
		logs.PrintErr(err)
	}
}

func WithDebug(b ...bool) Option {
	return func(p *Proxy) {
		p.Debug(b...)
	}
}

func WithMitm() Option {
	return func(p *Proxy) {
		p.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			return &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&p.ca)}, host
		}))
	}
}

func WithOptions(op ...Option) Option {
	return func(p *Proxy) {
		p.SetOptions(op...)
	}
}

func Default(op ...Option) *Proxy {
	return New(
		DefaultPort,
		WithCABytes([]byte(DefaultCrt), []byte(DefaultKey)),
		WithMitm(),
		WithOptions(op...),
	)
}

func New(port int, op ...Option) *Proxy {
	p := &Proxy{
		ProxyHttpServer: goproxy.NewProxyHttpServer(),
		ca:              goproxy.GoproxyCa,
		port:            port,
	}
	for _, v := range op {
		v(p)
	}
	return p
}

type Proxy struct {
	*goproxy.ProxyHttpServer
	ca    tls.Certificate
	port  int
	debug bool
}

func (this *Proxy) Run(ctx context.Context) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", this.port), this.ProxyHttpServer)
}

func (this *Proxy) SetOptions(op ...Option) {
	for _, v := range op {
		v(this)
	}
}

// SetCABytes 设置ca证书
func (this *Proxy) SetCABytes(crt, key []byte) error {
	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return err
	}
	this.ca = cert
	return nil
}

// SetCA 设置ca证书
func (this *Proxy) SetCA(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	this.ca = cert
	return nil
}

// SetProxy 设置代理
func (this *Proxy) SetProxy(u string) error {
	if len(u) == 0 {
		this.ProxyHttpServer.Tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		}
		return nil
	}

	proxyUrl, err := url.Parse(u)
	if err != nil {
		return err
	}
	t := &http.Transport{}
	switch proxyUrl.Scheme {
	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(proxyUrl, this)
		if err != nil {
			return err
		}
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	default: //"http", "https"
		t.Proxy = http.ProxyURL(proxyUrl)
	}
	this.ProxyHttpServer.Tr = t
	return nil
}

// Debug 打印通讯数据
func (this *Proxy) Debug(b ...bool) {
	this.debug = len(b) == 0 || b[0]
}

// Dial 实现接口
func (this *Proxy) Dial(network, addr string) (c net.Conn, err error) {
	return this.ProxyHttpServer.ConnectDial(network, addr)
}
