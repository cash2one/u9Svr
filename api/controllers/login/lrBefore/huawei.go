package lrBefore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"time"
)

type huaweiChannelRet struct {
	Gender          int    `json:"gender"`
	UserID          string `json:"userID"`
	UserName        string `json:"userName"`
	LanguageCode    string `json:"languageCode"`
	UserState       int    `json:"userState"`
	UserValidStatus int    `json:"userValidStatus"`
	Error           string `json:"error"`
}

type Huawei struct {
	base
}

func NewHuawei() *Huawei {
	ret := new(Huawei)
	return ret
}

func (this *Huawei) Init(param *Param) (err error) {
	this.base.Init(param)
	this.channelRet = new(huaweiChannelRet)
	return
}

func (this *Huawei) Exec() (ret string, err error) {
	this.IsHttps = true
	this.Method = "GET"
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	token := url.QueryEscape(this.param.Token)
	format := "?nsp_svc=OpenUP.User.getInfo&nsp_ts=%s&access_token=%s"
	this.Url = "https://api.vmall.com/rest.php" + fmt.Sprintf(format, ts, token)

	this.base.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*huaweiChannelRet)

	if channelRet.UserID == "" {
		err = errors.New(channelRet.Error)
		beego.Error(fmt.Sprintf("channelRet:%+v", channelRet))
		return
	}

	this.param.ChannelUserId = channelRet.UserID
	this.param.ChannelUserName = channelRet.UserName

	ret = this.Result
	return
}
