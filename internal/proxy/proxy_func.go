package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/elazarl/goproxy"
	"io"
	"net/http"
	"strings"
)

// SetCA 设置ca证书
func SetCA(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	goproxy.GoproxyCa = cert
	return nil
}

// ReplaceResponseBody 替换响应体内容
func ReplaceResponseBody(resp *http.Response, old, new string) (*http.Response, error) {
	if resp == nil || resp.Body == nil {
		return resp, nil
	}

	// 读取原始响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// 替换内容
	modified := strings.ReplaceAll(string(bodyBytes), old, new)

	// 设置新的响应体
	resp.Body = io.NopCloser(strings.NewReader(modified))
	resp.ContentLength = int64(len(modified))
	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(modified)))

	return resp, nil
}

/*

1.2 常见的媒体格式
text/html：HTML格式
text/plain：纯文本格式
text/XML：XML格式
image/gif：gif图片格式
image/jped：jpg图片格式
image/png：png图片格式
以application开头的媒体格式类型：

application/xhtml+xml：XHTML格式
application/xml：XML数据格式
application/atom+xml：Atom XML聚合格式
application/json：JSON数据格式【常用】
application/pdf：pdf数据格式
application/msword：Word文档格式
application/octet-stream：二进制流格式（如常见的文件下载）
application/x-www-form-urlencoded：<form encType=" ">中默认的encType，form表单数据编码为key/value格式发送到服务器（表单默认的提交数据格式）【常用】另外一种常见的媒体格式是上传文件时使用的：【常用】
multipart/form-data： ['mʌlti:pɑ:t] 需要在表单进行文件上传时，就需要使用该格式



*/

func NewResponse(r *http.Request, status int, body string, contentType string) *http.Response {
	return goproxy.NewResponse(r, contentType, status, body)
}

func NewTextResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, http.StatusAccepted, text, "text/plain")
}

func NewHtmlResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, http.StatusAccepted, text, "text/html")
}

func NewJsonResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, http.StatusAccepted, text, "application/json")
}
