package loginRequest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/tool"
)

//金立

type AmigoLrChannelRet struct {
	R   string `json:"r"`
	Wid string `json:"wid"`
}

type Amigo struct {
	Lr
	channelRet AmigoLrChannelRet
	args       *map[string]interface{}
}

func LrNewAmigo(channelUserId, token string, args *map[string]interface{}) *Amigo {
	ret := new(Amigo)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Amigo) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)

	this.args = args

	this.Method = "POST"
	this.IsHttps = true

	this.Url = "https://id.gionee.com/account/verify.do"
}

func (this *Amigo) InitParam() {
	this.Lr.InitParam()

	apiKey := (*this.args)["AMIGO_APIKEY"].(string)
	secretKey := (*this.args)["AMIGO_SECRETKEY"].(string)

	host := "id.gionee.com"
	port := "443"
	uri := "/account/verify.do"
	method := "POST"

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := tool.RandomAlphanumeric(8)
	context := ts + "\n" + nonce + "\n" + method + "\n" + uri + "\n" + host + "\n" + port + "\n" + "\n"
	sign := base64.StdEncoding.EncodeToString(tool.HmacSHA1Encrypt(context, secretKey))
	format := "MAC id=\"%s\",ts=\"%s\",nonce=\"%s\",mac=\"%s\""

	authorization := fmt.Sprintf(format, apiKey, ts, nonce, sign)

	this.Req.Header("Content-Type", "application/json")
	this.Req.Header("Authorization", authorization)
	this.Req.Body(this.token)
}

func (this *Amigo) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Amigo) CheckChannelRet() bool {
	return this.channelRet.R == "" || this.channelRet.R == "0"
}
