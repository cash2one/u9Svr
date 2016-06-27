package loginRequest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"net/url"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type OppoChannelRet struct {
	ResultCode   string `json:"resultCode"`
	ResultMsg    string `json:"resultMsg"`
	LoginToken   string `json:"loginToken"`
	Ssoid        string `json:"ssoid"`
	AppKey       string `json:"appKey"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
	CreateTime   string `json:"createTime"`
	UserStatus   string `json:"userStatus"`
}

type Oppo struct {
	Lr
	channelRet  OppoChannelRet
	signContent string
	sign        string
}

func LrNewOppo(mlr *models.LoginRequest, args *map[string]interface{}) *Oppo {
	ret := new(Oppo)
	ret.Init(mlr, args)
	return ret
}

func (this *Oppo) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)

	appSecret := (*args)["OPPO_APPSECRET"].(string)
	appkey := (*args)["OPPO_APPKEY"].(string)
	escapeToken := url.QueryEscape(this.mlr.Token)
	format := "http://i.open.game.oppomobile.com/gameopen/user/fileIdInfo?fileId=%s&token=%s"
	this.Url = fmt.Sprintf(format, this.mlr.ChannelUserid, escapeToken)

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	this.signContent = "oauthConsumerKey=" + appkey + "&oauthToken=" + escapeToken +
		"&oauthSignatureMethod=HMAC-SHA1" + "&oauthTimestamp=" + timeStamp +
		"&oauthNonce=" + timeStamp + strconv.Itoa(r.Intn(10)) + "&oauthVersion=1.0&"

	this.sign = base64.StdEncoding.EncodeToString(tool.HmacSHA1Encrypt(this.signContent, appSecret+"&"))
}

func (this *Oppo) InitParam() (err error) {
	if err = this.Lr.InitParam(); err != nil {
		return
	}
	this.Req.Header("param", this.signContent)
	this.Req.Header("oauthsignature", this.sign)
	return
}

func (this *Oppo) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Oppo) CheckChannelRet() bool {
	return this.channelRet.Ssoid == this.mlr.ChannelUserid
}
