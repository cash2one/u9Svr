package common

import (
	"bytes"
	"crypto/tls"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io"
	"net/http"
)

type Request struct {
	Url     string
	Req     *httplib.BeegoHTTPRequest
	Method  string
	Result  string
	IsHttps bool
}

func (this *Request) Init() {
	this.Method = "GET"
}

func (this *Request) InitParam() {
	this.initRequest()
}

func (this *Request) initRequest() {
	if this.Method == "GET" {
		this.Req = httplib.Get(this.Url)
	} else if this.Method == "POST" {
		this.Req = httplib.Post(this.Url)
	}
	if this.IsHttps {
		this.Req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
}

func (this *Request) Response() (err error) {
	var resp *http.Response
	resp, err = this.Req.Response()
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, resp.Body); err != nil {
		return
	}
	bytes := buffer.Bytes()
	this.Result = string(bytes)
	return
}
