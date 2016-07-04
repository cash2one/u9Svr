package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"time"
	. "u9/api/channel/third/ysdk"
	"u9/api/common"
)

var tencentUrlKeys []string = []string{"data"}

type tencentTradeData_0 struct {
	Amount        string `json:"amount"`
	OrderId       string `json:"orderId"`
	ChannelUserId string `json:"channelUserId"`
	ChannelRet    string `json:"channelRet"`
	ExtParam      string `json:"extParam"`
}

type Tencent struct {
	Base
	common.Request
	ver               string
	tencentChannelRet GetBalanceRet   //tencentTradeData_0
	clientChannelRet  GetBalanceRet   //tencentTradeData_0
	clientExtParam    GetBalanceParam //tencentTradeData_0
}

func (this *Tencent) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = ""

	this.ver = this.urlParams.Get("Ver")

	switch this.ver {
	case "1":
		this.channelTradeData = new(PayParam)
		this.channelRetData = new(PayRet)
	default:
		this.channelTradeData = new(tencentTradeData_0)
		this.channelRetData = nil
	}

	return
}

func (this *Tencent) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := "parseInputParam:err:%+v, \n\n"
			msg := fmt.Sprintf(format, err) + this.Dump()
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	content := this.urlParams.Get("data")
	if this.channelTradeContent, err = url.QueryUnescape(content); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	amount := ""
	discount := ""

	switch this.ver {
	case "1":
		channelTradeData := this.channelTradeData.(*PayParam)
		if this.channelParamKeys["_gameId"], _, this.channelParamKeys["_payKey"], err =
			GetParamName(channelTradeData.LoginType); err != nil {
			return
		}
		this.orderId = channelTradeData.Billno
		this.channelUserId = channelTradeData.ChannelUserId
		this.channelOrderId = ""
		amount = channelTradeData.Amt
	default:
		channelTradeData := this.channelTradeData.(*tencentTradeData_0)
		if err = json.Unmarshal([]byte(channelTradeData.ChannelRet), &this.clientChannelRet); err != nil {
			return
		}

		if err = json.Unmarshal([]byte(channelTradeData.ExtParam), &this.clientExtParam); err != nil {
			return
		}
		if this.channelParamKeys["_gameId"], _, this.channelParamKeys["_payKey"], err =
			GetParamName(this.clientExtParam.LoginType); err != nil {
			return
		}
		this.orderId = channelTradeData.OrderId
		this.channelUserId = channelTradeData.ChannelUserId
		this.channelOrderId = ""
		amount = channelTradeData.Amount
	}

	if err = this.parseChannelParam(); err != nil {
		return
	}

	if err = this.Base.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}

	return
}

func (this *Tencent) CheckSign(params ...interface{}) (err error) {
	return nil
}

func (this *Tencent) InitParam() (err error) {
	if err = this.Request.InitParam(); err != nil {
		return err
	}

	cookie := ""
	loginType := ""
	switch this.ver {
	case "1":
		channelTradeData := this.channelTradeData.(*PayParam)
		loginType = channelTradeData.LoginType

		if cookie, err = GetPayCookie(loginType); err != nil {
			return err
		}
	default:
		loginType = this.clientExtParam.LoginType

		if cookie, err = GetGetBalanceCookie(loginType); err != nil {
			return err
		}

	}
	this.Req.Header("cookie", cookie)
	return nil
}

func (this *Tencent) pay(amount string) (err error) {
	this.Request.Init()
	this.IsHttps = true

	channelTradeData := this.channelTradeData.(*tencentTradeData_0)

	loginType := this.clientExtParam.LoginType

	payParam := new(PayParam)
	payRet := new(PayRet)
	cookie := ""

	format := "pay err:%+v,\n\npayParam:%+v,\n\npayRet:%+v\n\n\n"
	defer func() {
		if err != nil {
			msg := fmt.Sprintf(format, err, payParam, payRet) + this.Dump()
			beego.Warn(msg)
		} else {
			msg := fmt.Sprintf(format, err, payParam, payRet) + this.Dump()
			beego.Trace(msg)
		}
	}()

	payParam.AppId = this.channelParams["_gameId"]
	payParam.PayKey = this.channelParams["_payKey"]

	payParam.Debug = this.clientExtParam.Debug
	payParam.LoginType = this.clientExtParam.LoginType
	payParam.OpenId = this.clientExtParam.OpenId
	payParam.OpenKey = this.clientExtParam.OpenKey
	payParam.PayToken = this.clientExtParam.PayToken
	payParam.Pf = this.clientExtParam.Pf
	payParam.PfKey = this.clientExtParam.PfKey
	payParam.ZoneId = this.clientExtParam.ZoneId

	payParam.Billno = channelTradeData.OrderId
	payParam.Amt = amount

	this.Url = GetPayUrl(payParam)

	if err = this.Request.InitParam(); err != nil {
		return err
	}

	if cookie, err = GetPayCookie(loginType); err != nil {
		return err
	}
	this.Req.Header("cookie", cookie)

	if err = this.GetResponse(); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(this.Result), payRet); err != nil {
		return err
	}

	if tradeState := payRet.Ret == 0; !tradeState {
		tradeFailDesc := `payRet.Ret!=0`
		return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
	}

	return
}

