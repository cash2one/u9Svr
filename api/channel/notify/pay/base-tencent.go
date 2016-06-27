package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"time"
	"u9/api/channel/api/createOrder"
	"u9/api/common"
)

var tencentUrlKeys []string = []string{"data"}

type tencentTradeData struct {
	Amount        string `json:"amount"`
	OrderId       string `json:"orderId"`
	ChannelUserId string `json:"channelUserId"`
	ChannelRet    string `json:"channelRet"`
	ExtParam      string `json:"extParam"`
}

type Tencent struct {
	Base
	common.Request
	clientChannelRet createOrder.TencentChannelRet
	clientExtParam   createOrder.TencentExtParam
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

	this.channelTradeData = new(tencentTradeData)
	this.channelRetData = nil

	return
}

func (this *Tencent) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := `ParseInputParam:err:%v, clientChannelRet:%v, clientExtParam:%v`
			msg := fmt.Sprintf(format,
				err, this.clientChannelRet, this.clientExtParam)
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

	channelTradeData := this.channelTradeData.(*tencentTradeData)
	if err = json.Unmarshal([]byte(channelTradeData.ChannelRet), &this.clientChannelRet); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(channelTradeData.ExtParam), &this.clientExtParam); err != nil {
		return
	}

	gameId_keyName := ""
	//gameKey_keyName := ""
	payKey_keyName := ""

	if gameId_keyName, _, payKey_keyName, err =
		createOrder.GetTencentPayParamName(this.clientExtParam.LoginType); err != nil {
		return
	}

	if this.channelParams["_gameId"], err = this.getChannelParam(gameId_keyName); err != nil {
		return
	}

	if this.channelParams["_payKey"], err = this.getChannelParam(payKey_keyName); err != nil {
		return
	}

	this.orderId = channelTradeData.OrderId
	this.channelUserId = channelTradeData.ChannelUserId
	this.channelOrderId = ""

	amount := channelTradeData.Amount
	discount := ""

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
	//beego.Trace("InitParam")
	if err = this.Request.InitParam(); err != nil {
		return err
	}

	cookie := ""
	if cookie, err = createOrder.GetTencentPayQueryCookie(this.clientExtParam.LoginType); err != nil {
		return err
	}
	this.Req.Header("cookie", cookie)
	return nil
}

func (this *Tencent) Handle() (err error) {
	var tencentChannelRet createOrder.TencentChannelRet
	defer func() {
		if err != nil {
			this.lastError = err_handleOrder

			format := `err:%v, channelParams:%v, ` +
				`clientChannelRet:%v, this.clientExtParam:%v, ` +
				`tencentChannelRet:%v, ` +
				`queryUrl:%s queryResult:%s`
			msg := fmt.Sprintf(format,
				err, this.channelParams,
				this.clientChannelRet, this.clientExtParam,
				tencentChannelRet,
				this.Url, this.Result)
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	this.Request.Init()
	this.IsHttps = true

	channelTradeData := this.channelTradeData.(*tencentTradeData)
	this.Url = createOrder.GetTencentPayQueryUrl(
		this.clientExtParam.Debug,
		channelTradeData.ChannelUserId,
		this.clientExtParam.OpenKey,
		this.clientExtParam.PayToken,
		this.clientExtParam.Pf,
		this.clientExtParam.PfKey,
		this.clientExtParam.ZoneId,
		this.channelParams["_gameId"],
		this.channelParams["_payKey"])

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

			beego.Trace("Result:" + this.Result)
			if requestErr := json.Unmarshal([]byte(this.Result), &tencentChannelRet); requestErr != nil {
				beego.Warn(requestErr)
			}

			if tencentChannelRet.SaveAmt > this.clientChannelRet.SaveAmt {
				beego.Trace("tencent pay query is finish.")
				if err = this.Base.Handle(); err != nil {
					return
				}
				return
			} else if cur >= 120 {
				beego.Trace("tencent pay query is timeout.")
				return
			} else {
				cur = cur + dur
				time.Sleep(time.Second * dur)
			}
		}

	}()
	return nil
}

func (this *Tencent) GetResult(params ...interface{}) (ret string) {
	format := `{"result":"%s"}`
	succMsg := "success"
	failMsg := "failure"
	return this.Base.GetResult(format, succMsg, failMsg)
}
