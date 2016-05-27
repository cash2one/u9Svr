package loginRequestHandle

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
	LRH
	channelRet huaweiChannelRet
}

func NewHuawei() *Huawei {
	ret := new(Huawei)
	return ret
}

func (this *Huawei) Init(param *Param) (err error) {
	this.LRH.Init(param)
	return
}

func (this *Huawei) Handle() (ret string, err error) {
	this.IsHttps = true
	this.Method = "GET"
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	token := url.QueryEscape(this.param.Token)
	format := "?nsp_svc=OpenUP.User.getInfo&nsp_ts=%s&access_token=%s"
	this.Url = "https://api.vmall.com/rest.php" + fmt.Sprintf(format, ts, token)

	this.LRH.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("result:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	if this.channelRet.UserID == "" {
		err = errors.New(this.channelRet.Error)
		beego.Error(fmt.Sprintf("channelRet:%+v", this.channelRet))
		return
	}

	this.param.ChannelUserId = this.channelRet.UserID
	this.param.ChannelUserName = this.channelRet.UserName

	data, _ := json.Marshal(this.channelRet)
	ret = string(data)
	return
}