func (this *Tencent) cancelPay() (err error) {
	this.Request.Init()
	this.IsHttps = true

	channelTradeData := this.channelTradeData.(*PayParam)

	loginType := channelTradeData.LoginType

	cancelPayParam := new(CancelPayParam)
	cancelPayRet := new(CancelPayRet)
	cookie := ""

	format := "cancelPay err:%+v,\n\ncancelPayParam:%+v,\n\ncancelPayRet:%+v\n\n\n"
	defer func() {
		if err != nil {
			msg := fmt.Sprintf(format, err, cancelPayParam, cancelPayRet) + this.Dump()
			beego.Warn(msg)
		} else {
			msg := fmt.Sprintf(format, err, cancelPayParam, cancelPayRet) + this.Dump()
			beego.Trace(msg)
		}
	}()

	cancelPayParam.AppId = this.channelParams["_gameId"]
	cancelPayParam.PayKey = this.channelParams["_payKey"]

	cancelPayParam.Debug = channelTradeData.Debug
	cancelPayParam.LoginType = channelTradeData.LoginType
	cancelPayParam.OpenId = channelTradeData.OpenId
	cancelPayParam.OpenKey = channelTradeData.OpenKey
	cancelPayParam.PayToken = channelTradeData.PayToken
	cancelPayParam.Pf = channelTradeData.Pf
	cancelPayParam.PfKey = channelTradeData.PfKey
	cancelPayParam.ZoneId = channelTradeData.ZoneId
	cancelPayParam.Billno = channelTradeData.Billno
	cancelPayParam.Amt = channelTradeData.Amt

	this.Url = GetCancelPayUrl(cancelPayParam)

	if err = this.Request.InitParam(); err != nil {
		return err
	}

	if cookie, err = GetPayCookie(loginType); err != nil {
		return err
	}
	this.Req.Header("cookie", cookie)

	if err = this.GetResponse(); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(this.Result), cancelPayRet); err != nil {
		return err
	}

	if tradeState := cancelPayRet.Ret == 0; !tradeState {
		tradeFailDesc := `cancelPayRet.Ret!=0`
		return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
	}

	return
}

func (this *Tencent) Handle() (err error) {
	format := "handle err:%+v, \n\n"
	defer func() {
		if err != nil {
			if this.lastError != err_noerror {
				this.lastError = err_handleOrder
			}
			msg := fmt.Sprintf(format, err) + this.Dump()
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	this.Request.Init()
	this.IsHttps = true

	switch this.ver {
	case "1":
		channelTradeData := this.channelTradeData.(*PayParam)

		channelTradeData.AppId = this.channelParams["_gameId"]
		channelTradeData.PayKey = this.channelParams["_payKey"]

		this.Url = GetPayUrl(channelTradeData)

		if err = this.InitParam(); err != nil {
			return err
		}

		if err = this.GetResponse(); err != nil {
			return err
		}
		if err = json.Unmarshal([]byte(this.Result), this.channelRetData); err != nil {
			return err
		}

		payRet := this.channelRetData.(*PayRet)
		if tradeState := payRet.Ret == 0; !tradeState {
			tradeFailDesc := `payRet.Ret!=0`
			return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
		}

		if err = this.Base.Handle(); err != nil {
			//call 取消支付
			msg := "pay request:\n\n" + this.Request.Dump() + "\n\n\n"
			beego.Error(msg)

			this.cancelPay()
			return
		}
	default:
		this.clientExtParam.AppId = this.channelParams["_gameId"]
		this.clientExtParam.PayKey = this.channelParams["_payKey"]

		this.Url = GetGetBalanceUrl(&this.clientExtParam)

		if err = this.InitParam(); err != nil {
			return err
		}

		go func() {
			var cur time.Duration = 0
			var dur time.Duration = 15

			for {
				if requestErr := this.GetResponse(); requestErr != nil {
					beego.Warn(requestErr)
				}

				if requestErr := json.Unmarshal([]byte(this.Result), &this.tencentChannelRet); requestErr != nil {
					beego.Warn(this.Request.Dump())
					beego.Warn(requestErr)
				}

				if this.tencentChannelRet.Balance > this.tencentChannelRet.Balance {
					msg := fmt.Sprintf(format, "pay query finish") +
						this.Dump()
					beego.Trace(msg)
					if err = this.Base.Handle(); err != nil {
						return
					}
					amount := strconv.Itoa(this.tencentChannelRet.Balance)
					this.pay(amount)
					return
				} else if cur >= 480 {
					msg := fmt.Sprintf(format, "pay query timeout") +
						this.Dump()
					beego.Warn(msg)
					amount := strconv.Itoa(this.clientChannelRet.Balance)
					this.pay(amount)
					return
				} else {
					cur = cur + dur
					time.Sleep(time.Second * dur)
				}
			}

		}()
	}

	return nil
}

func (this *Tencent) GetResult(params ...interface{}) (ret string) {
	switch this.ver {
	case "1":
		format := `{"result":"%s"}`
		succMsg := "success"
		failMsg := "failure"

		return this.Base.GetResult(format, succMsg, failMsg)
	default:
		format := `{"result":"%s"}`
		succMsg := "success"
		failMsg := "failure"
		return this.Base.GetResult(format, succMsg, failMsg)
	}
	return
}

func (this *Tencent) Dump() (ret string) {
	ret = "1 base:\n\n" + this.Base.Dump() + "\n\n\n" +
		"2 request:\n\n" + this.Request.Dump() + "\n\n\n" + "3 tencent:\n\n"

	switch this.ver {
	case "1":
		format := "ver: %+v,\n\n"
		ret = fmt.Sprintf(format, this.ver)
	default:
		format := "ver: %s,\n\nclientExtParam: %+v,\n\n" +
			"clientChannelRet: %+v,\n\ntencentChannelRet: %+v\n\n"
		ret = ret + fmt.Sprintf(format, this.ver,
			this.clientChannelRet,
			this.clientExtParam,
			this.tencentChannelRet)
	}

	return
}
