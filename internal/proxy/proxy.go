package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/elazarl/goproxy"
	"github.com/fatih/color"
	"github.com/injoyai/logs"
	"golang.org/x/net/proxy"
	"io"
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

func WithPort(port int) Option {
	return func(p *Proxy) {
		p.SetPort(port)
	}
}

func Default(op ...Option) *Proxy {
	return New(
		WithPort(DefaultPort),
		WithCABytes([]byte(DefaultCrt), []byte(DefaultKey)),
		WithMitm(),
		WithOptions(op...),
	)
}

func New(op ...Option) *Proxy {
	p := &Proxy{
		ProxyHttpServer: goproxy.NewProxyHttpServer(),
		log:             logs.New("").SetFormatter(logs.TimeFormatter).SetColor(color.FgGreen),
		ca:              goproxy.GoproxyCa,
		port:            DefaultPort,
	}
	for _, v := range op {
		v(p)
	}
	return p
}

type Proxy struct {
	*goproxy.ProxyHttpServer
	log   *logs.Entity
	ca    tls.Certificate
	port  int
	debug bool
}

func (this *Proxy) SetPort(port int) {
	this.port = port
}

func (this *Proxy) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
		Handler: this.ProxyHttpServer,
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	this.log.Printf("[信息] [%s] 代理开启成功...\n", srv.Addr)
	c := make(chan struct{})
	defer close(c)
	go func() {
		select {
		case <-ctx.Done():
		case <-c:
		}
		ln.Close()
	}()
	return srv.Serve(ln)
}

// SetOptions 设置选项
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

func (this *Proxy) OnRequestHost(host ...string) *goproxy.ReqProxyConds {
	return this.OnRequest(ReqHostIs(host...))
}

func (this *Proxy) OnResponse(c ...goproxy.RespCondition) *Action {
	return &Action{
		ProxyConds: this.ProxyHttpServer.OnResponse(c...),
		log:        this.log,
	}
}

func (this *Proxy) OnResponseHost(host ...string) *Action {
	return this.OnResponse(RespHostIs(host...))
}

func (this *Proxy) OnResponseHostReplace(host []string, old, new string) {
	this.OnResponse(RespHostIs(host...)).ReplaceBody(old, new)
}

type Condition = goproxy.ReqCondition

type Action struct {
	*goproxy.ProxyConds
	log *logs.Entity
}

func (this *Action) ReplaceBody(old, new string) {
	this.DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp2, err := RespReplaceBody(resp, old, new)
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		return resp2
	})
}

/*
Document 解析body
使用方法参考
https://blog.csdn.net/qq_38334677/article/details/129225231

	doc, err := r.Document()
	if err!=nil{
		return
	}
	//查找标签: 		doc.Find("body,div,...") 多个用,隔开
	//查找ID: 		doc.Find("#id1")
	//查找class: 	doc.Find(".class1")
	//查找属性: 		doc.Find("div[lang]") doc.Find("div[lang=zh]") doc.Find("div[id][lang=zh]")
	//查找子节点: 	doc.Find("body>div")
	//过滤数据: 		doc.Find("div:contains(xxx)")
	//过滤节点: 		dom.Find("span:has(div)")
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})

选择器					说明
Find(“div[lang]”)		筛选含有lang属性的div元素
Find(“div[lang=zh]”)	筛选lang属性为zh的div元素
Find(“div[lang!=zh]”)	筛选lang属性不等于zh的div元素
Find(“div[lang¦=zh]”)	筛选lang属性为zh或者zh-开头的div元素
Find(“div[lang*=zh]”)	筛选lang属性包含zh这个字符串的div元素
Find(“div[lang~=zh]”)	筛选lang属性包含zh这个单词的div元素，单词以空格分开的
Find(“div[lang$=zh]”)	筛选lang属性以zh结尾的div元素，区分大小写
Find(“div[lang^=zh]”)	筛选lang属性以zh开头的div元素，区分大小写
*/
func (this *Action) Document(f func(resp *http.Response, ctx *goproxy.ProxyCtx, doc *goquery.Document)) {
	this.ProxyConds.DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if f == nil {
			return resp
		}
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		resp.Body = io.NopCloser(bytes.NewReader(bs))
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		f(resp, ctx, doc)
		return resp
	})
	return
}
