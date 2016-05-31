package loginRequestHandle

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
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
	LRH
	appId      string
	channelRet lenovoChannelRet
}

func NewLenovo() *Lenovo {
	ret := new(Lenovo)
	return ret
}

func (this *Lenovo) Init(param *Param) (err error) {
	this.LRH.Init(param)
	return
}

func (this *Lenovo) Handle() (ret string, err error) {
	this.Method = "GET"
	this.IsHttps = false

	appId := ""
	channelId := this.param.ChannelId
	productId := this.param.ProductId

	pp := new(models.PackageParam)
	if appId, err = pp.GetXmlParam(channelId, productId, "lenovo.open.appid"); err != nil {
		beego.Error(err)
		beego.Error(this.param)
		return
	}

	format := "lpsust=%s&realm=%s"
	//token := url.QueryEscape(this.param.Token)
	url := "http://passport.lenovo.com/interserver/authen/1.2/getaccountid"
	this.Url = url + "?" + fmt.Sprintf(format, this.param.Token, appId)
	beego.Trace(this.Url)

	this.LRH.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("result:", this.Result)
	if err = xml.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	if this.channelRet.Code != "" {
		err = errors.New(this.Result)
		beego.Error(fmt.Sprintf("channelRet:%+v", this.channelRet))
		return
	}

	this.param.ChannelUserId = this.channelRet.AccountID
	this.param.ChannelUserName = this.channelRet.Username

	data, _ := json.Marshal(this.channelRet)

	ret = string(data)
	return
}
