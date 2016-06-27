package channelPayNotify

import (
	"encoding/json"
	"fmt"
	"u9/tool"
)

var kaopuUrlKeys []string = []string{"username", "kpordernum", "ywordernum",
	"status", "amount", "gamename", "sign"}

type kaopuChannelRet struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Sign string `json:"sign"`
}

//靠谱
type KaoPu struct {
	MD5
}

func (this *KaoPu) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &kaopuUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "KAOPU_SECRETKEY"

	this.channelTradeData = nil
	this.channelRetData = new(kaopuChannelRet)

	this.signHandleMethod = ""
	return
}

func (this *KaoPu) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "ywordernum"
	channelUserId_key := "username"
	channelOrderId_key := "kpordernum"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *KaoPu) CheckSign(params ...interface{}) (err error) {
	format := "%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s"
	this.signContent = fmt.Sprintf(format,
		this.channelUserId, this.channelOrderId, this.orderId,
		this.urlParams.Get("status"), this.urlParams.Get("paytype"),
		this.urlParams.Get("amount"), this.urlParams.Get("gameserver"),
		this.urlParams.Get("errdesc"), this.urlParams.Get("paytime"),
		this.urlParams.Get("gamename"), this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *KaoPu) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("status") == "1"
	tradeFailDesc := `urlParam(status)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *KaoPu) GetResult(params ...interface{}) (ret string) {
	channelRetData := this.channelRetData.(*kaopuChannelRet)
	switch this.lastError {
	case err_noerror:
		channelRetData.Code = "1000"
		channelRetData.Msg = "success"

	case err_checkSign:
		channelRetData.Code = "1002"
		channelRetData.Msg = "sign_err"
	case err_orderIsNotExist:
		channelRetData.Code = "1003"
		channelRetData.Msg = "order_err"
	case err_initChannelPayKey:
		channelRetData.Code = "1004"
		channelRetData.Msg = "param_err"
	case err_channelUserIsNotExist:
		channelRetData.Code = "1006"
		channelRetData.Msg = "user_err"
	case err_payAmountError:
		channelRetData.Code = "1009"
		channelRetData.Msg = "amount_err"
	default:
		channelRetData.Code = "1005"
		channelRetData.Msg = "system_err"
	}
	format := "%s|%s"
	content := fmt.Sprintf(format, channelRetData.Code, this.channelParams["_payKey"])
	channelRetData.Sign = tool.Md5([]byte(content))
	data, _ := json.Marshal(channelRetData)
	ret = string(data)

	this.MD5.GetResult()
	return
}
