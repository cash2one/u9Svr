package channelPayNotify

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"time"
	"u9/api/channel/api/createOrder"
	"u9/api/common"
)

var tencentUrlKeys []string = []string{"data"}

type tencentResData struct {
	Amount        string `json:"amount"`
	OrderId       string `json:"orderId"`
	ChannelUserId string `json:"channelUserId"`
	ChannelRet    string `json:"channelRet"`
	ExtParam      string `json:"extParam"`
}

type Tencent struct {
	Base
	common.Request
	clientChannelRet  createOrder.TencentChannelRet
	clientExtParam    createOrder.TencentExtParam
	tencentChannelRet createOrder.TencentChannelRet
	gameIdKeyName     string
	gameKeyKeyName    string
}

func NewTencent(channelId, productId int, urlParams *url.Values) *Tencent {
	ret := new(Tencent)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Tencent) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &tencentUrlKeys)
	this.data = new(tencentResData)
}

func (this *Tencent) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
			beego.Error(this.urlParams.Encode())
		}
	}()

	tencentData := this.urlParams.Get("data")
	if tencentData, err = url.QueryUnescape(tencentData); err != nil {
		return
	}
	//beego.Trace(fmt.Sprintf("tencentData:%v", tencentData))

	if err = json.Unmarshal([]byte(tencentData), &this.data); err != nil {
		return
	}

	data := this.data.(*tencentResData)
	//beego.Trace(fmt.Sprintf("data:%v", data))

	if err = json.Unmarshal([]byte(data.ChannelRet), &this.clientChannelRet); err != nil {
		return
	}
	//beego.Trace(fmt.Sprintf("clientChannelRet:%v", this.clientChannelRet))

	if err = json.Unmarshal([]byte(data.ExtParam), &this.clientExtParam); err != nil {
		return
	}
	//beego.Trace(fmt.Sprintf("clientExtParam:%v", this.clientExtParam))

	if this.gameIdKeyName, this.gameKeyKeyName, _, err =
		createOrder.GetTencentPayParamName(this.clientExtParam.LoginType); err != nil {
		return
	}

	this.orderId = data.OrderId
	this.channelUserId = data.ChannelUserId

	payAmount := 0
	if payAmount, err = strconv.Atoi(data.Amount); err != nil {
		return err
	} else {
		this.payAmount = payAmount
	}
	return
}

func (this *Tencent) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}

	if err = this.parseChannelGameID(this.gameIdKeyName); err != nil {
		return
	}

	if err = this.parseChannelGameKey(this.gameKeyKeyName); err != nil {
		return
	}
	if err = this.parseChannelPayKey("PAY_KEY"); err != nil {
		return
	}

	if err = this.Base.ParseParam(); err != nil {
		return
	}

	return
}

func (this *Tencent) CheckSign() (err error) {
	return nil
}

func (this *Tencent) InitParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_initParam
			beego.Error(err)
		}
	}()

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
	this.Request.Init()
	this.IsHttps = true

	data := this.data.(*tencentResData)
	this.Url = createOrder.GetTencentPayQueryUrl(this.clientExtParam.Debug, data.ChannelUserId,
		this.clientExtParam.OpenKey, this.clientExtParam.PayToken, this.clientExtParam.Pf,
		this.clientExtParam.PfKey, this.clientExtParam.ZoneId, this.channelGameId, this.channelPayKey)

	if err = this.InitParam(); err != nil {
		return err
	}

	go func() {
		var cur time.Duration = 0
		var dur time.Duration = 15
		for {
			var err error
			if err = this.GetResponse(); err != nil {
				beego.Warn(err)
				//return
			}

			beego.Trace("Result:" + this.Result)
			if err = json.Unmarshal([]byte(this.Result), &this.tencentChannelRet); err != nil {
				beego.Warn(err)
				//return
			}

			if this.tencentChannelRet.SaveAmt > this.clientChannelRet.SaveAmt {
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

func (this *Tencent) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = `{"result":"success"}`
	} else {
		ret = `{"result":"failure"}`
	}
	return
}
