package createOrder

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	// "strconv"
	// "time"
	"u9/models"
	"u9/tool"
)

type CaishenUrlParam struct {
	ProductId string `json:"product_id"`
	Price     int    `json:"price"`
	GameUid   string `json:"game_uid"`
	UId       string `json:"u_id"`
	GameId    string `json:"game_id"`
}

type Caishen struct {
	Cr
	channelParams *map[string]interface{}
	urlParam      CaishenUrlParam
	appSecret     string
}

func CoNewCaishen(lr *models.LoginRequest, orderId, host, urlJsonParam string, channelParams *map[string]interface{}) *Caishen {
	ret := new(Caishen)
	ret.Init(lr, orderId, host, urlJsonParam, channelParams)
	return ret
}

func (this *Caishen) Init(lr *models.LoginRequest, orderId, host, urlJsonParam string, channelParams *map[string]interface{}) {
	this.Cr.Init(lr, orderId, host, urlJsonParam)
	this.channelParams = channelParams
}

func (this *Caishen) InitParam() (err error) {
	if err = json.Unmarshal([]byte(this.urlJsonParam), &this.urlParam); err != nil {
		beego.Trace(err, ":", this.urlJsonParam)
		return
	}

	this.appSecret = (*this.channelParams)["PAYKEY"].(string)

	return
}

func (this *Caishen) ParseChannelRet() (err error) {
	return
}

func (this *Caishen) Response() (err error) {
	return
}

func (this *Caishen) GetResult() (ret string) {

	format := `%s_%s_%d_%s_%s_%s_%s`

	context := fmt.Sprintf(format, this.orderId, this.urlParam.ProductId, this.urlParam.Price,
		this.urlParam.GameUid, this.urlParam.UId, this.urlParam.GameId, this.appSecret)

	sign := tool.Md5([]byte(context))

	ret = sign
	return ret
}
