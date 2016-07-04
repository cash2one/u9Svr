package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego"
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
	this.IsHttps = false
}

func (this *Request) InitParam() (err error) {
	if this.Method == "GET" {
		this.Req = httplib.Get(this.Url)
	} else if this.Method == "POST" {
		this.Req = httplib.Post(this.Url)
	}
	if this.IsHttps {
		this.Req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return nil
}

func (this *Request) GetResponse() (err error) {
	var resp *http.Response
	format := "getResponse:err:%v \n"
	resp, err = this.Req.Response()
	if err != nil {
		msg := fmt.Sprintf(format, err) + this.Dump()
		beego.Error(msg)
		return
	}
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, resp.Body); err != nil {
		msg := fmt.Sprintf(format, err) + this.Dump()
		beego.Error(msg)
		return
	}
	bytes := buffer.Bytes()
	this.Result = string(bytes)
	return
}

func (this *Request) Dump() (ret string) {
	format := "method:%s,\nurl:%s,\n\nheader:%+v,\n\nresult:%+s"
	ret = fmt.Sprintf(format,
		this.Method, this.Url, this.Req.GetRequest().Header, this.Result)
	return
}
