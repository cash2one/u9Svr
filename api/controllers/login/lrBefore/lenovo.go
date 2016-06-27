package lrBefore

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

type lenovoChannelRet struct {
	AccountID string `xml:"AccountID" json:"AccountID"`
	Username  string `xml:"Username" json:"Username`
	DeviceID  string `xml:"DeviceID" json:"DeviceID`
	Verified  string `xml:"verified" json:"verified`
	Thirdname string `xml:"Thirdname" json:"Thirdname`
	Code      string `xml:"Code" json:"Code`
	Timestamp string `xml:"Timestamp" json:"Timestamp`
	Message   string `xml:"Message" json:"Message`
	Source    string `xml:"Source" json:"Source`
	URL       string `xml:"URL" json:"URL`
}

type Lenovo struct {
	base
}

func NewLenovo() *Lenovo {
	ret := new(Lenovo)
	return ret
}

func (this *Lenovo) Init(param *Param) (err error) {
	this.base.Init(param)
	this.channelRet = new(lenovoChannelRet)
	return
}

func (this *Lenovo) Exec() (ret string, err error) {
	this.Method = "GET"
	this.IsHttps = false

	defer func() {
		if err != nil {
			format := "exec: err:%v"
			msg := fmt.Sprintf(format, err) + this.Dump()
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	appId := ""
	if appId, err = this.getChannelParam("lenovo.open.appid"); err != nil {
		return
	}

	format := "lpsust=%s&realm=%s"
	//token := url.QueryEscape(this.param.Token)
	url := "http://passport.lenovo.com/interserver/authen/1.2/getaccountid"
	this.Url = url + "?" + fmt.Sprintf(format, this.param.Token, appId)

	this.base.InitParam()
	if err = this.GetResponse(); err != nil {
		return
	}

	if err = xml.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		return
	}

	channelRet := this.channelRet.(*lenovoChannelRet)

	if channelRet.Code != "" {
		msg := fmt.Sprintf(`channelRet.Code!=""`)
		err = errors.New(msg)
		return
	}

	this.param.ChannelUserId = channelRet.AccountID
	this.param.ChannelUserName = channelRet.Username

	data, _ := json.Marshal(this.channelRet)

	ret = string(data)
	return
}
