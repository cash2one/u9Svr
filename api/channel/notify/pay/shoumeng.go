package channelPayNotify

import (
	// "encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var shoumengUrlKeys []string = []string{"orderId", "uid", "amount", "coOrderId", "success"}

const (
	err_shoumengParsePayKey   = 10101
	err_shoumengResultFailure = 10102
)

//手盟
type ShouMeng struct {
	Base
	gameId   string
	payKey   string
	packetId string
	u9UserId string
	// tokenJson TokenJson
}

// type TokenJson struct {
// 	Login_account string `json:"login_account"`
// 	Session_id    string `json:"session_id"`
// }

func NewShouMeng(channelId, productId int, urlParams *url.Values) *ShouMeng {
	ret := new(ShouMeng)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *ShouMeng) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &shoumengUrlKeys)
}

func (this *ShouMeng) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_shoumengParsePayKey
			beego.Trace(err)
		}
	}()
	// this.gameId, err = this.getPackageParam("SHOUMENG_GAME_ID")
	// this.packetId, err = this.getPackageParam("SHOUMENG_PACKET_ID")
	this.payKey, err = this.getPackageParam("SHOUMENG_SERCETKEY")
	return
}

func (this *ShouMeng) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("coOrderId")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("orderId")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("amount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *ShouMeng) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("success"); result != "0" {
		this.callbackRet = err_shoumengResultFailure
	}
	return
}

func (this *ShouMeng) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parsePayKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
	// json.Unmarshal([]byte(this.loginRequest.Token), &TokenJson)
	// this.u9UserId = this.loginRequest.Userid
}

func (this *ShouMeng) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	// format := "user_id=%s&login_account=%s&game_id=%s&packet_id=%s&game_server_id=%s" +
	// 	"&game_role_id=%s&order_id=%s&game_coin=%s&total_fee=%s&pay_time=%s" +
	// 	"&pay_result=%s&secret=%s"

	// context := fmt.Sprintf(format,
	// 	this.channelUserId, this.tokenJson.Login_account, this.gameId, this.packetId, this.urlParams.Get("serverId"),
	// 	this.u9UserId, this.urlParams.Get("orderId"), this.urlParams.Get("result"),
	// 	this.urlParams.Get("ext"), this.payKey)

	format := "orderId=%s&uid=%s&amount=%s&coOrderId=%s&success=%s&secret=%s"
	context := fmt.Sprintf(format, this.channelOrderId, this.channelUserId, this.urlParams.Get("amount"),
		this.orderId, this.urlParams.Get("success"), this.payKey)

	if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {
		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *ShouMeng) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "SUCCESS"
	} else {
		ret = "FAILURE"
	}
	return
}

/*
  signature rule: md5("order=xxxx&money=xxxx&mid=xxxx&time=xxxx&result=x&ext=xxx&key=xxxx")
  test url:
  http://192.168.0.185/api/channelPayNotify/1000/101/?
  order=test20160116172500359&
  money=100.00&
  mid=test10086001&
  time=20160116172500&
  result=1&
  ext=game20160116175128772&
  signature=8f00a109716e819bfe0afb695c1addf1
*/
