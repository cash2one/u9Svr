package loginRequest

import (
	"encoding/json"
	// "fmt"
	"github.com/astaxie/beego"
	"u9/models"
	"net/url"
	"strings"
	// "u9/tool"
)

//靠谱
type KaoPuResultJson struct{
	CheckUrl  string `json:"checkUrl"`
}

type KaoPuChannelRet struct {
	Code    int 	  `json:"code"`
	Msg     string    `json:"msg"`
	sign 	string    `json:"sign"`
	R   	string    `json:"r"`
}

type KaoPu struct {
	Lr
	kaopuCheckUrl	KaoPuResultJson
	channelRet KaoPuChannelRet
}

func LrNewKaoPu(mlr *models.LoginRequest, args *map[string]interface{}) *KaoPu {
	ret := new(KaoPu)
	ret.Init(mlr, args)
	return ret
}

func (this *KaoPu) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.mlr.Ext,_ =url.QueryUnescape(this.mlr.Ext)
	this.mlr.Ext = strings.Replace(this.mlr.Ext, ",", "&", -1)
	json.Unmarshal([]byte(this.mlr.Ext), &this.kaopuCheckUrl)
	beego.Trace(this.mlr.Ext)
	this.Url = "http://121.201.26.35:8081/Api/CheckUserValidate?"+this.kaopuCheckUrl.CheckUrl
	beego.Trace(this.Url)

}

func (this *KaoPu) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *KaoPu) CheckChannelRet() bool {
	if this.channelRet.Code == 1{
		beego.Trace("check ok")
	}
	return this.channelRet.Code == 1
}
