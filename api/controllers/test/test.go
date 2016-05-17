package test

import (
	"bytes"
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io"
	"net/http"
)

var s = "s"

type Test struct {
	beego.Controller
}

func (this *Test) Test1() {

	ret := "test"

	defer func() {
		this.Data["json"] = ret
		//beego.Trace(ret)
		beego.Error("defer")
		this.ServeJSON(true)
	}()

	beego.Warn(s)
	s = s + `s`
	test1()
}

func test1() {
	var resp *http.Response
	beego.Warn("Response1")
	url := `https://ysdk.qq.com`

	req := httplib.Get(url)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	beego.Warn("Response2")
	resp, err := req.Response()
	beego.Warn("Response3")
	if err != nil {
		beego.Warn(err)
		beego.Warn("ResponseErr")
		return
	}
	beego.Warn("Response4")
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, resp.Body); err != nil {
		beego.Warn(err)
		return
	}
	beego.Warn("Response5")
	bytes := buffer.Bytes()
	result := string(bytes)

	beego.Warn(result)
	beego.Warn("Response6")
	return
}

/*
  test url:
  http://192.168.0.185/api/gamePayRequest/?
  ProductId=1000&
  UserId=23c16f5323755132272fba79ab2e11d8&
  ProductOrderId=game20160114142841787&
  Amount=32&
  CallbackUrl=http://www.baidu.com
*/
