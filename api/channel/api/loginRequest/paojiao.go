package loginRequest

import (
	"encoding/json"
	// "errors"
	"fmt"
	"github.com/astaxie/beego"
	// "net/url"
	// "strconv"
	"u9/models"
	"u9/tool"
)

type ChannelRet struct {
	Msg 	string    	`json:"msg"`
	Code  	string 		`json:"code"`
	Data    DataJson 	`json:"data"`
}
type DataJson struct{
	CreatedTime 	string    	`json:"createdTime"`
	NiceName  		string 		`json:"niceName"`
	Token     		string 		`json:"token"`
	Uid 			int    		`json:"uid"`
	UserName  		string 		`json:"userName"`
	Avatar     		string 		`json:"avatar"`
}

type PaoJiao struct {
	Lr
	appId      string
	secretKey  string
	channelRet ChannelRet
}

func LrNewPaoJiao(mlr *models.LoginRequest, args *map[string]interface{}) *PaoJiao {
	ret := new(PaoJiao)
	ret.Init(mlr, args)
	return ret
}

func (this *PaoJiao) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.appId = (*args)["PAOJIAO_APPID"].(string)
	this.secretKey = (*args)["PAOJIAO_SERVERSECRET"].(string)
	sign := tool.Md5([]byte(this.appId + this.mlr.Token + this.secretKey))

	this.Url = "http://ng.sdk.paojiao.cn/api/user/token.do?"
	format := "token=%s&appId=%s&sign=%s"
	this.Url = this.Url + fmt.Sprintf(format, this.mlr.Token ,this.appId, sign)
	beego.Trace(this.Url)
}

func (this *PaoJiao) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		return
	}

	// if this.appId != this.channelRet.AppID {
	// 	err = errors.New("channel.AppId isn't equal with xmlParms.AppId.")
	// 	return
	// }

	// content := this.appId + strconv.Itoa(this.channelRet.ResultCode) + this.channelRet.Content + this.secretKey
	// if sign := tool.Md5([]byte(content)); sign != this.channelRet.Sign {
	// 	msg := fmt.Sprintf("sign(%s) is equal channelRet's sign(%s)", sign, this.channelRet.Sign)
	// 	err = errors.New(msg)
	// 	return
	// }

	// if this.channelRet.ResultMsg, err = url.QueryUnescape(this.channelRet.ResultMsg); err != nil {
	// 	return
	// }

	// var enByte []byte
	// if enByte, err = base64.StdEncoding.DecodeString(this.channelRet.Content); err != nil {
	// 	return
	// }
	// this.channelRet.Content = string(enByte)

	//beego.Trace(fmt.Sprintf("%+v",this.channelRet))
	return
}

func (this *PaoJiao) CheckChannelRet() bool {
	return this.channelRet.Code == "1"
}
